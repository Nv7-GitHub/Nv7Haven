package base

import (
	"sync"

	"github.com/Nv7-Github/Nv7Haven/db"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

type Base struct {
	dat  map[string]types.ServerDat
	db   *db.DB
	lock *sync.RWMutex
	dg   *discordgo.Session
}

func NewBase(db *db.DB, dat map[string]types.ServerDat, dg *discordgo.Session, lock *sync.RWMutex) *Base {
	return &Base{
		db:   db,
		dat:  dat,
		lock: lock,
		dg:   dg,
	}
}
