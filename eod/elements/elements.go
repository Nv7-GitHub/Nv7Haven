package elements

import (
	"github.com/Nv7-Github/Nv7Haven/db"
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/polls"
	"github.com/bwmarrin/discordgo"
)

type Elements struct {
	*eodb.Data

	polls *polls.Polls
	db    *db.DB
	base  *base.Base
	dg    *discordgo.Session
}

func NewElements(data *eodb.Data, polls *polls.Polls, db *db.DB, base *base.Base, dg *discordgo.Session) *Elements {
	return &Elements{
		Data: data,

		polls: polls,
		db:    db,
		base:  base,
		dg:    dg,
	}
}
