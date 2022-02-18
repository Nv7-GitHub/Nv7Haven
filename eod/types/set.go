package types

func (dat *ServerData) SetComb(id string, comb Comb) {
	dat.Lock()
	dat.LastCombs[id] = comb
	dat.Unlock()
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

func (dat *ServerData) AddModal(id string, handler ModalHandler) {
	dat.Lock()
	dat.Modals[id] = handler
	dat.Unlock()
}

func (dat *ServerData) SavePageSwitcher(id string, ps PageSwitcher) {
	dat.Lock()
	dat.PageSwitchers[id] = ps
	dat.Unlock()
}

func (dat *ServerData) SetMsgElem(id string, elem int) {
	dat.Lock()
	dat.ElementMsgs[id] = elem
	dat.Unlock()
}
