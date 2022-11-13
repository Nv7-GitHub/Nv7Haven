package polls

import (
	"fmt"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

func (b *Polls) pollSuccess(p *types.Poll, dg *discordgo.Session) {
	fmt.Println("Success", p)

	b.deletePoll(p, dg)
}
