package categories

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

type catSortInfo struct {
	Name string
	Cnt  int
}

func (b *Categories) CatCmd(category string, sortKind string, hasUser bool, user string, m types.Msg, rsp types.Rsp) {
	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		return
	}

	category = strings.TrimSpace(category)

	if base.IsFoolsMode && !base.IsFool(category) {
		rsp.ErrorMessage(base.MakeFoolResp(category))
		return
	}

	id := m.Author.ID
	if hasUser {
		id = user
	}
	inv, res := dat.GetInv(id, !hasUser)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	cat, res := dat.GetCategory(category)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
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
			text = name + " " + types.Check
			found++
			fnd = 1
		} else {
			text = name + " " + types.X
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

	var o []string
	switch sortKind {
	case "catfound":
		sort.Slice(out, func(i, j int) bool {
			return out[i].found > out[j].found
		})

	case "catnotfound":
		sort.Slice(out, func(i, j int) bool {
			return out[i].found < out[j].found
		})

	case "catelemcount":
		rsp.ErrorMessage("Invalid sort!")
		return

	default:
		util.SortElemObj(out, len(out), func(index int, sort bool) string {
			if sort {
				return out[index].name
			}
			return out[index].text
		}, func(index int, val string) {
			out[index].text = val
		}, sortKind, dat)
	}

	o = make([]string, len(out))
	for i, val := range out {
		o[i] = val.text
	}

	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Thumbnail:  cat.Image,
		Title:      fmt.Sprintf("%s (%d, %s%%)", category, len(out), util.FormatFloat(float32(found)/float32(len(out))*100, 2)),
		PageGetter: b.base.InvPageGetter,
		Items:      o,
		Color:      cat.Color,
	}, m, rsp)
}

type catData struct {
	text  string
	name  string
	found float32
	count int
}

func (b *Categories) AllCatCmd(sortBy string, hasUser bool, user string, m types.Msg, rsp types.Rsp) {
	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		return
	}

	id := m.Author.ID
	if hasUser {
		id = user
	}
	inv, res := dat.GetInv(id, !hasUser)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	dat.Lock.RLock()
	out := make([]catData, len(dat.Categories))

	i := 0
	for _, cat := range dat.Categories {
		count := 0
		for elem := range cat.Elements {
			_, exists := inv[strings.ToLower(elem)]
			if exists {
				count++
			}
		}

		perc := float32(count) / float32(len(cat.Elements))
		text := "(" + util.FormatFloat(perc*100, 2) + "%)"
		if count == len(cat.Elements) {
			text = types.Check
		}
		out[i] = catData{
			text:  fmt.Sprintf("%s %s", cat.Name, text),
			name:  cat.Name,
			found: perc,
			count: len(cat.Elements),
		}
		i++
	}
	dat.Lock.RUnlock()

	switch sortBy {
	case "catfound":
		sort.Slice(out, func(i, j int) bool {
			return out[i].found > out[j].found
		})

	case "catnotfound":
		sort.Slice(out, func(i, j int) bool {
			return out[i].found < out[j].found
		})

	case "catelemcount":
		sort.Slice(out, func(i, j int) bool {
			return out[i].count > out[j].count
		})

	default:
		sort.Slice(out, func(i, j int) bool {
			return util.CompareStrings(out[i].name, out[j].name)
		})
	}

	names := make([]string, len(out))
	for i, dat := range out {
		names[i] = dat.text
	}

	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      fmt.Sprintf("All Categories (%d)", len(out)),
		PageGetter: b.base.InvPageGetter,
		Items:      names,
	}, m, rsp)
}
