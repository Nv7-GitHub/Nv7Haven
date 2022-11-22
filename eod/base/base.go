package base

import (
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/jmoiron/sqlx"
)

const configCmdId = "7" // TODO: Put in real value

type Base struct {
	s  *sevcord.Sevcord
	db *sqlx.DB

	lock *sync.RWMutex
	mem  map[string]*types.ServerMem // map[guild]data
}

func (b *Base) Init() {
	b.s.AddMiddleware(b.CheckCtx)
	b.s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"stats",
		"View the statistics of this server!",
		b.Stats,
	))
}

func NewBase(s *sevcord.Sevcord, db *sqlx.DB) *Base {
	b := &Base{
		lock: &sync.RWMutex{},
		mem:  make(map[string]*types.ServerMem),
		s:    s,
		db:   db,
	}
	b.Init()
	return b
}
