package categories

import (
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/polls"
	"github.com/bwmarrin/discordgo"
)

type Categories struct {
	*eodb.Data

	base  *base.Base
	dg    *discordgo.Session
	polls *polls.Polls
}

func NewCategories(data *eodb.Data, base *base.Base, dg *discordgo.Session, polls *polls.Polls) *Categories {
	return &Categories{
		Data: data,

		base:  base,
		dg:    dg,
		polls: polls,
	}
}
