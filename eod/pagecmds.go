package eod

import (
	"fmt"
	"sort"
	"strings"
)

func (b *EoD) invCmd(user string, m msg, rsp rsp, sorter string) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	dat.lock.RLock()
	inv, exists := dat.invCache[user]
	dat.lock.RUnlock()
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
	for k := range inv {
		dat.lock.RLock()
		items[i] = dat.elemCache[k].Name
		dat.lock.RUnlock()
		i++
	}

	switch sorter {
	case "id":
		sort.Slice(items, func(i, j int) bool {
			dat.lock.RLock()
			elem1, exists := dat.elemCache[strings.ToLower(items[i])]
			if !exists {
				return false
			}

			elem2, exists := dat.elemCache[strings.ToLower(items[j])]
			if !exists {
				return false
			}
			dat.lock.RUnlock()
			return elem1.CreatedOn.Before(elem2.CreatedOn)
		})

	case "madeby":
		count := 0
		outs := make([]string, len(items))
		for _, val := range items {
			creator := ""
			dat.lock.RLock()
			elem, exists := dat.elemCache[strings.ToLower(val)]
			dat.lock.RUnlock()
			if exists {
				creator = elem.Creator
			}
			if creator == user {
				outs[count] = val
				count++
			}
		}
		outs = outs[:count]
		sort.Strings(outs)
		items = outs

	default:
		sort.Strings(items)
	}

	name := m.Author.Username
	if m.Author.ID != user {
		u, err := b.dg.User(user)
		if rsp.Error(err) {
			return
		}
		name = u.Username
	}
	b.newPageSwitcher(pageSwitcher{
		Kind:       pageSwitchInv,
		Title:      fmt.Sprintf("%s's Inventory (%d, %s%%)", name, len(items), formatFloat(float32(len(items))/float32(len(dat.elemCache))*100, 2)),
		PageGetter: b.invPageGetter,
		Items:      items,
	}, m, rsp)
}

func (b *EoD) lbCmd(m msg, rsp rsp, sort string) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	dat.lock.RLock()
	_, exists = dat.invCache[m.Author.ID]
	dat.lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}

	b.newPageSwitcher(pageSwitcher{
		Kind:       pageSwitchLdb,
		Title:      "Top Most Elements",
		PageGetter: b.lbPageGetter,
		Sort:       sort,
		User:       m.Author.ID,
	}, m, rsp)
}