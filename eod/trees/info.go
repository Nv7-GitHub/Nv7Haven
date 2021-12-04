package trees

import (
	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

type InfoTree struct {
	Total int
	Found int

	DB    *eodb.DB
	added map[int]types.Empty
	inv   *types.Inventory
}

func (i *InfoTree) AddElem(elem int, unlock ...bool) (bool, string) {
	if len(unlock) == 0 {
		i.DB.RLock()
		i.inv.Lock.RLock()
		defer i.DB.RUnlock()
		defer i.inv.Lock.RUnlock()
	}

	_, exists := i.added[elem]
	if exists {
		return true, ""
	}

	el, res := i.DB.GetElement(elem, true)
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
	if i.inv.Contains(elem, true) {
		i.Found++
	}
	i.added[elem] = types.Empty{}
	return true, ""
}

func CalcElemInfo(elem int, user string, db *eodb.DB) (bool, string, InfoTree) {
	inv := db.GetInv(user)
	i := InfoTree{
		Total: 0,
		Found: 0,
		DB:    db,
		added: make(map[int]types.Empty),
		inv:   inv,
	}
	suc, msg := i.AddElem(elem)
	if !suc {
		return false, msg, InfoTree{}
	}
	return true, "", i
}
