package eod

import (
	"strings"

	"github.com/Nv7-Github/sevcord/v2"
)

var seps = []string{
	"\n",
	"+",
	",",
	"plus",
}

func (b *Bot) messageHandler(c sevcord.Ctx, content string) {
	for _, sep := range seps {
		if strings.Contains(content, sep) {
			// Check ctx
			if !b.base.CheckCtx(c) {
				return
			}
			if !b.base.IsPlayChannel(c) {
				return
			}

			// Combine
			elems := strings.Split(content, sep)
			b.elements.Combine(c, elems)
		}
	}
}
