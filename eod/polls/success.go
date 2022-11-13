package polls

import (
	"fmt"
	"log"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

func (b *Polls) pollSuccess(p *types.Poll, dg *discordgo.Session) {
	// Get news channel
	var news string
	var votecnt int
	err := b.db.QueryRow(`SELECT news, votecnt FROM config WHERE guild=$1`, p.Guild).Scan(&news, &votecnt)
	if err != nil {
		log.Println("news err", err)
		return
	}
	newsFunc := func(msg string) {
		if votecnt != 0 && float32(p.Downvotes)/float32(votecnt) >= 0.3 { // Controversial
			msg += " 🌩️"
		}

		// Send
		_, err = dg.ChannelMessageSend(news, msg)
		if err != nil {
			log.Println("news err", err)
		}
	}

	switch p.Kind {
	case types.PollKindCombo:
		err := b.elemCreate(p, newsFunc)
		if err != nil {
			log.Println("create error", err)
		}
	}

	b.deletePoll(p, dg)
}

func (b *Polls) pollContextMsg(p *types.Poll) string {
	return fmt.Sprintf("(Lasted **%s** • By <@%s>)", time.Since(p.CreatedOn).Round(time.Second).String(), p.Creator)
}
