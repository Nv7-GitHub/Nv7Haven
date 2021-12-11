package elements

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

func (b *Elements) SearchCmd(search string, sort string, source string, opt string, regex bool, m types.Msg, rsp types.Rsp) {
	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		return
	}
	rsp.Acknowledge()
	_, res := dat.GetInv(m.Author.ID, true)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	var list map[string]types.Empty
	switch source {
	case "elements":
		list = make(map[string]types.Empty, len(dat.Elements))
		for _, el := range dat.Elements {
			list[el.Name] = types.Empty{}
		}

	case "inventory":
		inv, res := dat.GetInv(opt, m.Author.ID == opt)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}

		list = make(map[string]types.Empty, len(inv.Elements))
		dat.Lock.RLock()
		for el := range inv.Elements {
			elem, res := dat.GetElement(el, true)
			if !res.Exists {
				list[el] = types.Empty{}
				continue
			}
			list[elem.Name] = types.Empty{}
		}
		dat.Lock.RUnlock()

	case "category":
		cat, res := dat.GetCategory(opt)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
		list = cat.Elements
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
	i := 0
	for k := range items {
		txt[i] = k
		i++
	}
	util.SortElemList(txt, sort, dat)

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
