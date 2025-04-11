package queries

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
)

// Format: user|query|sort|postfix|type
func (q *Queries) DownloadHandler(c sevcord.Ctx, params string) {

	parts := strings.Split(params, "|")
	if len(parts) != 5 {
		return
	}
	// Get query
	qu, ok := q.base.CalcQuery(c, parts[1])
	if !ok {
		return
	}
	sort := parts[2]

	conttext := ""
	if sort == "found" {
		conttext = `, id=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$3) cont`
	}
	sql := fmt.Sprintf(`SELECT name,id %s FROM elements WHERE guild=$1 AND id=ANY($2) ORDER BY %s `, conttext, types.SortSql[sort])
	//sql := `SELECT name ` + conttext + ` FROM elements WHERE guild=$1 AND id=ANY($2) ORDER BY ` + types.SortSql[sort]
	if parts[3] != "" {
		sql = fmt.Sprintf(`SELECT name, id %s postfix %s FROM elements WHERE guild = $1 AND id=ANY($2) ORDER BY %s `, conttext, types.PostfixSql[parts[3]], types.SortSql[sort])
	}

	// Get names
	var names []struct {
		Name    string `db:"name"`
		ID      int32  `db:"id"`
		Postfix string `db:"postfix"`
		Cont    string `db:"cont"`
	}
	var err error
	if sort != "found" {
		err = q.db.Select(&names, sql, c.Guild(), pq.Array(qu.Elements))
	} else {
		err = q.db.Select(&names, sql, c.Guild(), pq.Array(qu.Elements), c.Author().User.ID)
	}

	if err != nil {
		q.base.Error(c, err)
		return
	}

	// Make text
	out := &strings.Builder{}
	for _, name := range names {
		if parts[4] == "id" {
			out.WriteString(fmt.Sprintf("#%d", name.ID))
		} else {
			out.WriteString(name.Name)
		}

		if parts[3] != "" {
			out.WriteString(" - ")
			out.WriteString(types.GetPostfixVal(name.Postfix, parts[3]))
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
func (q *Queries) Download(c sevcord.Ctx, opts []any) {

	c.Acknowledge()

	sort := "id"
	if opts[1] != nil {
		sort = opts[1].(string)
	}
	postfix := ""
	if opts[2] != nil {
		postfix = opts[2].(string)
	}
	q.DownloadHandler(c, fmt.Sprintf("%s|%s|%s|%s|id", c.Author().GuildID, opts[0].(string), sort, postfix))
}
func (q *Queries) DownloadIDs(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	sort := "id"
	if opts[1] != nil {
		sort = opts[1].(string)
	}
	postfix := ""
	if opts[2] != nil {
		postfix = opts[2].(string)
	}
	q.DownloadHandler(c, fmt.Sprintf("%s|%s|%s|%s|id", c.Author().User.ID, opts[0].(string), sort, postfix))
}
