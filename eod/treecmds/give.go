package treecmds

import (
	"fmt"

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

	rsp.Resp(fmt.Sprintf(db.Config.LangProperty("GiveElem"), el.Name)
}

func (b *TreeCmds) GiveCatCmd(catName string, giveTree bool, user string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}

	inv := db.GetInv(user)

	cat, res := db.GetCat(catName)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	for elem := range cat.Elements {
		_, res := db.GetElement(elem)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}

		msg, suc := giveElem(db, giveTree, elem, inv)
		if !suc {
			rsp.ErrorMessage(fmt.Sprintf(db.Config.LangProperty("DoesntExist"), msg))
			return
		}
	}

	err := db.SaveInv(inv)
	if rsp.Error(err) {
		return
	}

	rsp.Resp(fmt.Sprintf(db.Config.LangProperty("GiveCat"), cat.Name)
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

	rsp.Resp(fmt.Sprintf(db.Config.LangProperty("GiveAll"), user)
}
