package api

import "github.com/Nv7-Github/Nv7Haven/eod/api/data"

func (a *API) MethodCategory(params map[string]any, id, gld string) data.Response {
	nm, exists := params["name"]
	if !exists {
		return data.RSPBadRequest
	}
	name, ok := nm.(string)
	if !ok {
		return data.RSPBadRequest
	}

	// Get data
	db, res := a.GetDB(gld)
	if !res.Exists {
		return data.RSPError(res.Message)
	}
	cat, res := db.GetCat(name)
	if !res.Exists {
		vcat, res := db.GetVCat(name)
		if !res.Exists {
			return data.RSPError(res.Message)
		}
		els, res := a.base.CalcVCat(vcat, db, true)
		if !res.Exists {
			return data.RSPError(res.Message)
		}
		elems := make([]int, len(els))
		i := 0
		for k := range cat.Elements {
			elems[i] = k
			i++
		}
		return data.RSPSuccess(map[string]any{"elems": elems})
	}

	cat.Lock.RLock()
	elems := make([]int, len(cat.Elements))
	i := 0
	for k := range cat.Elements {
		elems[i] = k
		i++
	}
	cat.Lock.RUnlock()

	return data.RSPSuccess(map[string]any{"elems": elems})
}
