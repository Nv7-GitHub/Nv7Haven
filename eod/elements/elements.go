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
}

func NewElements(s *sevcord.Sevcord, db *sqlx.DB, base *base.Base, polls *polls.Polls) *Elements {
	e := &Elements{
		db:    db,
		base:  base,
		polls: polls,
	}
	s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"info",
		"Get element info!",
		e.Info,
		sevcord.NewOption("element", "The ID of the element to view the info of!", sevcord.OptionKindInt, true).
			AutoComplete(e.Autocomplete),
	))
	s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"hint",
		"Learn how to make an element!",
		e.Hint,
		sevcord.NewOption("element", "An element to get the hint of!", sevcord.OptionKindInt, false).
			AutoComplete(e.Autocomplete),
	))
	s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"suggest",
		"Create a suggestion!",
		e.Suggest,
		sevcord.NewOption("result", "The result of the combination!", sevcord.OptionKindString, true),
		sevcord.NewOption("autocapitalize", "Whether or not to autocapitalize!", sevcord.OptionKindBool, false),
	))
	return e
}
