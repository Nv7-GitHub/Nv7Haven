package eod

import (
	"fmt"
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
			// Check if play channel
			var cnt bool
			err := b.db.QueryRow(`SELECT $1=ANY(play) FROM config WHERE guild=$2`, c.Channel(), c.Guild()).Scan(&cnt)
			if err != nil {
				fmt.Println(err)
				return
			}
			if !cnt {
				return
			}

			// Combine
			elems := strings.Split(content, sep)
			b.elements.Combine(c, elems)
		}
	}
}
