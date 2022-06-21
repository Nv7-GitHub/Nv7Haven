package treecmds

import (
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *TreeCmds) GiveCmd(elem string, giveTree bool, user string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}

	inv := db.GetInv(user)

	el, res := db.GetElementByName(elem)
	if !res.Exists {
		rsp.Resp(res.Message)
		return
	}

	msg, suc := giveElem(db, giveTree, el.ID, inv)
	if !suc {
		rsp.ErrorMessage(msg)
		return
	}

	opts := []bool{true, true}
	if giveTree {
		opts = []bool{true}
	}
	err := db.SaveInv(inv, opts...)
	if rsp.Error(err) {
		return
	}

	rsp.Resp(db.Config.LangProperty("GiveElem", el.Name))
}

func (b *TreeCmds) GiveCatCmd(catName string, giveTree bool, user string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}

	inv := db.GetInv(user)

	var els map[int]types.Empty
	var lock *sync.RWMutex
	catv, res := db.GetCat(catName)
	if !res.Exists {
		vcat, res := db.GetVCat(catName)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
		catName = vcat.Name
		els, res = b.base.CalcVCat(vcat, db, true)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
	} else {
		els = catv.Elements
		lock = catv.Lock
		catName = catv.Name
	}

	if lock != nil {
		lock.RLock()
	}
	for elem := range els {
		_, res := db.GetElement(elem)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}

		msg, suc := giveElem(db, giveTree, elem, inv)
		if !suc {
			rsp.ErrorMessage(db.Config.LangProperty("DoesntExist", msg))
			return
		}
	}
	if lock != nil {
		lock.RUnlock()
	}

	err := db.SaveInv(inv)
	if rsp.Error(err) {
		return
	}

	rsp.Resp(db.Config.LangProperty("GiveCat", catName))
}

func giveElem(db *eodb.DB, giveTree bool, elem int, out *types.Inventory) (string, bool) {
	el, res := db.GetElement(elem)
	if !res.Exists {
		return res.Message, false
	}
	if giveTree {
		for _, parent := range el.Parents {
			exists := out.Contains(parent)
			if !exists {
				msg, suc := giveElem(db, giveTree, parent, out)
				if !suc {
					return msg, false
				}
			}
		}
	}
	out.Add(elem)
	return "", true
}

func (b *TreeCmds) GiveAllCmd(user string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}

	inv := db.GetInv(user)

	db.RLock()
	for id := range db.Elements {
		inv.Add(id)
	}
	db.RUnlock()

	err := db.SaveInv(inv)
	if rsp.Error(err) {
		return
	}

	rsp.Resp(db.Config.LangProperty("GiveAll", user))
}
