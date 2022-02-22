package treecmds

import (
	"sync"

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
		Kind: types.PageSwitchInv,
		Title: db.Config.LangProperty("BreakdownTitle", map[string]interface{}{
			"Title": el.Name,
			"Count": tree.Total,
		}),
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

	var lock *sync.RWMutex
	var els map[int]types.Empty
	catv, res := db.GetCat(catName)
	if !res.Exists {
		vcat, res := db.GetVCat(catName)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
		catName = vcat.Name
		els, res = b.base.CalcVCat(vcat, db)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
	} else {
		els = catv.Elements
		catName = catv.Name
	}

	tree := &trees.BreakDownTree{
		DB:        db,
		Breakdown: make(map[string]int),
		Added:     make(map[int]types.Empty),
		Tree:      calcTree,
		Total:     0,
	}

	if lock != nil {
		lock.RLock()
	}
	for elem := range els {
		suc, err := tree.AddElem(elem)
		if !suc {
			rsp.ErrorMessage(err)
			return
		}
	}
	if lock != nil {
		lock.RUnlock()
	}

	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind: types.PageSwitchInv,
		Title: db.Config.LangProperty("BreakdownTitle", map[string]interface{}{
			"Title": catName,
			"Count": tree.Total,
		}),
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
		Kind: types.PageSwitchInv,
		Title: db.Config.LangProperty("InvBreakdownTitle", map[string]interface{}{
			"User":  name,
			"Count": tree.Total,
		}),
		PageGetter: b.base.InvPageGetter,
		Items:      tree.GetStringArr(),
	}, m, rsp)
}
