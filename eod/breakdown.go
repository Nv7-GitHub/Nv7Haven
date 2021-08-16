package eod

import (
	"fmt"

	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *EoD) elemBreakdownCmd(elem string, calcTree bool, m types.Msg, rsp types.Rsp) {
	rsp.Acknowledge()

	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("Guild isn't setup yet!")
		return
	}

	el, res := dat.GetElement(elem)
	if !exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	tree := &trees.BreakDownTree{
		Dat:       dat,
		Breakdown: make(map[string]int),
		Added:     make(map[string]types.Empty),
		Tree:      calcTree,
		Total:     0,
	}
	suc, err := tree.AddElem(el.Name)
	if !suc {
		rsp.ErrorMessage(err)
		return
	}

	b.newPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      fmt.Sprintf("%s Breakdown (%d)", el.Name, tree.Total),
		PageGetter: b.invPageGetter,
		Items:      tree.GetStringArr(),
	}, m, rsp)
}

func (b *EoD) catBreakdownCmd(catName string, calcTree bool, m types.Msg, rsp types.Rsp) {
	rsp.Acknowledge()

	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("Guild isn't setup yet!")
		return
	}

	cat, res := dat.GetCategory(catName)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	tree := &trees.BreakDownTree{
		Dat:       dat,
		Breakdown: make(map[string]int),
		Added:     make(map[string]types.Empty),
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

	b.newPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      fmt.Sprintf("%s Breakdown (%d)", cat.Name, tree.Total),
		PageGetter: b.invPageGetter,
		Items:      tree.GetStringArr(),
	}, m, rsp)
}
