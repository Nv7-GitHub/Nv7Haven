package types

import "strings"

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
