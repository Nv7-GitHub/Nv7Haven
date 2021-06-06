package eod

import (
	"fmt"
	"sort"
	"strings"
)

const x = "❌"
const check = "✅"

const (
	catSortAlphabetical = 0
	catSortByFound      = 1
	catSortByNotFound   = 2
)

func (b *EoD) catCmd(category string, sortKind int, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	inv, exists := dat.invCache[m.Author.ID]
	if !exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}

	cat, exists := dat.catCache[strings.ToLower(category)]
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
	case catSortAlphabetical:
		sort.Slice(out, func(i, j int) bool {
			return out[i].name < out[j].name
		})

	case catSortByFound:
		sort.Slice(out, func(i, j int) bool {
			return out[i].found > out[j].found
		})

	case catSortByNotFound:
		sort.Slice(out, func(i, j int) bool {
			return out[i].found < out[j].found
		})
	}

	o := make([]string, len(out))
	for i, val := range out {
		o[i] = val.text
	}

	b.newPageSwitcher(pageSwitcher{
		Kind:       pageSwitchInv,
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
}

func (b *EoD) allCatCmd(sortBy int, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	inv, exists := dat.invCache[m.Author.ID]
	if !exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}

	out := make([]catData, len(dat.catCache))

	i := 0
	for _, cat := range dat.catCache {
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
			return out[i].name < out[j].name
		})
	}

	names := make([]string, len(out))
	for i, dat := range out {
		names[i] = dat.text
	}

	b.newPageSwitcher(pageSwitcher{
		Kind:       pageSwitchInv,
		Title:      fmt.Sprintf("All Categories (%d)", len(out)),
		PageGetter: b.invPageGetter,
		Items:      names,
	}, m, rsp)
}
