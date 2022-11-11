package elements

import (
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/jmoiron/sqlx"
)

type Elements struct {
	s  *sevcord.Sevcord
	db *sqlx.DB
}

func NewElements(s *sevcord.Sevcord, db *sqlx.DB) *Elements {
	s.RegisterSlashCommand(sevcord.NewSlashCommand("ping", "Ping!", func(c sevcord.Ctx, opts []any) {
		c.Respond(sevcord.NewMessage("Pong!"))
	}))
	return &Elements{
		s:  s,
		db: db,
	}
}
