package queries

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
)

func (q *Queries) Download(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	// Get query
	qu, ok := q.base.CalcQuery(c, opts[0].(string))
	if !ok {
		return
	}

	// Get names
	var names []string
	err := q.db.Select(&names, `SELECT name FROM elements WHERE guild=$1 AND id=ANY($2)`, c.Guild(), pq.Array(qu.Elements))
	if err != nil {
		q.base.Error(c, err)
		return
	}

	// Send
	dm, err := c.Dg().UserChannelCreate(c.Author().User.ID)
	if err != nil {
		q.base.Error(c, err)
		return
	}
	msg := sevcord.NewMessage(fmt.Sprintf("📄 Query **%s**:", qu.Name)).
		AddFile("query.txt", "text/plain", strings.NewReader(strings.Join(names, "\n")))
	_, err = c.Dg().ChannelMessageSendComplex(dm.ID, msg.Dg())
	if err != nil {
		q.base.Error(c, err)
		return
	}

	// Respond
	c.Respond(sevcord.NewMessage("Sent query in DMs! 📄"))
}
