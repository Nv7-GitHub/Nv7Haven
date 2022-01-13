package categories

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/eodsort"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

type catSortInfo struct {
	Name string
	Cnt  int
}

func (b *Categories) CatCmd(category string, sortKind string, hasUser bool, user string, postfix bool, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
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
	inv := db.GetInv(id)

	cat, res := db.GetCat(category)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	category = cat.Name

	out := make([]struct {
		text string
		id   int
		name string
	}, len(cat.Elements))

	found := 0
	i := 0
	var text string

	db.RLock()
	for elem := range cat.Elements {
		exists := inv.Contains(elem)
		el, _ := db.GetElement(elem, true)
		if exists {
			text = el.Name + " " + types.Check
			found++
		} else {
			text = el.Name + " " + types.X
		}

		out[i] = struct {
			text string
			id   int
			name string
		}{
			text: text,
			id:   el.ID,
			name: el.Name,
		}

		i++
	}
	db.RUnlock()

	var o []string
	switch sortKind {
	case "catelemcount":
		rsp.ErrorMessage("Invalid sort!")
		return

	default:
		if sortKind == "found" {
			postfix = false
		}
		eodsort.Sort(out, len(out), func(index int) int {
			return out[index].id
		}, func(index int) string {
			return out[index].text
		}, func(index int, val string) {
			out[index].text = val
		}, sortKind, m.Author.ID, db, postfix)
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
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}

	id := m.Author.ID
	if hasUser {
		id = user
	}
	inv := db.GetInv(id)

	db.RLock()
	out := make([]catData, len(db.Cats()))

	i := 0
	for _, cat := range db.Cats() {
		count := 0
		for elem := range cat.Elements {
			exists := inv.Contains(elem)
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
	db.RUnlock()

	switch sortBy {
	case "found":
		sort.Slice(out, func(i, j int) bool {
			return out[i].found > out[j].found
		})

	case "catelemcount":
		sort.Slice(out, func(i, j int) bool {
			return out[i].count > out[j].count
		})

	default:
		sort.Slice(out, func(i, j int) bool {
			return eodsort.CompareStrings(out[i].name, out[j].name)
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
