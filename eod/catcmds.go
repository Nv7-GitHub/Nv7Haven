package eod

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

const x = "❌"
const check = "✅"

const (
	catSortAlphabetical   = 0
	catSortByFound        = 1
	catSortByNotFound     = 2
	catSortByElementCount = 3
)

func (b *EoD) catCmd(category string, sortKind int, hasUser bool, user string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	if isFoolsMode && !isFool(category) {
		rsp.ErrorMessage(makeFoolResp(category))
		return
	}

	msg := "You don't have an inventory!"
	id := m.Author.ID
	if hasUser {
		id = user
		msg = fmt.Sprintf("User <@%s> doesn't have an inventory!", user)
	}
	dat.Lock.RLock()
	inv, exists := dat.InvCache[id]
	dat.Lock.RUnlock()
	if !exists {
		rsp.ErrorMessage(msg)
		return
	}

	cat, exists := dat.CatCache[strings.ToLower(category)]
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Category **%s** doesn't exist!", category))
		return
	}
	category = cat.Name

	out := make([]struct {
		found int
		text  string
		name  string
	}, len(cat.Elements))

	found := 0
	i := 0
	fnd := 0
	var text string

	for name := range cat.Elements {
		_, exists := inv[strings.ToLower(name)]
		if exists {
			text = name + " " + check
			found++
			fnd = 1
		} else {
			text = name + " " + x
			fnd = 0
		}

		out[i] = struct {
			found int
			text  string
			name  string
		}{
			found: fnd,
			text:  text,
			name:  name,
		}

		i++
	}

	switch sortKind {
	case catSortByFound:
		sort.Slice(out, func(i, j int) bool {
			return out[i].found > out[j].found
		})

	case catSortByNotFound:
		sort.Slice(out, func(i, j int) bool {
			return out[i].found < out[j].found
		})

	default:
		sort.Slice(out, func(i, j int) bool {
			return compareStrings(out[i].name, out[j].name)
		})
	}

	o := make([]string, len(out))
	for i, val := range out {
		o[i] = val.text
	}

	b.newPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Thumbnail:  cat.Image,
		Title:      fmt.Sprintf("%s (%d, %s%%)", category, len(out), formatFloat(float32(found)/float32(len(out))*100, 2)),
		PageGetter: b.invPageGetter,
		Items:      o,
	}, m, rsp)
}

type catData struct {
	text  string
	name  string
	found float32
	count int
}

func (b *EoD) allCatCmd(sortBy int, hasUser bool, user string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	msg := "You don't have an inventory!"
	id := m.Author.ID
	if hasUser {
		id = user
		msg = fmt.Sprintf("User <@%s> doesn't have an inventory!", user)
	}
	dat.Lock.RLock()
	inv, exists := dat.InvCache[id]
	dat.Lock.RUnlock()
	if !exists {
		rsp.ErrorMessage(msg)
		return
	}

	out := make([]catData, len(dat.CatCache))

	i := 0
	for _, cat := range dat.CatCache {
		count := 0
		for elem := range cat.Elements {
			_, exists := inv[strings.ToLower(elem)]
			if exists {
				count++
			}
		}

		perc := float32(count) / float32(len(cat.Elements))
		text := "(" + formatFloat(perc*100, 2) + "%)"
		if count == len(cat.Elements) {
			text = check
		}
		out[i] = catData{
			text:  fmt.Sprintf("%s %s", cat.Name, text),
			name:  cat.Name,
			found: perc,
			count: len(cat.Elements),
		}
		i++
	}

	switch sortBy {
	case catSortByFound:
		sort.Slice(out, func(i, j int) bool {
			return out[i].found > out[j].found
		})

	case catSortByNotFound:
		sort.Slice(out, func(i, j int) bool {
			return out[i].found < out[j].found
		})

	case catSortAlphabetical:
		sort.Slice(out, func(i, j int) bool {
			return compareStrings(out[i].name, out[j].name)
		})

	case catSortByElementCount:
		sort.Slice(out, func(i, j int) bool {
			return out[i].count > out[j].count
		})
	}

	names := make([]string, len(out))
	for i, dat := range out {
		names[i] = dat.text
	}

	b.newPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      fmt.Sprintf("All Categories (%d)", len(out)),
		PageGetter: b.invPageGetter,
		Items:      names,
	}, m, rsp)
}
