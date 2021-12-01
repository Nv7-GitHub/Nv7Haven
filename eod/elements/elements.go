package elements

import (
	"sync"

	"github.com/Nv7-Github/Nv7Haven/db"
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/polls"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

type Elements struct {
	dat   map[string]types.ServerDat
	lock  *sync.RWMutex
	polls *polls.Polls
	db    *db.DB
	base  *base.Base
	dg    *discordgo.Session
}

func NewElements(dat map[string]types.ServerDat, lock *sync.RWMutex, polls *polls.Polls, db *db.DB, base *base.Base, dg *discordgo.Session) *Elements {
	return &Elements{
		dat:   dat,
		lock:  lock,
		polls: polls,
		db:    db,
		base:  base,
		dg:    dg,
	}
}
