package categories

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
)

func (c *Categories) Info(ctx sevcord.Ctx, opts []any) {
	ctx.Acknowledge()

	// Get category
	var cat types.Category
	err := c.db.Get(&cat, "SELECT * FROM categories WHERE LOWER(name)=$1 AND guild=$2", strings.ToLower(opts[0].(string)), ctx.Guild())
	if err != nil {
		c.base.Error(ctx, err, "Category **"+opts[0].(string)+"** doesn't exist!")
		return
	}

	// Get progress
	var common int
	err = c.db.QueryRow(`SELECT COALESCE(array_length(elements & (SELECT inv FROM inventories WHERE guild=$1 AND "user"=$3), 1), 0) FROM categories WHERE guild=$1 AND name=$2`, ctx.Guild(), cat.Name, ctx.Author().User.ID).Scan(&common)
	if err != nil {
		c.base.Error(ctx, err)
		return
	}

	// Description
	description := "📝 **Mark**\n" + cat.Comment

	// Embed
	emb := sevcord.NewEmbed().
		Title(cat.Name+" Info").
		Description(description).
		Color(cat.Color).
		AddField("💼 Element Count", humanize.Comma(int64(len(cat.Elements))), true).
		AddField("📊 Progress", humanize.FormatFloat("", float64(common)/float64(len(cat.Elements))*100)+"%", true)

	// Optional things
	if cat.Image != "" {
		emb = emb.Thumbnail(cat.Image)
	}
	if cat.Commenter != "" {
		emb = emb.AddField("💬 Commenter", fmt.Sprintf("<@%s>", cat.Commenter), true)
	}
	if cat.Colorer != "" {
		emb = emb.AddField("🖌️ Colorer", fmt.Sprintf("<@%s>", cat.Colorer), true)
	}
	if cat.Imager != "" {
		emb = emb.AddField("🖼️ Imager", fmt.Sprintf("<@%s>", cat.Imager), true)
	}
	emb = emb.AddField("🎨 Color", util.FormatHex(cat.Color), true)
	// Respond
	ctx.Respond(sevcord.NewMessage("").AddEmbed(emb))
}
