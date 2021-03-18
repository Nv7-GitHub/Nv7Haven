package eod

import (
	"bytes"
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

func (b *EoD) calcTreeCmd(elem string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RLock()
	if !exists {
		return
	}
	txt, suc, msg := calcTree(dat.elemCache, elem)
	if !suc {
		rsp.ErrorMessage(fmt.Sprintf("Element %s doesn't exist!", msg))
	}
	if len(txt) <= 2000 {
		rsp.Resp(txt)
		return
	}
	rsp.Resp(txt) // Get <user used command> thing to show up

	buf := bytes.NewBufferString(txt)
	b.dg.ChannelFileSend(m.ChannelID, "path.txt", buf)
}

// Treecalc
func calcTree(elemCache map[string]element, elem string) (string, bool, string) {
	t := tree{
		text:      "",
		elemCache: elemCache,
		calced:    make(map[string]empty),
		num:       1,
	}
	suc, msg := t.addElem(elem)
	if len(t.text) > 2000 {
		return t.rawTxt, suc, msg
	}
	return t.text, suc, msg
}

type tree struct {
	text      string
	rawTxt    string
	elemCache map[string]element
	calced    map[string]empty
	num       int
}

func (t *tree) addElem(elem string) (bool, string) {
	_, exists := t.calced[strings.ToLower(elem)]
	if !exists {
		el, exists := t.elemCache[strings.ToLower(elem)]
		if !exists {
			return false, elem
		}
		for _, parent := range el.Parents {
			suc, msg := t.addElem(parent)
			if !suc {
				return false, msg
			}
		}
		if len(el.Parents) == 2 {
			t.text += fmt.Sprintf("%d. %s + %s = **%s**\n", t.num, el.Parents[0], el.Parents[1], el.Name)
			t.rawTxt += fmt.Sprintf("%d. %s + %s = %s\n", t.num, el.Parents[0], el.Parents[1], el.Name)
			t.num++
		}
		t.calced[strings.ToLower(elem)] = empty{}
	}
	return true, ""
}
