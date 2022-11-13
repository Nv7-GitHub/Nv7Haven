package elements

import (
	"database/sql"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
)

func (e *Elements) Suggest(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	e.base.IncrementCommandStat(c, "suggest")

	// Autocapitalization
	autocap := true
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
		c.Respond(sevcord.NewMessage(res.Message))
		return
	}

	// Get res
	var idV any
	if id == -1 {
		idV = opts[0].(string)
	} else {
		idV = float64(id)
	}

	// Create suggestion
	err = e.polls.CreatePoll(c, &types.Poll{
		Kind: types.PollKindCombo,
		Data: types.PollData{
			"els":    util.Map(v, func(a int) float64 { return float64(a) }),
			"result": idV,
		},
	})
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Make text
	names, err := e.base.GetNames(v, c.Guild())
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

	// Message
	c.Respond(sevcord.NewMessage(text.String()))
}
