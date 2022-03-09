package base

import (
	"regexp"
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *Base) CatOpPollTitle(c types.CategoryOperation, db *eodb.DB) string {
	switch c {
	case types.CatOpUnion:
		return db.Config.LangProperty("UnionPoll", nil)

	case types.CatOpIntersect:
		return db.Config.LangProperty("IntersectPoll", nil)

	case types.CatOpDiff:
		return db.Config.LangProperty("DiffPoll", nil)

	default:
		return "unknown"
	}
}

var Elemlock = &sync.RWMutex{}

var Allelements = make(map[string]map[int]types.Empty)

var Madebylock = &sync.RWMutex{}

var Madeby = make(map[string]map[string]map[int]types.Empty)

func (b *Base) VCatDependencies(cat string, deps *map[string]types.Empty, db *eodb.DB, notfirst ...bool) {
	_, exists := (*deps)[cat]
	if exists {
		return
	}
	if len(notfirst) > 0 {
		(*deps)[cat] = types.Empty{}
	}
	vcat, res := db.GetVCat(cat)
	if !res.Exists {
		return
	}
	if vcat.Rule != types.VirtualCategoryRuleSetOperation {
		return
	}
	lhs := vcat.Data["lhs"].(string)
	rhs := vcat.Data["rhs"].(string)
	b.VCatDependencies(lhs, deps, db, true)
	b.VCatDependencies(rhs, deps, db, true)
}

func (b *Base) CalcVCat(vcat *types.VirtualCategory, db *eodb.DB) (map[int]types.Empty, types.GetResponse) {
	var out map[int]types.Empty
	switch vcat.Rule {
	case types.VirtualCategoryRuleRegex:
		if vcat.Cache != nil { // Has cache
			out = make(map[int]types.Empty, len(vcat.Cache))
			for k := range vcat.Cache {
				out[k] = types.Empty{}
			}
			break
		}

		// Populate cache
		reg := regexp.MustCompile(vcat.Data["regex"].(string))
		out = make(map[int]types.Empty)
		db.RLock()
		for _, elem := range db.Elements {
			if reg.MatchString(elem.Name) {
				out[elem.ID] = types.Empty{}
			}
		}
		db.RUnlock()

		vcat.Cache = out

		// Save
		err := db.SaveCatCache(vcat.Name, vcat.Cache)
		if err != nil {
			return nil, types.GetResponse{
				Exists:  false,
				Message: err.Error(),
			}
		}

	case types.VirtualCategoryRuleInvFilter:
		inv := db.GetInv(vcat.Data["user"].(string))
		switch vcat.Data["filter"].(string) {
		case "madeby":
			// Get cat
			Madebylock.RLock()
			gld, exists := Madeby[db.Guild]
			Madebylock.RUnlock()
			if exists {
				Madebylock.RLock()
				out, exists = gld[vcat.Data["user"].(string)]
				Madebylock.RUnlock()
				if exists {
					break
				}
			}

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

			// Save to cache
			Madebylock.Lock()
			gld, exists = Madeby[db.Guild]
			if !exists {
				gld = make(map[string]map[int]types.Empty)
			}
			Madeby[db.Guild] = gld
			gld[vcat.Data["user"].(string)] = out
			Madebylock.Unlock()

		default:
			out = make(map[int]types.Empty, len(inv.Elements))
			inv.Lock.RLock()
			for k := range inv.Elements {
				out[k] = types.Empty{}
			}
			inv.Lock.RUnlock()
		}

	case types.VirtualCategoryRuleSetOperation:
		deps := make(map[string]types.Empty)
		b.VCatDependencies(vcat.Name, &deps, db)
		_, exists := deps[vcat.Name]
		if exists {
			out = make(map[int]types.Empty)
			break
		}

		// Calc lhs
		var lhselems map[int]types.Empty
		lhs := vcat.Data["lhs"].(string)
		cat, res := db.GetCat(lhs)
		if !res.Exists {
			vcat, res := db.GetVCat(lhs)
			if !res.Exists {
				lhselems = make(map[int]types.Empty)
			} else {
				lhselems, res = b.CalcVCat(vcat, db)
				if !res.Exists {
					lhselems = make(map[int]types.Empty)
				}
			}
		} else {
			lhselems = make(map[int]types.Empty, len(cat.Elements))
			cat.Lock.RLock()
			for k := range cat.Elements {
				lhselems[k] = types.Empty{}
			}
			cat.Lock.RUnlock()
		}

		// Calc rhs
		var rhselems map[int]types.Empty
		rhs := vcat.Data["rhs"].(string)
		cat, res = db.GetCat(rhs)
		if !res.Exists {
			vcat, res := db.GetVCat(rhs)
			if !res.Exists {
				rhselems = make(map[int]types.Empty)
			} else {
				rhselems, res = b.CalcVCat(vcat, db)
				if !res.Exists {
					rhselems = make(map[int]types.Empty)
				}
			}
		} else {
			rhselems = make(map[int]types.Empty, len(cat.Elements))
			cat.Lock.RLock()
			for k := range cat.Elements {
				rhselems[k] = types.Empty{}
			}
			cat.Lock.RUnlock()
		}

		// Operations
		switch types.CategoryOperation(vcat.Data["operation"].(string)) {
		case types.CatOpUnion:
			out = make(map[int]types.Empty, len(lhselems)+len(rhselems))
			for k := range lhselems {
				out[k] = types.Empty{}
			}
			for k := range rhselems {
				out[k] = types.Empty{}
			}

		case types.CatOpIntersect:
			out = make(map[int]types.Empty)
			for k := range lhselems {
				if _, ok := rhselems[k]; ok {
					out[k] = types.Empty{}
				}
			}
			for k := range rhselems {
				if _, ok := lhselems[k]; ok {
					out[k] = types.Empty{}
				}
			}

		case types.CatOpDiff:
			out = make(map[int]types.Empty, len(lhselems))
			for k := range lhselems {
				if _, ok := rhselems[k]; !ok {
					out[k] = types.Empty{}
				}
			}
		}

	case types.VirtualCategoryRuleAllElements:
		// Check if available in cache
		var exists bool
		Elemlock.RLock()
		out, exists = Allelements[db.Guild]
		Elemlock.RUnlock()
		if !exists {
			// Calculate
			out = make(map[int]types.Empty, len(db.Elements))
			db.RLock()
			for _, el := range db.Elements {
				out[el.ID] = types.Empty{}
			}
			db.RUnlock()

			Elemlock.Lock()
			Allelements[db.Guild] = out
			Elemlock.Unlock()
		}
	}

	return out, types.GetResponse{Exists: true}
}
