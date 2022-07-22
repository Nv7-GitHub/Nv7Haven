package eod

import (
	"fmt"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *EoD) optimize(m types.Msg, rsp types.Rsp) {
	b.Data.RLock()
	defer b.Data.RUnlock()

	id := rsp.Message(fmt.Sprintf("Optimizing [0/%d]...", len(b.DB)))

	taken := time.Duration(0)
	i := 0
	lastUpdated := 0
	for _, db := range b.DB {
		if (len(db.Elements) > 100) || (i-lastUpdated > 10) { // If it has enough elements to take a significant amount of time
			hasEdited := false
			gld, err := b.dg.Guild(db.Guild)
			if err == nil {
				isCommunity := false
				for _, feature := range gld.Features {
					if feature == "COMMUNITY" {
						isCommunity = true
						break
					}
				}

				if isCommunity {
					b.dg.ChannelMessageEdit(m.ChannelID, id, fmt.Sprintf("<@%s> Optimizing **%s**... [%d/%d]", m.Author.ID, gld.Name, i+1, len(b.DB)))
					hasEdited = true
				}
			}

			if !hasEdited {
				b.dg.ChannelMessageEdit(m.ChannelID, id, fmt.Sprintf("<@%s> Optimizing... [%d/%d]", m.Author.ID, i+1, len(b.DB)))
			}

			lastUpdated = i
		}

		start := time.Now()
		err := db.Optimize()
		if rsp.Error(err) {
			return
		}
		err = db.OptimizeCats()
		if rsp.Error(err) {
			return
		}
		taken += time.Since(start)

		i++
	}

	b.dg.ChannelMessageEdit(m.ChannelID, id, fmt.Sprintf("<@%s> Optimized in **%s**.", m.Author.ID, taken.String()))
}
