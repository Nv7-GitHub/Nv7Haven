package elements

import (
	"database/sql"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/timing"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
)

func (e *Elements) Suggest(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	timer := timing.GetTimer("suggest")

	// Autocapitalization
	autocap := false
	if strings.ToLower(opts[0].(string)) == opts[0].(string) { // Autocapitalize if all lowercase
		autocap = true
	}
	if opts[1] != nil {
		autocap = opts[1].(bool)
	}
	if autocap {
		opts[0] = util.Capitalize(opts[0].(string))
	}

	// Check if result exists already
	var id int
	var name string
	err := e.db.QueryRow(`SELECT id, name FROM elements WHERE guild=$1 AND LOWER(name)=$2`, c.Guild(), strings.ToLower(opts[0].(string))).Scan(&id, &name)
	if err != nil {
		if err == sql.ErrNoRows {
			id = -1
		} else {
			e.base.Error(c, err)
			return
		}
	}

	// Get els
	v, res := e.base.GetCombCache(c)
	if !res.Ok {
		c.Respond(res.Response())
		return
	}

	// Check if combo has result
	var exists bool
	err = e.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM combos WHERE guild=$1 AND els=$2)`, c.Guild(), pq.Array(v.Elements)).Scan(&exists)
	if err != nil {
		e.base.Error(c, err)
		return
	}
	if exists {
		c.Respond(sevcord.NewMessage("This combo already has a result! " + types.RedCircle))
		return
	}

	// Get res
	var idV any
	if id == -1 {
		idV = opts[0].(string)

		// Check if valid
		var ok types.Resp
		idV, ok = base.CheckName(idV.(string))
		if !ok.Ok {
			c.Respond(ok.Response())
			return
		}
	} else {
		idV = float64(id)
	}

	// Create suggestion
	res = e.polls.CreatePoll(c, &types.Poll{
		Kind: types.PollKindCombo,
		Data: types.PgData{
			"els":    util.Map(v.Elements, func(a int) any { return float64(a) }),
			"result": idV,
		},
	})
	if !res.Ok {
		c.Respond(res.Response())
		return
	}

	// Make text
	names, err := e.base.GetNames(v.Elements, c.Guild())
	if err != nil {
		e.base.Error(c, err)
		return
	}
	text := &strings.Builder{}
	text.WriteString("Suggested **")
	text.WriteString(strings.Join(names, " + "))
	text.WriteString(" = ")
	if id != -1 {
		text.WriteString(name)
	} else {
		text.WriteString(opts[0].(string))
	}
	text.WriteString("** ")
	if id != -1 {
		text.WriteString("🌟")
	} else {
		text.WriteString("✨")
	}

	timer.Stop()

	// Message
	c.Respond(sevcord.NewMessage(text.String()))
}
