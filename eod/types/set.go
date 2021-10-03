package types

import (
	"log"
	"os"
	"strings"
)

func (dat *ServerData) SetCategory(cat Category) {
	dat.Lock.Lock()
	dat.Categories[strings.ToLower(cat.Name)] = cat
	dat.Lock.Unlock()
}

func (dat *ServerData) SetElement(elem Element) {
	dat.Lock.Lock()
	dat.Elements[strings.ToLower(elem.Name)] = elem
	dat.Lock.Unlock()
}

func (dat *ServerData) SetComb(id string, comb Comb) {
	dat.Lock.Lock()
	dat.LastCombs[id] = comb
	dat.Lock.Unlock()
}

func (dat *ServerData) SetInv(id string, inv Container) {
	dat.Lock.Lock()
	dat.Inventories[id] = inv
	dat.Lock.Unlock()
}

func (dat *ServerData) DeleteCategory(name string) {
	dat.Lock.Lock()
	delete(dat.Categories, strings.ToLower(name))
	dat.Lock.Unlock()
}

func (dat *ServerData) DeleteElement(name string) {
	dat.Lock.Lock()
	delete(dat.Elements, strings.ToLower(name))
	dat.Lock.Unlock()
}

func (dat *ServerData) DeleteComb(id string) {
	dat.Lock.Lock()
	delete(dat.LastCombs, id)
	dat.Lock.Unlock()
}

func (dat *ServerData) AddComponentMsg(id string, msg ComponentMsg) {
	dat.Lock.Lock()
	dat.ComponentMsgs[id] = msg
	dat.Lock.Unlock()
}

func (dat *ServerData) SavePageSwitcher(id string, ps PageSwitcher) {
	dat.Lock.Lock()
	dat.PageSwitchers[id] = ps
	dat.Lock.Unlock()
}

func (dat *ServerData) SavePoll(id string, poll Poll) {
	log.SetOutput(os.Stdout)
	log.Println("save poll", id)
	dat.Lock.Lock()
	dat.Polls[id] = poll
	dat.Lock.Unlock()
}

func (dat *ServerData) AddComb(elems string, elem3 string) {
	dat.Lock.Lock()
	dat.Combos[elems] = elem3
	dat.Lock.Unlock()
}

func (dat *ServerData) SetMsgElem(id string, elem string) {
	dat.Lock.Lock()
	if dat.ElementMsgs == nil {
		dat.ElementMsgs = make(map[string]string)
	}
	dat.ElementMsgs[id] = strings.ToLower(elem)
	dat.Lock.Unlock()
}
