package eod

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

func (b *EoD) invCmd(user string, m types.Msg, rsp types.Rsp, sorter string, filter string) {
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
	sortElemList(items, sorter, dat)
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
		Title:      fmt.Sprintf("%s's Inventory (%d, %s%%)", name, len(items), util.FormatFloat(float32(len(items))/float32(len(dat.Elements))*100, 2)),
		PageGetter: b.invPageGetter,
		Items:      items,
	}, m, rsp)
}

func (b *EoD) lbCmd(m types.Msg, rsp types.Rsp, sort string, user string) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	_, res := dat.GetInv(user, user == m.Author.ID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	b.newPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchLdb,
		Title:      "Top Most Elements",
		PageGetter: b.lbPageGetter,
		Sort:       sort,
		User:       user,
	}, m, rsp)
}

func (b *EoD) foundCmd(elem string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	rsp.Acknowledge()

	el, res := dat.GetElement(elem)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	var foundCnt int
	err := b.db.QueryRow(`SELECT COUNT(1) as cnt FROM eod_inv WHERE guild=? AND (JSON_EXTRACT(inv, CONCAT('$."', LOWER(?), '"')) IS NOT NULL)`, m.GuildID, el.Name).Scan(&foundCnt)
	if rsp.Error(err) {
		return
	}

	found, err := b.db.Query(`SELECT user as cnt FROM eod_inv WHERE guild=? AND (JSON_EXTRACT(inv, CONCAT('$."', LOWER(?), '"')) IS NOT NULL)`, m.GuildID, el.Name)
	if rsp.Error(err) {
		return
	}
	defer found.Close()

	out := make([]string, foundCnt)
	i := 0

	var user string
	for found.Next() {
		err = found.Scan(&user)
		if rsp.Error(err) {
			return
		}

		out[i] = fmt.Sprintf("<@%s>", user)
		i++
	}

	b.newPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      fmt.Sprintf("%s Found (%d)", el.Name, len(out)),
		PageGetter: b.invPageGetter,
		Items:      out,
		User:       m.Author.ID,
	}, m, rsp)
}

func (b *EoD) elemSearchCmd(search string, m types.Msg, rsp types.Rsp) {
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
	if util.IsWildcard(search) {
		for val := range util.Wildcards {
			search = strings.ReplaceAll(search, string([]rune{val}), string([]rune{'\\', val}))
		}
	}

	b.newPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchSearch,
		Title:      "Element Search",
		PageGetter: b.searchPageGetter,
		Search:     search,
		User:       m.Author.ID,
	}, m, rsp)
}
