package treecmds

import (
	"fmt"

	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *TreeCmds) ElemBreakdownCmd(elem string, calcTree bool, m types.Msg, rsp types.Rsp) {
	rsp.Acknowledge()

	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	el, res := db.GetElementByName(elem)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	tree := &trees.BreakDownTree{
		DB:        db,
		Breakdown: make(map[string]int),
		Added:     make(map[int]types.Empty),
		Tree:      calcTree,
		Total:     0,
	}
	suc, err := tree.AddElem(el.ID)
	if !suc {
		rsp.ErrorMessage(err)
		return
	}

	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      fmt.Sprintf("%s Breakdown (%d)", el.Name, tree.Total),
		PageGetter: b.base.InvPageGetter,
		Items:      tree.GetStringArr(),
	}, m, rsp)
}

func (b *TreeCmds) CatBreakdownCmd(catName string, calcTree bool, m types.Msg, rsp types.Rsp) {
	rsp.Acknowledge()

	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	cat, res := db.GetCat(catName)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	tree := &trees.BreakDownTree{
		DB:        db,
		Breakdown: make(map[string]int),
		Added:     make(map[int]types.Empty),
		Tree:      calcTree,
		Total:     0,
	}

	for elem := range cat.Elements {
		suc, err := tree.AddElem(elem)
		if !suc {
			rsp.ErrorMessage(err)
			return
		}
	}

	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      fmt.Sprintf("%s Breakdown (%d)", cat.Name, tree.Total),
		PageGetter: b.base.InvPageGetter,
		Items:      tree.GetStringArr(),
	}, m, rsp)
}

func (b *TreeCmds) InvBreakdownCmd(user string, calcTree bool, m types.Msg, rsp types.Rsp) {
	rsp.Acknowledge()

	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	inv := db.GetInv(user)

	tree := &trees.BreakDownTree{
		DB:        db,
		Breakdown: make(map[string]int),
		Added:     make(map[int]types.Empty),
		Tree:      calcTree,
		Total:     0,
	}

	for elem := range inv.Elements {
		/*suc, err :=*/ tree.AddElem(elem, true)
		/*	if !suc {
			rsp.ErrorMessage(err)
			return
		}*/
	}

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
		Title:      fmt.Sprintf("%s's Inventory Breakdown (%d)", name, tree.Total),
		PageGetter: b.base.InvPageGetter,
		Items:      tree.GetStringArr(),
	}, m, rsp)
}
