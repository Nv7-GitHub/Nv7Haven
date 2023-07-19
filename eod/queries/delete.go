package queries

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
)

func (q *Queries) DeleteQuery(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	var name string
	err := q.db.QueryRow("SELECT name FROM queries WHERE LOWER(name)=$1 AND guild=$2", strings.ToLower(opts[0].(string)), c.Guild()).Scan(&name)
	if err != nil {
		q.base.Error(c, err, "Query **"+opts[0].(string)+"** doesn't exist!")
		return
	}

	// Check if used
	var usedName string
	err = q.db.QueryRow(`SELECT name FROM queries WHERE data->>'query'=$1 OR data->>'left'=$1 OR data->>'right'=$1 AND guild=$2`, name, c.Guild()).Scan(&usedName)
	if err != nil && err != sql.ErrNoRows {
		q.base.Error(c, err)
		return
	}
	if err == nil {
		c.Respond(sevcord.NewMessage(fmt.Sprintf("Cannot delete query **%s** because it is used in query **%s**! "+types.RedCircle, name, usedName)))
		return
	}

	// Delete
	res := q.polls.CreatePoll(c, &types.Poll{
		Kind: types.PollKindDelQuery,
		Data: types.PgData{
			"query": name,
		},
	})
	if !res.Ok {
		c.Respond(res.Response())
		return
	}

	// Respond
	c.Respond(sevcord.NewMessage("Suggested to delete query! 🗑️"))
}
