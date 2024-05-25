package queries

import (
	"fmt"
	"log"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/lib/pq"
)

func (q *Queries) editNewsMessage(c sevcord.Ctx, message string) {
	var news string
	err := q.db.QueryRow(`SELECT news FROM config WHERE guild=$1`, c.Guild()).Scan(&news)
	if err != nil {
		log.Println("news err", err)
		return
	}
	_, err = c.Dg().ChannelMessageSend(news, fmt.Sprintf("🔨 "+message))
	if err != nil {
		log.Println("news err", err)
	}
}

func (q *Queries) editCmd(c sevcord.Ctx, opts []any, field string, name ...string) {
	c.Acknowledge()
	qu, ok := q.base.CalcQuery(c, opts[0].(string))
	if !ok {
		return
	}

	_, err := q.db.Exec("UPDATE elements SET "+field+"=$3 WHERE id=ANY($1) AND guild=$2", pq.Array(qu.Elements), c.Guild(), opts[1])
	if err != nil {
		q.base.Error(c, err)
		return
	}
	nameV := field
	if len(name) > 0 {
		nameV = name[0]
	}
	c.Respond(sevcord.NewMessage("Successfully edited elements in query " + nameV + "! ✅"))
	q.editNewsMessage(c, fmt.Sprintf("Edited Query Elements %s - **%s**", util.Capitalize(nameV), qu.Name))
}

func (q *Queries) EditElementImageCmd(c sevcord.Ctx, opts []any) {
	q.editCmd(c, opts, "image")
}

func (q *Queries) EditElementColorCmd(c sevcord.Ctx, opts []any) {
	q.editCmd(c, opts, "color")
}

func (q *Queries) EditElementCommentCmd(c sevcord.Ctx, opts []any) {
	q.editCmd(c, opts, "comment")
}

func (q *Queries) EditElementCreatorCmd(c sevcord.Ctx, opts []any) {
	opts[1] = opts[1].(*discordgo.User).ID
	q.editCmd(c, opts, "creator")
}

func (q *Queries) EditElementCreatedonCmd(c sevcord.Ctx, opts []any) {
	opts[1] = time.Unix(opts[1].(int64), 0)
	q.editCmd(c, opts, "createdon", "creation date")
}
