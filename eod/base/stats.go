package base

import (
	"fmt"
	"time"

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
	fmt.Println("D")
	cmds += int64(b.getMem(c).CommandStatsTODOCnt) // Include count not pushed to DB
	fmt.Println("E")

	// Embed
	emb := sevcord.NewEmbed().Title("Stats").
		Color(7506394). // Blurple
		AddField("🔢 Element Count", humanize.Comma(elemcnt), true).
		AddField("🔄 Combination Count", humanize.Comma(combcnt), true).
		AddField("🧑‍🤝‍🧑 User Count", humanize.Comma(invcnt), true).
		AddField("🔍 Elements Found", humanize.Comma(foundcnt), true).
		AddField("📁 Elements Categorized", humanize.Comma(categorized), true).
		AddField("👨‍💻 Commands Used", humanize.Comma(cmds), true)

	fmt.Println("F")

	// Respond
	fmt.Println(c.Respond(sevcord.NewMessage("").
		AddEmbed(emb).
		AddComponentRow(
			sevcord.NewButton("View More Stats", sevcord.ButtonStyleLink, "", "").WithEmoji(sevcord.ComponentEmojiCustom("stats", "1197216720897712209", false)).SetURL("https://nv7haven.com/eod"),
		)))
	fmt.Println("G")
}

func (b *Base) SaveStats() {
	var lastTime time.Time
	err := b.db.QueryRow("SELECT time FROM stats ORDER BY time DESC LIMIT 1").Scan(&lastTime)
	if err != nil {
		fmt.Println(err)
	}

	if time.Since(lastTime).Hours() > 24 {
		var categorized int
		var found int
		var elemCnt int
		var comboCnt int
		var userCnt int
		var serverCnt int
		var commandStatsRaw []struct {
			Command string `db:"command"`
			Count   int    `db:"count"`
		}

		// Select
		err := b.db.QueryRow(`SELECT
		(SELECT SUM(array_length(elements, 1)) FROM categories),
		(SELECT SUM(array_length(inv, 1)) FROM inventories),
		(SELECT COUNT(*) FROM elements),
		(SELECT COUNT(*) FROM combos),
		(SELECT COUNT(*) FROM inventories),
		(SELECT COUNT(DISTINCT(guild)) FROM config)
		`).Scan(&categorized, &found, &elemCnt, &comboCnt, &userCnt, &serverCnt)
		if err != nil {
			fmt.Println(err)
		}
		err = b.db.Select(&commandStatsRaw, "SELECT command, count FROM command_stats")
		if err != nil {
			fmt.Println(err)
		}

		// Create command stats
		commandStats := make(map[string]int)
		for _, v := range commandStatsRaw {
			commandStats[v.Command] += v.Count
		}

		// Insert
		_, err = b.db.Exec("INSERT INTO stats VALUES ($1, $2, $3, $4, $5, $6, $7)", time.Now(), elemCnt, comboCnt, userCnt, found, categorized, serverCnt)
		if err != nil {
			fmt.Println(err)
		}

		// Save command stats
		for k, v := range commandStats {
			_, err = b.db.Exec("INSERT INTO command_stats VALUES ($1, $2, $3)", time.Now(), k, v)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
