package basecmds

import (
	"sync"

	"github.com/Nv7-Github/Nv7Haven/db"
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

type BaseCmds struct {
	dat  map[string]types.ServerDat
	lock *sync.RWMutex
	base *base.Base
	dg   *discordgo.Session
	db   *db.DB
}

func NewBaseCmds(dat map[string]types.ServerDat, base *base.Base, dg *discordgo.Session, db *db.DB, lock *sync.RWMutex) *BaseCmds {
	return &BaseCmds{
		dat:  dat,
		lock: lock,
		base: base,
		dg:   dg,
		db:   db,
	}
}
