package queries

import (
	"fmt"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
	"github.com/lib/pq"
)

func (q *Queries) Info(ctx sevcord.Ctx, opts []any) {
	ctx.Acknowledge()

	// Get query
	qu, ok := q.base.CalcQuery(ctx, opts[0].(string))
	if !ok {
		return
	}

	// Get progress
	var common int
	err := q.db.QueryRow(`SELECT COALESCE(array_length($2 & (SELECT inv FROM inventories WHERE guild=$1 AND "user"=$3), 1), 0)`, ctx.Guild(), pq.Array(qu.Elements), ctx.Author().User.ID).Scan(&common)
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

	// Add query data
	switch qu.Kind {
	case types.QueryKindElement:
		emb = emb.AddField("Kind", "Element", true)
		name, err := q.base.GetName(ctx.Guild(), int(qu.Data["elem"].(float64)))
		if err != nil {
			q.base.Error(ctx, err)
			return
		}
		emb = emb.AddField("Element", name, true)

	case types.QueryKindCategory:
		emb = emb.AddField("Kind", "Category", true)
		emb = emb.AddField("Category", qu.Data["cat"].(string), true)

	case types.QueryKindProducts:
		emb = emb.AddField("Kind", "Products", true)
		emb = emb.AddField("Query", qu.Data["query"].(string), true)

	case types.QueryKindParents:
		emb = emb.AddField("Kind", "Parents", true)
		emb = emb.AddField("Query", qu.Data["query"].(string), true)

	case types.QueryKindInventory:
		emb = emb.AddField("Kind", "Inventory", true)
		emb = emb.AddField("User", fmt.Sprintf("<@%s>", qu.Data["user"].(string)), true)

	case types.QueryKindElements:
		emb = emb.AddField("Kind", "Elements", true)

	case types.QueryKindRegex:
		emb = emb.AddField("Kind", "Regex", true)
		emb = emb.AddField("Query", qu.Data["query"].(string), true)
		emb = emb.AddField("Regex", "```"+qu.Data["regex"].(string)+"```", false)

	case types.QueryKindComparison:
		emb = emb.AddField("Kind", "Comparison", true)
		emb = emb.AddField("Field", "`"+qu.Data["field"].(string)+"`", true)
		emb = emb.AddField("Operator", qu.Data["typ"].(string), true)
		emb = emb.AddField("Value", fmt.Sprintf("%v", qu.Data["value"]), true)
	}

	// Respond
	ctx.Respond(sevcord.NewMessage("").AddEmbed(emb))
}
