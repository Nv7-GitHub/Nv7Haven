package eod

import (
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/elements"
)

func (b *Bot) Init() {
	b.base = base.NewBase(b.s, b.db)
	b.elements = elements.NewElements(b.s, b.db)
}
