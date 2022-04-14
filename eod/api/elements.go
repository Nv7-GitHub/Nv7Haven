package api

import (
	"encoding/json"

	"github.com/Nv7-Github/Nv7Haven/eod/api/data"
)

func (a *API) MethodElem(params map[string]any, id, gld string) data.Response {
	// Process params
	name, ok := params["name"]
	if !ok {
		return data.RSPError("Bad request")
	}
	nm, ok := name.(string)
	if !ok {
		return data.RSPError("Bad request")
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
	v, ok := params["id"]
	if !ok {
		return data.RSPError("Bad request")
	}
	elid, ok := v.(int)
	if !ok {
		return data.RSPError("Bad request")
	}

	// Get data
	db, res := a.GetDB(gld)
	if !res.Exists {
		return data.RSPError(res.Message)
	}
	elem, res := db.GetElement(elid)
	if !res.Exists {
		return data.RSPError(res.Message)
	}
	val, err := json.Marshal(elem)
	if err != nil {
		return data.RSPError(err.Error())
	}

	// NOTE: Data is string, json marshalled
	return data.RSPSuccess(map[string]any{"data": string(val)})
}

func (a *API) MethodCombo(params map[string]any, id, gld string) data.Response {
	// Process params
	vals, ok := params["elems"]
	if !ok {
		return data.RSPError("Bad request")
	}
	v, ok := vals.([]interface{})
	if !ok {
		return data.RSPError("Bad request")
	}
	elems := make([]int, len(v))
	for i, el := range v {
		elems[i], ok = el.(int)
		if !ok {
			return data.RSPError("Bad request")
		}
	}

	// Get data
	db, res := a.GetDB(gld)
	if !res.Exists {
		return data.RSPError(res.Message)
	}
	el3, res := db.GetCombo(elems)
	if !res.Exists {
		return data.RSPError(res.Message)
	}

	// Save to inv
	inv := db.GetInv(id)
	inv.Add(el3)
	err := db.SaveInv(inv)
	if err != nil {
		return data.RSPError(err.Error())
	}

	return data.RSPSuccess(map[string]any{"id": el3})
}
