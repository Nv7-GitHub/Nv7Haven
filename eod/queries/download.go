package queries

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
)

func (q *Queries) Download(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	sort := "id"
	if opts[1] != nil {
		sort = opts[1].(string)
	}
	if sort == "found" {
		e := types.Fail("Cannot sort by found!")
		c.Respond(e.Response())
		return
	}

	// Get query
	qu, ok := q.base.CalcQuery(c, opts[0].(string))
	if !ok {
		return
	}
	sql := `SELECT name FROM elements WHERE guild=$1 AND id=ANY($2) ORDER BY ` + types.SortSql[sort]
	if opts[2] != nil {
		sql = `SELECT name, ` + types.PostfixSql[opts[2].(string)] + ` postfix FROM elements WHERE guild=$1 AND id=ANY($2) ORDER BY ` + types.SortSql[sort]
	}

	// Get names
	var names []struct {
		Name    string `db:"name"`
		Postfix string `db:"postfix"`
	}
	err := q.db.Select(&names, sql, c.Guild(), pq.Array(qu.Elements))
	if err != nil {
		q.base.Error(c, err)
		return
	}

	// Make text
	out := &strings.Builder{}
	for _, name := range names {
		out.WriteString(name.Name)
		if opts[2] != nil {
			out.WriteString(" - ")
			out.WriteString(types.GetPostfixVal(name.Postfix, opts[2].(string)))
		}
		out.WriteRune('\n')
	}

	// Send
	dm, err := c.Dg().UserChannelCreate(c.Author().User.ID)
	if err != nil {
		q.base.Error(c, err)
		return
	}
	msg := sevcord.NewMessage(fmt.Sprintf("📄 Query **%s**:", qu.Name)).
		AddFile("query.txt", "text/plain", strings.NewReader(out.String()), out.Len())
	_, err = c.Dg().ChannelMessageSendComplex(dm.ID, msg.Dg())
	if err != nil {
		q.base.Error(c, err)
		return
	}

	// Respond
	c.Respond(sevcord.NewMessage("Sent query in DMs! 📄"))
}
