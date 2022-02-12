package elements

import (
	"sort"

	"github.com/Nv7-Github/Nv7Haven/eod/eodsort"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

func (b *Elements) InvCmd(user string, m types.Msg, rsp types.Rsp, sorter string, filter string, postfix bool, defaul bool) {
	rsp.Acknowledge()

	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	inv := db.GetInv(user)
	type invItem struct {
		name string
		id   int
	}
	items := make([]invItem, len(inv.Elements))
	i := 0
	db.RLock()
	for k := range inv.Elements {
		el, _ := db.GetElement(k, true)
		items[i] = invItem{el.Name, el.ID}
		i++
	}

	switch filter {
	case "madeby":
		count := 0
		outs := make([]invItem, len(items))
		for _, val := range items {
			creator := ""
			elem, res := db.GetElement(val.id, true)
			if res.Exists {
				creator = elem.Creator
			}
			if creator == user {
				outs[count] = invItem{elem.Name, elem.ID}
				count++
			}
		}
		outs = outs[:count]
		items = outs
	}
	if defaul && m.Author.ID != user {
		postfix = true
	}
	eodsort.Sort(items, len(items), func(index int) int {
		return items[index].id
	}, func(index int) string {
		return items[index].name
	}, func(index int, val string) {
		items[index].name = val
	}, sorter, m.Author.ID, db, postfix)

	text := make([]string, len(items))
	for i, v := range items {
		text[i] = v.name
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
		Kind: types.PageSwitchInv,
		Title: db.Config.LangProperty("UserInventory", map[string]interface{}{
			"Username": name,
			"Count":    len(items),
			"Percent":  util.FormatFloat(float32(len(items))/float32(len(db.Elements))*100, 2),
		}),
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
	var sortFn func(a, b int) bool
	switch sorter {
	case "made":
		sortFn = func(a, b int) bool {
			return invs[a].MadeCnt > invs[b].MadeCnt
		}

	case "signed":
		sortFn = func(a, b int) bool {
			return invs[a].SignedCnt > invs[b].SignedCnt
		}

	case "imaged":
		sortFn = func(a, b int) bool {
			return invs[a].ImagedCnt > invs[b].ImagedCnt
		}

	case "colored":
		sortFn = func(a, b int) bool {
			return invs[a].ColoredCnt > invs[b].ColoredCnt
		}

	case "catimaged":
		sortFn = func(a, b int) bool {
			return invs[a].CatImagedCnt > invs[b].CatImagedCnt
		}

	case "catcolored":
		sortFn = func(a, b int) bool {
			return invs[a].CatColoredCnt > invs[b].CatColoredCnt
		}

	default:
		sortFn = func(a, b int) bool {
			return len(invs[a].Elements) > len(invs[b].Elements)
		}
	}
	sort.Slice(invs, sortFn)

	// Convert to right format
	users := make([]string, len(invs))
	cnts := make([]int, len(invs))
	userpos := 0
	for i, v := range invs {
		users[i] = v.User
		switch sorter {
		case "made":
			cnts[i] = v.MadeCnt

		case "signed":
			cnts[i] = v.SignedCnt

		case "imaged":
			cnts[i] = v.ImagedCnt

		case "colored":
			cnts[i] = v.ColoredCnt

		case "catimaged":
			cnts[i] = v.CatImagedCnt

		case "catcolored":
			cnts[i] = v.CatColoredCnt

		default:
			cnts[i] = len(v.Elements)
		}
		if v.User == user {
			userpos = i
		}
	}

	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchLdb,
		Title:      db.Config.LangProperty("LbTitleElem", nil),
		PageGetter: b.base.LbPageGetter,

		User:    user,
		Users:   users,
		UserPos: userpos,
		Cnts:    cnts,
	}, m, rsp)
}
