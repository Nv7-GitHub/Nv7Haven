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

func (b *Bot) getElementId(c sevcord.Ctx, val string) (int64, bool) {
	var id int64
	err := b.db.QueryRow("SELECT id FROM elements WHERE LOWER(name)=$1 AND guild=$2", strings.ToLower(strings.TrimSpace(val)), c.Guild()).Scan(&id)
	if err != nil {
		b.base.Error(c, err)
		return 0, false
	}
	return id, true
}

func (b *Bot) textCommandHandler(c sevcord.Ctx, name string, content string) {
	switch name {
	case "s", "suggest":
		b.base.IncrementCommandStat(c, "suggest")
		b.elements.Suggest(c, []any{any(content), nil})

	case "h", "hint":
		b.base.IncrementCommandStat(c, "hint")
		val := any(nil)
		if content != "" {
			v, ok := b.getElementId(c, content)
			if !ok {
				return
			}
			val = any(v)
		}
		b.elements.Hint(c, []any{val, nil})

	case "cat":
		b.base.IncrementCommandStat(c, "cat")
		b.pages.Cat(c, []any{any(content), nil})

	case "p", "products":
		b.base.IncrementCommandStat(c, "products")
		id, ok := b.getElementId(c, content)
		if !ok {
			return
		}
		b.elements.Products(c, []any{any(id), nil})
	}
}

func (b *Bot) messageHandler(c sevcord.Ctx, content string) {
	if strings.HasPrefix(content, "!") {
		parts := strings.SplitN(content[1:], " ", 2)
		b.textCommandHandler(c, parts[0], parts[1])
	}
	if strings.HasPrefix(content, "?") {
		b.elements.InfoMsgCmd(c, content[1:])
		return
	}
	for _, sep := range seps {
		if strings.Contains(content, sep) {
			// Check ctx
			if !b.base.CheckCtx(c, "message") {
				return
			}
			if !b.base.IsPlayChannel(c) {
				return
			}

			// Combine
			elems := strings.Split(content, sep)
			b.elements.Combine(c, elems)
			return
		}
	}
}
