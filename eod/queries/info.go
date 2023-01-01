package queries

import (
	"fmt"

	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
	"github.com/lib/pq"
)

func (q *Queries) Info(ctx sevcord.Ctx, opts []any) {
	ctx.Acknowledge()

	// Get query
	qu, err := q.CalcQuery(ctx, opts[0].(string))
	if err != nil {
		q.base.Error(ctx, err)
		return
	}

	// Get progress
	var common int
	err = q.db.QueryRow(`SELECT COALESCE(array_length($2 & (SELECT inv FROM inventories WHERE guild=$1 AND "user"=$3), 1), 0)`, ctx.Guild(), pq.Array(qu.Elements), ctx.Author().User.ID).Scan(&common)
	if err != nil {
		q.base.Error(ctx, err)
		return
	}

	// Description
	description := "**Mark**\n" + qu.Comment

	// Embed
	emb := sevcord.NewEmbed().
		Title(qu.Name+" Info").
		Description(description).
		Color(qu.Color).
		AddField("Element Count", humanize.Comma(int64(len(qu.Elements))), true).
		AddField("Progress", humanize.FormatFloat("", float64(common)/float64(len(qu.Elements))*100)+"%", true)

	// Optional things
	if qu.Image != "" {
		emb = emb.Thumbnail(qu.Image)
	}
	if qu.Commenter != "" {
		emb = emb.AddField("Commenter", fmt.Sprintf("<@%s>", qu.Commenter), true)
	}
	if qu.Colorer != "" {
		emb = emb.AddField("Colorer", fmt.Sprintf("<@%s>", qu.Colorer), true)
	}
	if qu.Imager != "" {
		emb = emb.AddField("Imager", fmt.Sprintf("<@%s>", qu.Imager), true)
	}

	// Respond
	ctx.Respond(sevcord.NewMessage("").AddEmbed(emb))
}
