package queries

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
)

func (q *Queries) createCmd(c sevcord.Ctx, name string, kind types.QueryKind, data map[string]any) {
	c.Acknowledge()

	// Check if name exists
	var edit bool
	err := q.db.Get(&edit, "SELECT EXISTS(SELECT 1 FROM queries WHERE LOWER(name)=$1 AND guild=$2)", strings.ToLower(name), c.Guild())
	if err != nil {
		q.base.Error(c, err)
		return
	}
	if edit {
		err = q.db.QueryRow(`SELECT name FROM queries WHERE LOWER(name)=$1 AND guild=$2`, strings.ToLower(name), c.Guild()).Scan(&name)
		if err != nil {
			q.base.Error(c, err)
			return
		}
	} else {
		// Fix name
		var ok types.Resp
		name, ok = base.CheckName(name)
		if !ok.Ok {
			c.Respond(sevcord.NewMessage(ok.Message + " " + types.RedCircle))
			return
		}
	}

	// Check if data already exists
	var existsName string
	err = q.db.QueryRow("SELECT name FROM queries WHERE data@>$1 AND data<@$1 AND guild=$2", types.PgData(data), c.Guild()).Scan(&existsName)
	if err != nil && err != sql.ErrNoRows {
		q.base.Error(c, err)
		return
	}
	if err == nil {
		c.Respond(sevcord.NewMessage(fmt.Sprintf("Query **%s** already exists with this data! "+types.RedCircle, existsName)))
		return
	}

	// Create
	err = q.polls.CreatePoll(c, &types.Poll{
		Kind: types.PollKindQuery,
		Data: types.PgData{
			"query": name,
			"edit":  edit,
			"kind":  string(kind),
			"data":  any(data),
		},
	})
	if err != nil {
		q.base.Error(c, err)
		return
	}

	// Respond
	word := "create"
	if edit {
		word = "edit"
	}
	c.Respond(sevcord.NewMessage(fmt.Sprintf("Suggested to %s query! 🧮", word)))
}

func (q *Queries) CreateElementsCmd(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	q.createCmd(c, opts[0].(string), types.QueryKindElements, map[string]any{})
}
