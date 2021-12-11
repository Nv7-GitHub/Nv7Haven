package elements

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/eodsort"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *Elements) SearchCmd(search string, sort string, source string, opt string, regex bool, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}
	rsp.Acknowledge()

	var list map[string]types.Empty
	switch source {
	case "elements":
		list = make(map[string]types.Empty, len(db.Elements))
		for _, el := range db.Elements {
			list[el.Name] = types.Empty{}
		}

	case "inventory":
		inv := db.GetInv(opt)

		list = make(map[string]types.Empty, len(inv.Elements))
		inv.Lock.RLock()
		db.RLock()
		for el := range inv.Elements {
			elem, res := db.GetElement(el, true)
			if !res.Exists {
				continue
			}
			list[elem.Name] = types.Empty{}
		}
		db.RUnlock()
		inv.Lock.RUnlock()

	case "category":
		cat, res := db.GetCat(opt)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}

		list = make(map[string]types.Empty, len(cat.Elements))
		cat.Lock.RLock()
		db.RLock()
		for el := range cat.Elements {
			elem, res := db.GetElement(el, true)
			if !res.Exists {
				continue
			}
			list[elem.Name] = types.Empty{}
		}
		db.RUnlock()
		cat.Lock.RUnlock()
	}

	items := make(map[string]types.Empty)
	if regex {
		reg, err := regexp.Compile(search)
		if rsp.Error(err) {
			return
		}
		for el := range list {
			m := reg.Find([]byte(el))
			if m != nil {
				items[el] = types.Empty{}
			}
		}
	} else {
		s := strings.ToLower(search)
		for el := range list {
			if strings.Contains(strings.ToLower(el), s) {
				items[el] = types.Empty{}
			}
		}
	}

	txt := make([]string, len(items))
	ids := make([]int, len(items))
	i := 0
	db.RLock()
	for k := range items {
		txt[i] = k
		i++
		el, res := db.GetElementByName(k, true)
		if !res.Exists {
			continue
		}
		ids[i-1] = el.ID
	}
	db.RUnlock()

	eodsort.SortElemObj(txt, len(txt), func(index int) int {
		return ids[index]
	}, func(index int) string {
		return txt[index]
	}, func(index int, val string) {
		txt[index] = val
	}, sort, db)

	if len(txt) == 0 {
		rsp.Message("No results!")
		return
	}

	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      fmt.Sprintf("Element Search (%d)", len(txt)),
		PageGetter: b.base.InvPageGetter,
		Items:      txt,
		User:       m.Author.ID,
	}, m, rsp)
}
