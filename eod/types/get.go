package types

import "github.com/Nv7-Github/Nv7Haven/eod/translation"

type GetResponse struct {
	Exists  bool
	Message string
}

func (dat *ServerData) GetComb(id string) (Comb, GetResponse) {
	dat.RLock()
	comb, exists := dat.LastCombs[id]
	dat.RUnlock()
	if !exists {
		return Comb{}, GetResponse{
			Exists:  false,
			Message: "You haven't combined anything!",
		}
	}
	return comb, GetResponse{Exists: true}
}

func (dat *ServerData) GetPageSwitcher(id string) (PageSwitcher, GetResponse) {
	dat.RLock()
	ps, exists := dat.PageSwitchers[id]
	dat.RUnlock()
	if !exists {
		return PageSwitcher{}, GetResponse{
			Exists:  false,
			Message: "Page switcher doesn't exist!",
		}
	}
	return ps, GetResponse{Exists: true}
}

func (dat *ServerData) GetMsgElem(id string) (int, GetResponse) {
	dat.RLock()
	elem, exists := dat.ElementMsgs[id]
	dat.RUnlock()
	if !exists {
		return 0, GetResponse{
			Exists:  false,
			Message: "Message doesn't have an element!",
		}
	}
	return elem, GetResponse{Exists: true}
}

func (l *ServerConfig) LangProperty(key string, params interface{}) string {
	return translation.LangProperty(l.LanguageFile, key, params)
}
