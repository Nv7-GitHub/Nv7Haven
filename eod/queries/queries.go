package queries

import (
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/jmoiron/sqlx"
)

type Queries struct {
	db   *sqlx.DB
	base *base.Base
	s    *sevcord.Sevcord
}

func (q *Queries) Init() {

}

func NewQueries(s *sevcord.Sevcord, db *sqlx.DB, base *base.Base) *Queries {
	q := &Queries{
		db:   db,
		base: base,
		s:    s,
	}
	q.Init()
	return q
}
