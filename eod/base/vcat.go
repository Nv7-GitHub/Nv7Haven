package base

import (
	"regexp"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *Base) CalcVCat(vcat *types.VirtualCategory, db *eodb.DB) (map[int]types.Empty, types.GetResponse) {
	var out map[int]types.Empty
	switch vcat.Rule {
	case types.VirtualCategoryRuleRegex:
		reg := regexp.MustCompile(vcat.Data["regex"].(string))

		out = make(map[int]types.Empty)
		db.RLock()
		for _, elem := range db.Elements {
			if reg.MatchString(elem.Name) {
				out[elem.ID] = types.Empty{}
			}
		}
		db.RUnlock()

	case types.VirtualCategoryRuleInvFilter:
		inv := db.GetInv(vcat.Data["user"].(string))
		switch vcat.Data["filter"].(string) {
		case "madeby":
			out = make(map[int]types.Empty)
			inv.Lock.RLock()
			db.RLock()
			for k := range inv.Elements {
				el, res := db.GetElement(k, true)
				if res.Exists && el.Creator == inv.User {
					out[k] = types.Empty{}
				}
			}
			db.RUnlock()
			inv.Lock.RUnlock()

		default:
			out = make(map[int]types.Empty, len(inv.Elements))
			inv.Lock.RLock()
			for k := range inv.Elements {
				out[k] = types.Empty{}
			}
			inv.Lock.RUnlock()
		}

	case types.VirtualCategoryRuleSetOperation:
		// TODO: implement

	case types.VirtualCategoryRuleAllElements:
		out = make(map[int]types.Empty, len(db.Elements))
		db.RLock()
		for k := range db.Elements {
			out[k] = types.Empty{}
		}
		db.RUnlock()
	}

	return out, types.GetResponse{Exists: true}
}
