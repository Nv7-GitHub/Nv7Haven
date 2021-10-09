package polls

import (
	"database/sql"
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

type Polls struct {
	dat  map[string]types.ServerData
	lock *sync.RWMutex
	dg   *discordgo.Session
	db   *sql.DB
	base *base.Base
}

func NewPolls(dat map[string]types.ServerData, dg *discordgo.Session, db *sql.DB, base *base.Base, lock *sync.RWMutex) *Polls {
	return &Polls{
		dat:  dat,
		lock: lock,
		dg:   dg,
		db:   db,
		base: base,
	}
}
