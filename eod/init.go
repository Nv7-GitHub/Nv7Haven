package eod

import (
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/elements"
	"github.com/Nv7-Github/Nv7Haven/eod/polls"
)

func (b *Bot) Init() {
	b.base = base.NewBase(b.s, b.db)
	b.polls = polls.NewPolls(b.db, b.base, b.s)
	b.elements = elements.NewElements(b.s, b.db, b.base, b.polls)
	b.s.SetMessageHandler(b.messageHandler)
}
