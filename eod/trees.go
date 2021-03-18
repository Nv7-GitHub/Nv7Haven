package eod

import (
	"fmt"
	"strings"
)

func (b *EoD) giveCmd(elem string, giveTree bool, user string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	inv, exists := dat.invCache[user]
	if !exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}

	el, exists := dat.elemCache[strings.ToLower(elem)]
	if !exists {
		rsp.Resp(fmt.Sprintf("Element %s doesn't exist!", elem))
		return
	}

	msg, suc := giveElem(dat.elemCache, giveTree, elem, &inv)
	if !suc {
		rsp.Resp(fmt.Sprintf("Element %s doesn't exist!", msg))
		return
	}

	dat.invCache[user] = inv
	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()
	rsp.Resp("Successfully gave element " + el.Name + "!")
}

func giveElem(elemCache map[string]element, giveTree bool, elem string, out *map[string]empty) (string, bool) {
	el, exists := elemCache[strings.ToLower(elem)]
	if !exists {
		return elem, false
	}
	if giveTree {
		for _, parent := range el.Parents {
			msg, suc := giveElem(elemCache, giveTree, parent, out)
			if !suc {
				return msg, false
			}
		}
	}
	(*out)[strings.ToLower(el.Name)] = empty{}
	return "", true
}
