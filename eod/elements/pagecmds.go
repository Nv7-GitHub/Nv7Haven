package elements

import (
	"fmt"
	"sort"

	"github.com/Nv7-Github/Nv7Haven/eod/eodsort"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

func (b *Elements) InvCmd(user string, m types.Msg, rsp types.Rsp, sorter string, filter string) {
	rsp.Acknowledge()

	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	inv := db.GetInv(user)
	items := make([]int, len(inv.Elements))
	i := 0
	db.RLock()
	for k := range inv.Elements {
		el, _ := db.GetElement(k, true)
		items[i] = el.ID
		i++
	}

	switch filter {
	case "madeby":
		count := 0
		outs := make([]int, len(items))
		for _, val := range items {
			creator := ""
			elem, res := db.GetElement(val, true)
			if res.Exists {
				creator = elem.Creator
			}
			if creator == user {
				outs[count] = elem.ID
				count++
			}
		}
		outs = outs[:count]
	}
	eodsort.SortElemList(items, sorter, db)

	text := make([]string, len(items))
	for i, v := range items {
		elem, res := db.GetElement(v, true)
		if !res.Exists {
			continue
		}
		text[i] = elem.Name
	}
	db.RUnlock()

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
		Title:      fmt.Sprintf("%s's Inventory (%d, %s%%)", name, len(items), util.FormatFloat(float32(len(items))/float32(len(db.Elements))*100, 2)),
		PageGetter: b.base.InvPageGetter,
		Items:      text,
	}, m, rsp)
}

func (b *Elements) LbCmd(m types.Msg, rsp types.Rsp, sorter string, user string) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	db.GetInv(user) // Make user exist

	// Sort invs
	invs := make([]*types.Inventory, len(db.Invs()))
	i := 0
	for _, v := range db.Invs() {
		invs[i] = v
		i++
	}
	sortFn := func(a, b int) bool {
		return len(invs[a].Elements) > len(invs[b].Elements)
	}
	if sorter == "made" {
		sortFn = func(a, b int) bool {
			return invs[a].MadeCnt > invs[b].MadeCnt
		}
	}
	sort.Slice(invs, sortFn)

	// Convert to right format
	users := make([]string, len(invs))
	cnts := make([]int, len(invs))
	userpos := 0
	for i, v := range invs {
		users[i] = v.User
		if sorter == "count" {
			cnts[i] = len(v.Elements)
		} else {
			cnts[i] = v.MadeCnt
		}
		if v.User == user {
			userpos = i
		}
	}

	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchLdb,
		Title:      "Top Most Elements",
		PageGetter: b.base.LbPageGetter,

		User:    user,
		Users:   users,
		UserPos: userpos,
		Cnts:    cnts,
	}, m, rsp)
}
