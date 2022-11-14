package categories

import (
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/jmoiron/sqlx"
)

type Categories struct {
	db   *sqlx.DB
	base *base.Base
}

func NewCategories(db *sqlx.DB, base *base.Base, s *sevcord.Sevcord) *Categories {
	c := &Categories{
		db:   db,
		base: base,
	}
	return c
}
