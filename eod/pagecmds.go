package eod

import (
	"fmt"
	"sort"
	"strings"

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
	dat.Lock.RLock()
	inv, exists := dat.InvCache[user]
	dat.Lock.RUnlock()
	if !exists {
		if user == m.Author.ID {
			rsp.ErrorMessage("You don't have an inventory!")
		} else {
			rsp.ErrorMessage(fmt.Sprintf("User <@%s> doesn't have an inventory!", user))
		}
		return
	}
	items := make([]string, len(inv))
	i := 0
	dat.Lock.RLock()
	for k := range inv {
		items[i] = dat.ElemCache[k].Name
		i++
	}

	switch sorter {
	case "id":
		sort.Slice(items, func(i, j int) bool {
			elem1, exists := dat.ElemCache[strings.ToLower(items[i])]
			if !exists {
				return false
			}

			elem2, exists := dat.ElemCache[strings.ToLower(items[j])]
			if !exists {
				return false
			}
			return elem1.CreatedOn.Before(elem2.CreatedOn)
		})

	case "madeby":
		count := 0
		outs := make([]string, len(items))
		for _, val := range items {
			creator := ""
			elem, exists := dat.ElemCache[strings.ToLower(val)]
			if exists {
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
		Title:      fmt.Sprintf("%s's Inventory (%d, %s%%)", name, len(items), formatFloat(float32(len(items))/float32(len(dat.ElemCache))*100, 2)),
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
	dat.Lock.RLock()
	_, exists = dat.InvCache[m.Author.ID]
	dat.Lock.RUnlock()
	if !exists {
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
