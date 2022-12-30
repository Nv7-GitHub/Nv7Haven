package queries

import (
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/polls"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/jmoiron/sqlx"
)

type Queries struct {
	db    *sqlx.DB
	base  *base.Base
	s     *sevcord.Sevcord
	polls *polls.Polls
}

func (q *Queries) Init() {
	q.s.RegisterSlashCommand(sevcord.NewSlashCommandGroup("newquery", "Create a new query!",
		sevcord.NewSlashCommand(
			"elements",
			"Create a query that contains every element!",
			q.CreateElementsCmd,
			sevcord.NewOption("name", "The name of the query!", sevcord.OptionKindString, true),
		),
	))
}

func NewQueries(s *sevcord.Sevcord, db *sqlx.DB, base *base.Base, polls *polls.Polls) *Queries {
	q := &Queries{
		db:    db,
		base:  base,
		s:     s,
		polls: polls,
	}
	q.Init()
	return q
}
