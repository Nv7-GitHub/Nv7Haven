package basecmds

import (
	"github.com/Nv7-Github/Nv7Haven/db"
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/bwmarrin/discordgo"
)

type BaseCmds struct {
	*eodb.Data

	db   *db.DB
	base *base.Base
	dg   *discordgo.Session
}

func NewBaseCmds(base *base.Base, db *db.DB, dg *discordgo.Session, data *eodb.Data) *BaseCmds {
	return &BaseCmds{
		Data: data,

		base: base,
		db:   db,
		dg:   dg,
	}
}
