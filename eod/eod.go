package eod

import (
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/elements"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/jmoiron/sqlx"
)

type Bot struct {
	s  *sevcord.Sevcord
	db *sqlx.DB

	// Modules
	base     *base.Base
	elements *elements.Elements
}

func InitEod(db *sqlx.DB, token string) {
	s, err := sevcord.New(token)
	if err != nil {
		panic(err)
	}
	b := Bot{
		s:  s,
		db: db,
	}
	b.Init()
	b.s.Listen()
}
