package pages

import (
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/jmoiron/sqlx"
)

type Pages struct {
	base *base.Base
	db   *sqlx.DB
}

func NewPages(base *base.Base, db *sqlx.DB, s *sevcord.Sevcord) *Pages {
	p := &Pages{
		base: base,
		db:   db,
	}
	s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"inv",
		"View your inventory!",
		p.Inv,
		sevcord.NewOption("user", "The user to view the inventory of!", sevcord.OptionKindUser, false),
		sevcord.NewOption("sort", "The sort order of the inventory!", sevcord.OptionKindString, false).
			AddChoices(types.Sorts...),
	))
	s.AddButtonHandler("inv", p.InvHandler)
	return p
}
