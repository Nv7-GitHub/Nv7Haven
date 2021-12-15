package polls

import (
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/bwmarrin/discordgo"
)

type Polls struct {
	*eodb.Data

	dg   *discordgo.Session
	base *base.Base
}

func NewPolls(data *eodb.Data, dg *discordgo.Session, base *base.Base) *Polls {
	return &Polls{
		Data: data,

		dg:   dg,
		base: base,
	}
}
