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

	case types.PollKindImage:
		err := b.elemImageSuccess(p, newsFunc)
		if err != nil {
			log.Println("image error", err)
		}

	case types.PollKindCategorize:
		err := b.categorizeSuccess(p, newsFunc)
		if err != nil {
			log.Println("categorize error", err)
		}

	case types.PollKindUncategorize:
		err := b.unCategorizeSuccess(p, newsFunc)
		if err != nil {
			log.Println("uncategorize error", err)
		}

	case types.PollKindComment:
		err := b.elemMarkSuccess(p, newsFunc)
		if err != nil {
			log.Println("comment error", err)
		}
	}

	b.deletePoll(p, dg)
}

func (b *Polls) pollContextMsg(p *types.Poll) string {
	return fmt.Sprintf("(Lasted **%s** • By <@%s>)", time.Since(p.CreatedOn).Round(time.Second).String(), p.Creator)
}
