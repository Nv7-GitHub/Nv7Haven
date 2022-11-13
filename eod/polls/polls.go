package polls

import (
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
)

type Polls struct {
	db   *sqlx.DB
	base *base.Base
}

func NewPolls(d *sqlx.DB, b *base.Base, s *sevcord.Sevcord) *Polls {
	p := &Polls{
		db:   d,
		base: b,
	}
	s.Dg().AddHandler(p.reactionHandler)
	s.Dg().AddHandler(p.unReactionHandler)
	s.Dg().Identify.Intents |= discordgo.IntentsGuildMessageReactions
	return p
}
