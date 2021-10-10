package elements

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

func (b *Elements) InvCmd(user string, m types.Msg, rsp types.Rsp, sorter string, filter string) {
	rsp.Acknowledge()

	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("Guild not setup!")
		return
	}

	inv, res := dat.GetInv(user, user == m.Author.ID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	items := make([]string, len(inv))
	i := 0
	dat.Lock.RLock()
	for k := range inv {
		el, _ := dat.GetElement(k, true)
		items[i] = el.Name
		i++
	}

	switch filter {
	case "madeby":
		count := 0
		outs := make([]string, len(items))
		for _, val := range items {
			creator := ""
			elem, res := dat.GetElement(val, true)
			if res.Exists {
				creator = elem.Creator
			}
			if creator == user {
				outs[count] = val
				count++
			}
		}
		outs = outs[:count]
		items = outs
	}
	util.SortElemList(items, sorter, dat)
	dat.Lock.RUnlock()

	name := m.Author.Username
	if m.Author.ID != user {
		u, err := b.dg.User(user)
		if rsp.Error(err) {
			return
		}
		name = u.Username
	}
	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      fmt.Sprintf("%s's Inventory (%d, %s%%)", name, len(items), util.FormatFloat(float32(len(items))/float32(len(dat.Elements))*100, 2)),
		PageGetter: b.base.InvPageGetter,
		Items:      items,
	}, m, rsp)
}

func (b *Elements) LbCmd(m types.Msg, rsp types.Rsp, sort string, user string) {
	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		return
	}
	_, res := dat.GetInv(user, user == m.Author.ID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchLdb,
		Title:      "Top Most Elements",
		PageGetter: b.base.LbPageGetter,
		Sort:       sort,
		User:       user,
	}, m, rsp)
}

func (b *Elements) ElemSearchCmd(search string, m types.Msg, rsp types.Rsp) {
	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		return
	}
	_, res := dat.GetInv(m.Author.ID, true)
	if !res.Exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}
	if util.IsWildcard(search) {
		for val := range util.Wildcards {
			search = strings.ReplaceAll(search, string([]rune{val}), string([]rune{'\\', val}))
		}
	}

	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchSearch,
		Title:      "Element Search",
		PageGetter: b.base.SearchPageGetter,
		Search:     search,
		User:       m.Author.ID,
	}, m, rsp)
}
