package base

import (
	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/bwmarrin/discordgo"
)

type Base struct {
	*eodb.Data

	dg *discordgo.Session
}

func NewBase(data *eodb.Data, dg *discordgo.Session) *Base {
	return &Base{
		Data: data,
		dg:   dg,
	}
}
