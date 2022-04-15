package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Nv7-Github/Nv7Haven/eod/api/data"
)

func (a *API) MethodElem(params map[string]any, id, gld string) data.Response {
	// Process params
	name, ok := params["name"]
	if !ok {
		return data.RSPBadRequest
	}
	nm, ok := name.(string)
	if !ok {
		return data.RSPBadRequest
	}

	// Get data
	db, res := a.GetDB(gld)
	if !res.Exists {
		return data.RSPError(res.Message)
	}
	elem, res := db.GetElementByName(nm)
	if !res.Exists {
		return data.RSPError(res.Message)
	}
	return data.RSPSuccess(map[string]any{"id": elem.ID})
}

func (a *API) MethodElemInfo(params map[string]any, id, gld string) data.Response {
	// Process params
	v, ok := params["ids"]
	if !ok {
		return data.RSPBadRequest
	}
	idsV, ok := v.([]interface{})
	if !ok {
		return data.RSPBadRequest
	}
	ids := make([]int, len(idsV))
	for i, v := range idsV {
		val, ok := v.(float64)
		if !ok {
			return data.RSPBadRequest
		}
		ids[i] = int(val)
	}

	// Get data
	db, res := a.GetDB(gld)
	if !res.Exists {
		return data.RSPError(res.Message)
	}

	out := make(map[string]any)
	for _, id := range ids {
		elem, res := db.GetElement(id)
		if !res.Exists {
			return data.RSPError(res.Message)
		}
		val, err := json.Marshal(elem)
		if err != nil {
			return data.RSPError(err.Error())
		}
		out[strconv.Itoa(id)] = string(val)
	}

	// NOTE: Data is string, json marshalled
	return data.RSPSuccess(out)
}

func (a *API) MethodCombo(params map[string]any, id, gld string) data.Response {
	// Process params
	vals, ok := params["elems"]
	if !ok {
		return data.RSPBadRequest
	}
	v, ok := vals.([]interface{})
	if !ok {
		return data.RSPBadRequest
	}
	elems := make([]int, len(v))
	for i, el := range v {
		v, ok := el.(float64)
		if !ok {
			return data.RSPBadRequest
		}
		elems[i] = int(v)
	}

	// Get data
	db, res := a.GetDB(gld)
	if !res.Exists {
		return data.RSPError(res.Message)
	}

	// Check if you have everything
	inv := db.GetInv(id)
	for _, el := range elems {
		if !inv.Contains(el) {
			el, res := db.GetElement(el)
			if !res.Exists {
				return data.RSPError(res.Message)
			}
			return data.RSPError(fmt.Sprintf("You don't have %s!", el.Name))
		}
	}

	el3, res := db.GetCombo(elems)
	if !res.Exists {
		return data.RSPError(res.Message)
	}

	// Save to inv
	exists := inv.Contains(el3)
	if !exists {
		inv.Add(el3)
		err := db.SaveInv(inv)
		if err != nil {
			return data.RSPError(err.Error())
		}
	}

	return data.RSPSuccess(map[string]any{"id": el3, "exists": exists})
}

func (a *API) MethodInv(id, gld string) data.Response {
	// Get data
	db, res := a.GetDB(gld)
	if !res.Exists {
		return data.RSPError(res.Message)
	}
	inv := db.GetInv(id)
	els := make([]int, len(inv.Elements))
	i := 0
	inv.Lock.RLock()
	for k := range inv.Elements {
		els[i] = k
		i++
	}
	inv.Lock.RUnlock()
	return data.RSPSuccess(map[string]any{"elems": els})
}
