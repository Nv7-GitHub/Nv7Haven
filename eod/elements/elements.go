package elements

import (
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/jmoiron/sqlx"
)

type Elements struct {
	s    *sevcord.Sevcord
	db   *sqlx.DB
	base *base.Base
}

func NewElements(s *sevcord.Sevcord, db *sqlx.DB, base *base.Base) *Elements {
	e := &Elements{
		s:    s,
		db:   db,
		base: base,
	}
	s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"info",
		"Get element info!",
		e.Info,
		sevcord.NewOption("element", "The ID of the element to view the info of!", sevcord.OptionKindInt, true).
			AutoComplete(e.Autocomplete),
	))
	return e
}
