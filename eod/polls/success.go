package polls

import (
	"log"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

func (b *Polls) pollSuccess(p *types.Poll, dg *discordgo.Session) {
	switch p.Kind {
	case types.PollKindCombo:
		err := b.elemCreate(p)
		if err != nil {
			log.Println("create error", err)
		}
	}

	b.deletePoll(p, dg)
}
