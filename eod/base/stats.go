package base

import (
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
)

func (b *Base) Stats(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	var elemcnt int64
	err := b.db.QueryRow("SELECT COUNT(*) FROM elements WHERE guild=$1", c.Guild()).Scan(&elemcnt)
	if err != nil {
		b.Error(c, err)
		return
	}

	var combcnt int64
	err = b.db.QueryRow("SELECT COUNT(*) FROM combos WHERE guild=$1", c.Guild()).Scan(&combcnt)
	if err != nil {
		b.Error(c, err)
		return
	}

	var invcnt int64
	err = b.db.QueryRow("SELECT COUNT(*) FROM inventories WHERE guild=$1", c.Guild()).Scan(&invcnt)
	if err != nil {
		b.Error(c, err)
		return
	}

	var foundcnt int64
	err = b.db.QueryRow("SELECT SUM(array_length(inv, 1)) FROM inventories WHERE guild=$1", c.Guild()).Scan(&foundcnt)
	if err != nil {
		b.Error(c, err)
		return
	}

	var categorized int64
	err = b.db.QueryRow("SELECT SUM(array_length(elements, 1)) FROM categories WHERE guild=$1", c.Guild()).Scan(&categorized)
	if err != nil {
		b.Error(c, err)
		return
	}

	var cmds int64
	err = b.db.QueryRow("SELECT SUM(count) FROM command_stats WHERE guild=$1", c.Guild()).Scan(&cmds)
	if err != nil {
		b.Error(c, err)
		return
	}

	// Embed
	emb := sevcord.NewEmbed().Title("Stats").
		Color(7506394). // Blurple
		AddField("Element Count", humanize.Comma(elemcnt), true).
		AddField("Combination Count", humanize.Comma(combcnt), true).
		AddField("User Count", humanize.Comma(invcnt), true).
		AddField("Elements Found", humanize.Comma(foundcnt), true).
		AddField("Elements Categorized", humanize.Comma(categorized), true).
		AddField("Commands Used", humanize.Comma(cmds), true)

	// Respond
	c.Respond(sevcord.NewMessage("").
		AddEmbed(emb).
		AddComponentRow(
			sevcord.NewButton("View More Stats", sevcord.ButtonStyleLink, "", "").SetURL("https://nv7haven.com/eod"),
		))
}
