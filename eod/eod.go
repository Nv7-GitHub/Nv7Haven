package eod

import (
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/jmoiron/sqlx"
)

type Bot struct {
	s  *sevcord.Sevcord
	db *sqlx.DB
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
