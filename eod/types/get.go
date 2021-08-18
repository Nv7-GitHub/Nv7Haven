package types

import (
	"fmt"
	"strings"
)

type GetResponse struct {
	Exists  bool
	Message string
}

func (dat *ServerData) GetElement(name string, noLock ...bool) (Element, GetResponse) {
	if len(noLock) == 0 {
		dat.Lock.RLock()
	}
	el, exists := dat.Elements[strings.ToLower(name)]
	if len(noLock) == 0 {
		dat.Lock.RUnlock()
	}
	if !exists {
		return Element{}, GetResponse{
			Exists:  false,
			Message: fmt.Sprintf("Element **%s** doesn't exist!", name),
		}
	}
	return el, GetResponse{Exists: true}
}

func (dat *ServerData) GetInv(id string, you bool) (Container, GetResponse) {
	dat.Lock.RLock()
	inv, exists := dat.Inventories[id]
	dat.Lock.RUnlock()
	if !exists {
		var response GetResponse
		if you {
			response = GetResponse{
				Exists:  false,
				Message: "You don't have an inventory!",
			}
		} else {
			response = GetResponse{
				Exists:  false,
				Message: fmt.Sprintf("User <@%s> doesn't have an inventory!", id),
			}
		}
		return nil, response
	}
	return inv, GetResponse{Exists: true}
}

func (dat *ServerData) GetCategory(name string, noLock ...bool) (Category, GetResponse) {
	if len(noLock) == 0 {
		dat.Lock.RLock()
	}
	cat, exists := dat.Categories[strings.ToLower(name)]
	if len(noLock) == 0 {
		dat.Lock.RUnlock()
	}
	if !exists {
		return Category{}, GetResponse{
			Exists:  false,
			Message: fmt.Sprintf("Category **%s** doesn't exist!", name),
		}
	}
	return cat, GetResponse{Exists: true}
}

func (dat *ServerData) GetComb(id string) (Comb, GetResponse) {
	dat.Lock.RLock()
	comb, exists := dat.LastCombs[id]
	dat.Lock.RUnlock()
	if !exists {
		return Comb{}, GetResponse{
			Exists:  false,
			Message: "You haven't combined anything!",
		}
	}
	return comb, GetResponse{Exists: true}
}

func (dat *ServerData) GetPageSwitcher(id string) (PageSwitcher, GetResponse) {
	dat.Lock.RLock()
	ps, exists := dat.PageSwitchers[id]
	dat.Lock.RUnlock()
	if !exists {
		return PageSwitcher{}, GetResponse{
			Exists:  false,
			Message: "Page switcher doesn't exist!",
		}
	}
	return ps, GetResponse{Exists: true}
}

func (dat *ServerData) GetPoll(id string) (Poll, GetResponse) {
	dat.Lock.RLock()
	poll, exists := dat.Polls[id]
	dat.Lock.RUnlock()
	if !exists {
		return Poll{}, GetResponse{
			Exists:  false,
			Message: "Poll doesn't exist!",
		}
	}
	return poll, GetResponse{Exists: true}
}

func (dat *ServerData) GetCombo(elems string) (string, GetResponse) {
	dat.Lock.RLock()
	elem3, exists := dat.Combos[elems]
	dat.Lock.RUnlock()
	if !exists {
		return "", GetResponse{
			Exists: false,
		}
	}
	return elem3, GetResponse{Exists: true}
}
