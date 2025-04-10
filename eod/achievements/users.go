package achievements

import (
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/categories"
	"github.com/Nv7-Github/Nv7Haven/eod/elements"
	"github.com/Nv7-Github/Nv7Haven/eod/queries"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/jmoiron/sqlx"
)

type Users struct {
	base       *base.Base
	db         *sqlx.DB
	categories *categories.Categories
	elements   *elements.Elements
	queries    *queries.Queries
	s          *sevcord.Sevcord
}

func (u *Users) Init() {

	u.s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"profile",
		"View your profile",
		u.Profile,
		sevcord.NewOption("user", "The user to view the profile of!", sevcord.OptionKindUser, false),
	))
}

func NewUsers(base *base.Base, db *sqlx.DB, s *sevcord.Sevcord, categories *categories.Categories, elements *elements.Elements, queries *queries.Queries) *Users {
	u := &Users{
		base:       base,
		db:         db,
		categories: categories,
		elements:   elements,
		queries:    queries,
		s:          s,
	}
	u.Init()
	return u
}
