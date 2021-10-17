package elements

import (
	"fmt"
	"regexp"
	"sort"
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
	items := make([]string, len(inv.Elements))
	i := 0
	dat.Lock.RLock()
	for k := range inv.Elements {
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

func (b *Elements) LbCmd(m types.Msg, rsp types.Rsp, sorter string, user string) {
	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		return
	}
	_, res := dat.GetInv(user, user == m.Author.ID) // Check if user exists
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	// Sort invs
	invs := make([]types.Inventory, len(dat.Inventories))
	i := 0
	for _, v := range dat.Inventories {
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

func (b *Elements) ElemSearchCmd(search string, sort string, regex bool, m types.Msg, rsp types.Rsp) {
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

	items := make(map[string]types.Empty)
	if regex {
		reg, err := regexp.Compile(search)
		if rsp.Error(err) {
			return
		}
		for _, el := range dat.Elements {
			m := reg.Find([]byte(el.Name))
			if m != nil {
				items[el.Name] = types.Empty{}
			}
		}
	} else {
		s := strings.ToLower(search)
		for e, el := range dat.Elements {
			if strings.Contains(e, s) {
				items[el.Name] = types.Empty{}
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

	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      "Element Search",
		PageGetter: b.base.InvPageGetter,
		Items:      txt,
		User:       m.Author.ID,
	}, m, rsp)
}
