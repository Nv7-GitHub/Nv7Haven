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

func (dat *ServerData) SetComb(id string, comb Comb) {
	dat.Lock()
	dat.LastCombs[id] = comb
	dat.Unlock()
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

func (dat *ServerData) DeleteComb(id string) {
	dat.Lock()
	delete(dat.LastCombs, id)
	dat.Unlock()
}

func (dat *ServerData) AddComponentMsg(id string, msg ComponentMsg) {
	dat.Lock()
	dat.ComponentMsgs[id] = msg
	dat.Unlock()
}

func (dat *ServerData) SavePageSwitcher(id string, ps PageSwitcher) {
	dat.Lock()
	dat.PageSwitchers[id] = ps
	dat.Unlock()
}

func (dat *ServerDat) SavePoll(id string, poll OldPoll) {
	dat.Lock.Lock()
	dat.Polls[id] = poll
	dat.Lock.Unlock()
}

func (dat *ServerDat) AddComb(elems string, elem3 string) {
	dat.Lock.Lock()
	dat.Combos[elems] = elem3
	dat.Lock.Unlock()
}

func (dat *ServerData) SetMsgElem(id string, elem int) {
	dat.Lock()
	dat.ElementMsgs[id] = elem
	dat.Unlock()
}
