package elements

import (
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/polls"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/jmoiron/sqlx"
)

type Elements struct {
	db    *sqlx.DB
	base  *base.Base
	polls *polls.Polls
	s     *sevcord.Sevcord
}

func (e *Elements) Init() {
	e.s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"suggest",
		"Create a suggestion!",
		e.Suggest,
		sevcord.NewOption("result", "The result of the combination!", sevcord.OptionKindString, true),
		sevcord.NewOption("autocapitalize", "Whether or not to autocapitalize!", sevcord.OptionKindBool, false),
	))
}

func NewElements(s *sevcord.Sevcord, db *sqlx.DB, base *base.Base, polls *polls.Polls) *Elements {
	e := &Elements{
		db:    db,
		base:  base,
		polls: polls,
		s:     s,
	}
	e.Init()
	return e
}
