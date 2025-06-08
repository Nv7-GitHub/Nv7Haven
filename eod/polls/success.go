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

	case types.PollKindColor:
		err := b.elemColorSuccess(p, newsFunc)
		if err != nil {
			log.Println("color error", err)
		}

	case types.PollKindCatImage:
		err := b.catImageSuccess(p, newsFunc)
		if err != nil {
			log.Println("cat image error", err)
		}

	case types.PollKindCatComment:
		err := b.catMarkSuccess(p, newsFunc)
		if err != nil {
			log.Println("cat comment error", err)
		}

	case types.PollKindCatColor:
		err := b.catColorSuccess(p, newsFunc)
		if err != nil {
			log.Println("cat color error", err)
		}

	case types.PollKindQuery:
		err := b.queryCreateSuccess(p, newsFunc)
		if err != nil {
			log.Println("query create error", err)
		}

	case types.PollKindDelQuery:
		err := b.queryDeleteSuccess(p, newsFunc)
		if err != nil {
			log.Println("query delete error", err)
		}

	case types.PollKindQueryComment:
		err := b.queryMarkSuccess(p, newsFunc)
		if err != nil {
			log.Println("query comment error", err)
		}

	case types.PollKindQueryColor:
		err := b.queryColorSuccess(p, newsFunc)
		if err != nil {
			log.Println("query color error", err)
		}

	case types.PollKindQueryImage:
		err := b.queryImageSuccess(p, newsFunc)
		if err != nil {
			log.Println("query image error", err)
		}
	case types.PollKindCatRename:
		err := b.catRenameSuccess(p, newsFunc)
		if err != nil {
			log.Println("cat rename error", err)
		}
	case types.PollKindQueryRename:
		err := b.queryRenameSuccess(p, newsFunc)
		if err != nil {
			log.Println("query rename error", err)
		}
	}

	b.deletePoll(p, dg)
}

func (b *Polls) pollContextMsg(p *types.Poll) string {
	return fmt.Sprintf("(Lasted **%s** • By <@%s>)", time.Since(p.CreatedOn).Round(time.Second).String(), p.Creator)
}
