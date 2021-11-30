package basecmds

import (
	"fmt"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/bwmarrin/discordgo"
)

func (b *BaseCmds) StatsCmd(m types.Msg, rsp types.Rsp) {
	rsp.Acknowledge()

	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("Guild not setup yet!")
		return
	}
	gd, err := b.dg.State.Guild(m.GuildID)
	if rsp.Error(err) {
		return
	}

	found := 0
	dat.Lock.RLock()
	for _, val := range dat.Inventories {
		found += len(val.Elements)
	}

	categorized := 0
	for _, val := range dat.Categories {
		categorized += len(val.Elements)
	}

	rsp.Message(fmt.Sprintf("Element Count: **%s**\nCombination Count: **%s**\nMember Count: **%s**\nElements Found: **%s**\nElements Categorized: **%s**", util.FormatInt(len(dat.Elements)), util.FormatInt(len(dat.Combos)), util.FormatInt(gd.MemberCount), util.FormatInt(found), util.FormatInt(categorized)), discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label: "View More Stats",
				URL:   "https://nv7haven.com/?page=eod",
				Style: discordgo.LinkButton,
			},
		},
	})
	dat.Lock.RUnlock()
}

// takes time, found, categorized
var saveStatsQuery = `INSERT INTO eod_stats VALUES (?, (SELECT COUNT(1) FROM eod_elements), (SELECT COUNT(1) FROM eod_combos), (SELECT COUNT(DISTINCT user) FROM eod_inv), ?, ?, (SELECT COUNT(DISTINCT guild) FROM eod_serverdata))`

func (b *BaseCmds) SaveStats() {
	var lastTime int64
	err := b.db.QueryRow("SELECT time FROM eod_stats ORDER BY time DESC LIMIT 1").Scan(&lastTime)
	if err != nil {
		fmt.Println(err)
	}

	if time.Since(time.Unix(lastTime, 0)).Hours() > 24 {
		b.lock.RLock()
		categorized := 0
		found := 0
		for _, dat := range b.dat {
			for _, val := range dat.Categories {
				categorized += len(val.Elements)
			}
			for _, val := range dat.Inventories {
				found += len(val.Elements)
			}
		}
		b.lock.RUnlock()

		_, err = b.db.Exec(saveStatsQuery, time.Now().Unix(), found, categorized)
		if err != nil {
			fmt.Println(err)
		}
	}
}
