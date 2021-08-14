package eod

import (
	"fmt"
	"sort"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *EoD) invCmd(user string, m types.Msg, rsp types.Rsp, sorter string) {
	rsp.Acknowledge()

	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
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

	switch sorter {
	case "id":
		dat.Lock.RLock()
		sort.Slice(items, func(i, j int) bool {
			elem1, res := dat.GetElement(items[i], true)
			if !res.Exists {
				return false
			}

			elem2, res := dat.GetElement(items[j])
			if !res.Exists {
				return false
			}
			return elem1.CreatedOn.Before(elem2.CreatedOn)
		})
		dat.Lock.RUnlock()

	case "madeby":
		count := 0
		outs := make([]string, len(items))
		dat.Lock.RLock()
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
		sortStrings(outs)
		items = outs

	case "length":
		sort.Slice(items, func(i, j int) bool {
			return len(items[i]) < len(items[j])
		})

	default:
		sortStrings(items)
	}
	dat.Lock.RUnlock()

	name := m.Author.Username
	if m.Author.ID != user {
		u, err := b.dg.User(user)
		if rsp.Error(err) {
			return
		}
		name = u.Username
	}
	b.newPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      fmt.Sprintf("%s's Inventory (%d, %s%%)", name, len(items), formatFloat(float32(len(items))/float32(len(dat.Elements))*100, 2)),
		PageGetter: b.invPageGetter,
		Items:      items,
	}, m, rsp)
}

func (b *EoD) lbCmd(m types.Msg, rsp types.Rsp, sort string) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	_, res := dat.GetInv(m.Author.ID, true)
	if !res.Exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}

	b.newPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchLdb,
		Title:      "Top Most Elements",
		PageGetter: b.lbPageGetter,
		Sort:       sort,
		User:       m.Author.ID,
	}, m, rsp)
}
