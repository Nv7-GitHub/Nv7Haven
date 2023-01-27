package eod

import (
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
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
		b.base.Error(c, err, "Element **"+val+"** doesn't exist!")
		return 0, false
	}
	return id, true
}

func (b *Bot) textCommandHandler(c sevcord.Ctx, name string, content string) {
	switch name {
	case "s", "suggest":
		if !b.base.CheckCtx(c, "suggest") {
			return
		}
		b.elements.Suggest(c, []any{any(content), nil})

	case "h", "hint":
		if !b.base.CheckCtx(c, "hint") {
			return
		}
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
		if !b.base.CheckCtx(c, "cat") {
			return
		}
		b.pages.Cat(c, []any{any(content), nil})

	case "p", "products":
		if !b.base.CheckCtx(c, "products") {
			return
		}
		id, ok := b.getElementId(c, content)
		if !ok {
			return
		}
		b.elements.Products(c, []any{any(id), nil})

	case "q", "query":
		if !b.base.CheckCtx(c, "query") {
			return
		}
		b.pages.Query(c, []any{any(content), nil})

	case "ac", "rc":
		if !b.base.CheckCtx(c, "cat") {
			return
		}
		parts := strings.SplitN(content, "|", 2)
		if len(parts) != 2 {
			c.Respond(sevcord.NewMessage("Invalid format!"))
			return
		}
		els := make([]int, 0)
		added := false
		for sep := range seps {
			if strings.Contains(parts[1], seps[sep]) {
				vals := strings.Split(parts[1], seps[sep])
				for _, val := range vals {
					id, ok := b.getElementId(c, val)
					if !ok {
						return
					}
					els = append(els, int(id))
				}
				added = true
				break
			}
		}
		if !added {
			id, ok := b.getElementId(c, parts[1])
			if !ok {
				return
			}
			els = append(els, int(id))
		}
		if len(els) == 0 {
			c.Respond(sevcord.NewMessage("Invalid format!"))
			return
		}
		if name == "ac" {
			b.categories.CatEditCmd(c, parts[0], els, types.PollKindCategorize, "Suggested to add **%s** to **%s** 🗃️", false)
		} else {
			b.categories.CatEditCmd(c, parts[0], els, types.PollKindUncategorize, "Suggested to remove **%s** from **%s** 🗃️", true)
		}

	case "img":
		if !b.base.CheckCtx(c, "image") {
			return
		}

		// Get image
		var image string
		m := c.(*sevcord.MessageCtx).Message()
		if len(m.Attachments) < 1 {
			c.Respond(sevcord.NewMessage("No image attached! " + types.RedCircle))
			return
		}
		if len(m.Attachments) > 1 {
			c.Respond(sevcord.NewMessage("Too many images attached! " + types.RedCircle))
			return
		}
		if !strings.HasPrefix(m.Attachments[0].ContentType, "image/") {
			c.Respond(sevcord.NewMessage("Invalid image format! " + types.RedCircle))
			return
		}
		image = m.Attachments[0].URL

		// Parse
		parts := strings.SplitN(content, " ", 2)
		if len(parts) != 2 {
			c.Respond(sevcord.NewMessage("Use `!image [element/category/query] <element/category/query name>`! " + types.RedCircle))
			return
		}

		// Run command
		switch content {
		case "e", "element":
			// Get ID
			var id int
			err := b.db.QueryRow("SELECT id FROM elements WHERE LOWER(name)=$1 AND guild=$2", strings.ToLower(parts[1])).Scan(&id)
			if err != nil {
				b.base.Error(c, err, "Element **"+content+"** doesn't exist! "+types.RedCircle)
				return
			}

			// Command
			b.elements.ImageCmd(c, id, image)

		case "c", "cat", "category":
			b.categories.ImageCmd(c, parts[1], image)

		case "q", "query":
			b.queries.ImageCmd(c, parts[1], image)
		}
	}
}

func (b *Bot) messageHandler(c sevcord.Ctx, content string) {
	if strings.HasPrefix(content, "=") {
		if len(content) < 2 {
			return
		}
		if !b.base.CheckCtx(c, "suggest") {
			return
		}
		b.elements.Suggest(c, []any{any(content[1:]), nil})
	}
	if strings.HasPrefix(content, "+") {
		if len(content) < 2 {
			return
		}
		if !b.base.CheckCtx(c, "message") {
			return
		}
		if !b.base.IsPlayChannel(c) {
			return
		}

		comb, ok := b.base.GetCombCache(c)
		if !ok.Ok {
			c.Respond(sevcord.NewMessage(ok.Message + " " + types.RedCircle))
			return
		}
		name, err := b.base.GetName(c.Guild(), comb.Result)
		if err != nil {
			b.base.Error(c, err)
			return
		}
		b.elements.Combine(c, []string{name, strings.TrimSpace(content[1:])})
		return
	}
	if strings.HasPrefix(content, "!") {
		if len(content) < 2 {
			return
		}
		parts := strings.SplitN(content[1:], " ", 2)
		if len(parts) < 2 {
			parts = append(parts, "")
		}
		b.textCommandHandler(c, strings.ToLower(parts[0]), parts[1])
		return
	}
	if strings.HasPrefix(content, "?") {
		if len(content) < 2 {
			return
		}
		if !b.base.CheckCtx(c, "info") {
			return
		}
		b.elements.InfoMsgCmd(c, content[1:])
		return
	}
	if strings.HasPrefix(content, "*") {
		if len(content) < 2 {
			return
		}
		if !b.base.CheckCtx(c, "message") {
			return
		}
		if !b.base.IsPlayChannel(c) {
			return
		}

		parts := strings.SplitN(content[1:], " ", 2)
		cnt, err := strconv.Atoi(parts[0])
		if err != nil {
			c.Respond(sevcord.NewMessage("Invalid number of repeats! " + types.RedCircle))
			return
		}
		if cnt < 1 {
			return
		}
		if len(parts) == 2 {
			inps := make([]string, 0, cnt)
			for i := 0; i < cnt; i++ {
				inps = append(inps, strings.TrimSpace(parts[1]))
			}
			b.elements.Combine(c, inps)
			return
		} else {
			// Get prev
			comb, ok := b.base.GetCombCache(c)
			if !ok.Ok {
				c.Respond(sevcord.NewMessage(ok.Message + " " + types.RedCircle))
				return
			}
			if comb.Result == -1 {
				c.Respond(sevcord.NewMessage("You haven't combined anything! " + types.RedCircle))
				return
			}
			name, err := b.base.GetName(c.Guild(), comb.Result)
			if err != nil {
				b.base.Error(c, err)
				return
			}
			new := make([]string, 0, cnt)
			for i := 0; i < cnt; i++ {
				new = append(new, name)
			}
			b.elements.Combine(c, new)
			return
		}
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
