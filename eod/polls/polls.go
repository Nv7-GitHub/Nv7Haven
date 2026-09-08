package polls

import (
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
)

type Polls struct {
	db   *sqlx.DB
	base *base.Base
	s    *sevcord.Sevcord

	lock    *sync.RWMutex
	triaged map[string]bool // map[guild]triaged
}

func (p *Polls) Init() {
	p.s.Dg().AddHandler(p.reactionHandler)
	p.s.Dg().AddHandler(p.unReactionHandler)
	p.s.Dg().Identify.Intents |= discordgo.IntentsGuildMessageReactions
}

func NewPolls(d *sqlx.DB, b *base.Base, s *sevcord.Sevcord) *Polls {
	p := &Polls{
		db:   d,
		base: b,
		s:    s,

		lock:    &sync.RWMutex{},
		triaged: make(map[string]bool),
	}
	p.Init()
	return p
}
