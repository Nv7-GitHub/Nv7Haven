package treecmds

import (
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/bwmarrin/discordgo"
)

type TreeCmds struct {
	*eodb.Data

	base *base.Base
	dg   *discordgo.Session
}

func NewTreeCmds(data *eodb.Data, dg *discordgo.Session, base *base.Base) *TreeCmds {
	return &TreeCmds{
		Data: data,

		base: base,
		dg:   dg,
	}
}
