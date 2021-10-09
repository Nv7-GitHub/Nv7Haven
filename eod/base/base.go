package base

import (
	"database/sql"
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

type Base struct {
	dat  map[string]types.ServerData
	db   *sql.DB
	lock *sync.RWMutex
	dg   *discordgo.Session
}

func NewBase(db *sql.DB, dat map[string]types.ServerData, dg *discordgo.Session, lock *sync.RWMutex) *Base {
	return &Base{
		db:   db,
		dat:  dat,
		lock: lock,
		dg:   dg,
	}
}
