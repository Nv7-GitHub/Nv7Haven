package types

import (
	"strings"
)

func (dat *ServerDat) SetCategory(cat OldCategory) {
	dat.Lock.Lock()
	dat.Categories[strings.ToLower(cat.Name)] = cat
	dat.Lock.Unlock()
}

func (dat *ServerDat) SetElement(elem OldElement) {
	dat.Lock.Lock()
	dat.Elements[strings.ToLower(elem.Name)] = elem
	dat.Lock.Unlock()
}

func (dat *ServerDat) SetComb(id string, comb Comb) {
	dat.Lock.Lock()
	dat.LastCombs[id] = comb
	dat.Lock.Unlock()
}

func (dat *ServerDat) SetInv(id string, inv Inventory) {
	dat.Lock.Lock()
	dat.Inventories[id] = inv
	dat.Lock.Unlock()
}

func (dat *ServerDat) DeleteCategory(name string) {
	dat.Lock.Lock()
	delete(dat.Categories, strings.ToLower(name))
	dat.Lock.Unlock()
}

func (dat *ServerDat) DeleteElement(name string) {
	dat.Lock.Lock()
	delete(dat.Elements, strings.ToLower(name))
	dat.Lock.Unlock()
}

func (dat *ServerDat) DeleteComb(id string) {
	dat.Lock.Lock()
	delete(dat.LastCombs, id)
	dat.Lock.Unlock()
}

func (dat *ServerDat) AddComponentMsg(id string, msg ComponentMsg) {
	dat.Lock.Lock()
	dat.ComponentMsgs[id] = msg
	dat.Lock.Unlock()
}

func (dat *ServerDat) SavePageSwitcher(id string, ps PageSwitcher) {
	dat.Lock.Lock()
	dat.PageSwitchers[id] = ps
	dat.Lock.Unlock()
}

func (dat *ServerDat) SavePoll(id string, poll Poll) {
	dat.Lock.Lock()
	dat.Polls[id] = poll
	dat.Lock.Unlock()
}

func (dat *ServerDat) AddComb(elems string, elem3 string) {
	dat.Lock.Lock()
	dat.Combos[elems] = elem3
	dat.Lock.Unlock()
}

func (dat *ServerDat) SetMsgElem(id string, elem string) {
	dat.Lock.Lock()
	if dat.ElementMsgs == nil {
		dat.ElementMsgs = make(map[string]string)
	}
	dat.ElementMsgs[id] = strings.ToLower(elem)
	dat.Lock.Unlock()
}
