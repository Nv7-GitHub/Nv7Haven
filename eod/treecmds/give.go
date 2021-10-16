package treecmds

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *TreeCmds) GiveCmd(elem string, giveTree bool, user string, m types.Msg, rsp types.Rsp) {
	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		return
	}

	inv, res := dat.GetInv(user, true)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	el, res := dat.GetElement(elem)
	if !res.Exists {
		rsp.Resp(res.Message)
		return
	}

	msg, suc := giveElem(dat, giveTree, elem, &inv)
	if !suc {
		rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", msg))
		return
	}

	dat.SetInv(user, inv)

	b.lock.Lock()
	b.dat[m.GuildID] = dat
	b.lock.Unlock()
	b.base.SaveInv(m.GuildID, user, true, true)

	rsp.Resp("Successfully gave element **" + el.Name + "**!")
}

func (b *TreeCmds) GiveCatCmd(catName string, giveTree bool, user string, m types.Msg, rsp types.Rsp) {
	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		return
	}

	inv, res := dat.GetInv(user, true)
	if !exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	cat, res := dat.GetCategory(catName)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	for elem := range cat.Elements {
		_, res := dat.GetElement(elem)
		if !res.Exists {
			rsp.Resp(fmt.Sprintf("Element **%s** doesn't exist!", elem))
			return
		}

		msg, suc := giveElem(dat, giveTree, elem, &inv)
		if !suc {
			rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", msg))
			return
		}
	}

	dat.SetInv(user, inv)

	b.lock.Lock()
	b.dat[m.GuildID] = dat
	b.lock.Unlock()
	b.base.SaveInv(m.GuildID, user, true, true)

	rsp.Resp("Successfully gave all elements in category **" + cat.Name + "**!")
}

func giveElem(dat types.ServerData, giveTree bool, elem string, out *types.Inventory) (string, bool) {
	el, res := dat.GetElement(elem)
	if !res.Exists {
		return elem, false
	}
	if giveTree {
		for _, parent := range el.Parents {
			if len(strings.TrimSpace(parent)) == 0 {
				continue
			}
			exists := out.Elements.Contains(parent)
			if !exists {
				msg, suc := giveElem(dat, giveTree, parent, out)
				if !suc {
					return msg, false
				}
			}
		}
	}
	(*out).Elements.Add(el.Name)
	return "", true
}

func (b *TreeCmds) GiveAllCmd(user string, m types.Msg, rsp types.Rsp) {
	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		return
	}
	inv, res := dat.GetInv(user, user == m.Author.ID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	for k := range dat.Elements {
		inv.Elements.Add(k)
	}

	dat.SetInv(user, inv)

	b.lock.Lock()
	b.dat[m.GuildID] = dat
	b.lock.Unlock()
	b.base.SaveInv(m.GuildID, user, true, true)
	rsp.Resp("Successfully gave every element to <@" + user + ">!")
}
