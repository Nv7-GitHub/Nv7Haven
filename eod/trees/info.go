package trees

import (
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

type InfoTree struct {
	Total int
	Found int

	dat   types.ServerDat
	added map[string]types.Empty
	inv   types.Container
}

func (i *InfoTree) AddElem(name string, unlock ...bool) (bool, string) {
	if len(unlock) == 0 {
		i.dat.Lock.RLock()
		defer i.dat.Lock.RUnlock()
	}

	elem := strings.ToLower(name)
	_, exists := i.added[elem]
	if exists {
		return true, ""
	}

	el, res := i.dat.GetElement(elem, true)
	if !res.Exists {
		return false, res.Message
	}
	for _, par := range el.Parents {
		suc, err := i.AddElem(par, true)
		if !suc {
			return suc, err
		}
	}
	i.Total++
	if i.inv.Contains(elem) {
		i.Found++
	}
	i.added[elem] = types.Empty{}
	return true, ""
}

func CalcElemInfo(elem string, user string, dat types.ServerDat) (bool, string, InfoTree) {
	inv, res := dat.GetInv(user, true)
	if !res.Exists {
		return false, res.Message, InfoTree{}
	}
	i := InfoTree{
		Total: 0,
		Found: 0,
		dat:   dat,
		added: make(map[string]types.Empty),
		inv:   inv.Elements,
	}
	suc, msg := i.AddElem(elem)
	if !suc {
		return false, msg, InfoTree{}
	}
	return true, "", i
}
