package eod

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
)

var seps = []string{
	"\n",
	"+",
	",",
	"plus",
}

func (b *Bot) PingCmd(c sevcord.Ctx, opts []any) {
	t1 := time.Now()
	c.Acknowledge()
	t2 := time.Now()
	ping := t2.Sub(t1)
	milliseconds := float64(ping) / 1000000
	if milliseconds > 1000 {
		seconds := milliseconds / 1000
		c.Respond(sevcord.NewMessage("🏓 Pong! Latency: **" + humanize.Ftoa(math.Floor(seconds*100)/100) + "s**"))
	} else {
		c.Respond(sevcord.NewMessage("🏓 Pong! Latency: **" + humanize.Ftoa(math.Floor(milliseconds*100)/100) + "ms**"))
	}
}

func (b *Bot) textCommandHandler(c sevcord.Ctx, name string, content string) {
	switch name {
	case "s", "suggest":
		if !b.base.CheckCtx(c, "suggest") {
			return
		}
		val := content
		b.MsgSugElement(c, val)

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
	case "hq", "qh", "hintquery":
		if !b.base.CheckCtx(c, "hint") {
			return
		}
		b.elements.Hint(c, []any{nil, content})
	case "ic", "ci", "infocategory", "infocat", "catinfo", "categoryinfo":
		if !b.base.CheckCtx(c, "info") {
			return
		}
		b.categories.Info(c, []any{content})
	case "iq", "qi", "infoquery", "queryinfo":
		if !b.base.CheckCtx(c, "info") {
			return
		}
		b.queries.Info(c, []any{content})
	case "cat", "c", "category":

		if !b.base.CheckCtx(c, "cat") {
			return
		}
		if content != "" {
			parts := strings.SplitN(content, "|", 2)
			if len(parts) == 2 {
				sort := getSort(strings.ToLower(strings.TrimSpace(parts[1])))
				b.pages.Cat(c, []any{any(strings.TrimSpace(parts[0])), sort})
			} else {
				b.pages.Cat(c, []any{any(content), nil})
			}

		} else {
			b.pages.CatList(c, []any{"name"})
		}
	case "ping":
		b.PingCmd(c, []any{nil})
	case "stats":
		if !b.base.CheckCtx(c, "stats") {
			return
		}
		b.base.Stats(c, []any{nil})
	case "commandlb", "clb":
		if !b.base.CheckCtx(c, "commandlb") {
			return
		}
		b.pages.CommandLb(c, []any{nil})
	case "inv", "inventory":
		if !b.base.CheckCtx(c, "inv") {
			return
		}

		var id string
		if strings.HasPrefix(content, "<@") && strings.HasSuffix(content, ">") {
			id = content[2 : len(content)-1]
		} else {
			id = content
			if content == "" {
				b.pages.Inv(c, []any{nil, nil})
				return
			}
		}
		_, err := strconv.Atoi(id)
		if err != nil {
			c.Respond(sevcord.NewMessage("Invalid user!"))
			return
		}
		user, err := c.Dg().User(id)
		if err != nil {
			b.base.Error(c, err)
			return
		}
		b.pages.Inv(c, []any{any(user), nil})

	case "lb", "leaderboard":
		if !b.base.CheckCtx(c, "lb") {
			return
		}
		parts := strings.Split(content, " ")
		Lbsort := getLbSort(strings.ToLower(parts[0]))
		b.pages.Lb(c, []any{Lbsort, nil, nil})

	case "p", "products", "ih", "inversehint", "invhint":
		if !b.base.CheckCtx(c, "products") {
			return
		}
		id, ok := b.getElementId(c, content)
		if !ok {
			return
		}
		b.pages.Products(c, []any{any(id), nil})

	case "q", "query":
		parts := strings.SplitN(content, "|", 2)

		if !b.base.CheckCtx(c, "query") {
			return
		}
		if content != "" {
			if len(parts) == 2 {

				sort := getSort(strings.ToLower(strings.TrimSpace(parts[1])))
				b.pages.Query(c, []any{any(strings.TrimSpace(parts[0])), sort})

			} else {
				b.pages.Query(c, []any{any(parts[0]), nil})
			}

		} else {
			b.pages.QueryList(c, []any{"name"})
		}

	case "ac", "rc":
		if !b.base.CheckCtx(c, "cat") {
			return
		}
		parts := strings.SplitN(content, "|", 2)
		if len(parts) != 2 {
			c.Respond(sevcord.NewMessage("Invalid format! " + types.RedCircle))
			return
		}
		var inputs []string

		els := make([]int, 0)

		for sep := range seps {
			if strings.Contains(parts[1], seps[sep]) {
				vals := strings.Split(parts[1], seps[sep])
				inputs = append(inputs, vals...)
				break
			}
		}
		if len(inputs) == 0 {
			inputs = append(inputs, parts[1])
		}

		ids, ok := b.getElementIds(c, inputs)
		for i := 0; i < len(ids); i++ {
			if !slices.Contains(els, int(ids[i])) {
				els = append(els, int(ids[i]))
			}

		}
		if !ok {
			return
		}

		if name == "ac" {
			b.categories.CatEditCmd(c, strings.TrimSpace(parts[0]), els, types.PollKindCategorize, "Suggested to add **%s** to **%s** 🗃️", false)
		} else {
			b.categories.CatEditCmd(c, strings.TrimSpace(parts[0]), els, types.PollKindUncategorize, "Suggested to remove **%s** from **%s** 🗃️", true)
		}
	case "dc":
		if !b.base.CheckCtx(c, "cat") {
			return
		}
		parts := strings.SplitN(content, "|", 2)
		if len(parts) != 1 {
			c.Respond(sevcord.NewMessage("Invalid format! " + types.RedCircle))
			return
		}
		b.categories.DelCat(c, []any{parts[0]})
	case "sign", "mark":
		if !b.base.CheckCtx(c, "sign") {
			return
		}
		parts := strings.SplitN(content, "|", 3)
		if len(parts) != 3 {
			c.Respond(sevcord.NewMessage("Use `!sign [e/c/q]|[element/category/query name]|<text>`! " + types.RedCircle))
			return
		}
		// check for signing element/category/query
		switch strings.ToLower(strings.TrimSpace(parts[0])) {
		case "e", "element":
			b.elements.MsgSignCmd(c, strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2]))

		case "c", "cat", "category":
			b.categories.MsgSignCmd(c, strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2]))

		case "q", "query":
			b.queries.MsgSignCmd(c, strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2]))

		default:
			c.Respond(sevcord.NewMessage("Use `!sign [e/c/q]|[element/category/query name]|<text>`! " + types.RedCircle))
		}
	case "col", "color", "colour":
		if !b.base.CheckCtx(c, "color") {
			return
		}
		// check part amount
		parts := strings.SplitN(content, "|", 3)
		if len(parts) != 3 {
			c.Respond(sevcord.NewMessage("Use `!color [e/c/q]|[element/category/query name]|<hex code>`! " + types.RedCircle))
			return
		}
		// check for coloring element/category/query
		switch strings.ToLower(strings.TrimSpace(parts[0])) {
		case "e", "element":
			id, ok := b.getElementId(c, parts[1])
			if !ok {
				return
			}
			b.elements.ColorCmd(c, []any{id, strings.TrimSpace(parts[2])})

		case "c", "cat", "category":
			b.categories.ColorCmd(c, []any{strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2])})

		case "q", "query":
			b.queries.ColorCmd(c, []any{strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2])})

		default:
			c.Respond(sevcord.NewMessage("Use `!color [e/c/q]|[element/category/query name]|<hex code>`! " + types.RedCircle))
		}
	case "n", "next":
		if !b.base.CheckCtx(c, "next") {
			return
		}
		val := any(nil)
		part := strings.TrimSpace(content)
		if part != "" {
			val = any(part)
		}
		b.elements.Next(c, []any{val})
	case "elemcats":
		id, ok := b.getElementId(c, content)
		if ok {
			b.pages.ElemCats(c, []any{any(id)})
		}

	case "img", "image":
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
		switch strings.ToLower(parts[0]) {
		case "e", "element":
			// Get ID
			var id int
			err := b.db.QueryRow("SELECT id FROM elements WHERE LOWER(name)=$1 AND guild=$2", strings.ToLower(parts[1]), c.Guild()).Scan(&id)
			if err != nil {
				b.base.Error(c, err, "Element **"+parts[1]+"** doesn't exist!")
				return
			}

			// Command
			b.elements.ImageCmd(c, id, image)

		case "c", "cat", "category":
			b.categories.ImageCmd(c, parts[1], image)

		case "q", "query":
			b.queries.ImageCmd(c, parts[1], image)

		default:
			c.Respond(sevcord.NewMessage("Use `!image [element/category/query] <element/category/query name>`! " + types.RedCircle))
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
		val := strings.TrimSpace(content[1:])
		b.MsgSugElement(c, val)
		return
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
			c.Respond(ok.Response())
			return
		}
		name, err := b.base.GetName(c.Guild(), comb.Result)
		if err != nil {
			b.base.Error(c, err)
			return
		}

		// Get parts
		parts := []string{content[1:]}
		for _, sep := range seps {
			if strings.Contains(content[1:], sep) {
				parts = strings.Split(content[1:], sep)
				break
			}
		}

		b.combineElements(c, append([]string{name}, parts...))
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
		id, ok := b.getElementId(c, strings.TrimSpace(content[1:]))
		if ok {
			b.elements.Info(c, int(id))
		}

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
		if cnt > types.MaxComboLength {
			c.Respond(sevcord.NewMessage(fmt.Sprintf("You can only combine up to %d elements! "+types.RedCircle, types.MaxComboLength)))
			return
		}
		if cnt < 2 {
			c.Respond(sevcord.NewMessage("You need to combine at least 2 elements! " + types.RedCircle))
			return
		}
		if len(parts) == 2 {
			inps := make([]string, 0, cnt)
			for i := 0; i < cnt; i++ {
				inps = append(inps, strings.TrimSpace(parts[1]))
			}
			b.combineElements(c, inps)
			return
		} else {
			// Get prev
			comb, ok := b.base.GetCombCache(c)
			if !ok.Ok {
				c.Respond(ok.Response())
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
			b.combineElements(c, new)
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
			b.combineElements(c, elems)
			return
		}
	}
}
