package base

import (
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/jmoiron/sqlx"
)

const configCmdId = ""

type Base struct {
	s  *sevcord.Sevcord
	db *sqlx.DB
}

func (b *Base) Init() {
	b.s.AddMiddleware(b.CheckCtx)
}

func NewBase(s *sevcord.Sevcord, db *sqlx.DB) *Base {
	b := &Base{
		s:  s,
		db: db,
	}
	b.Init()
	return b
}
