package api

import (
	"encoding/json"
	"fmt"

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
	elid, ok := v.(float64)
	if !ok {
		return data.RSPError("Bad request")
	}

	// Get data
	db, res := a.GetDB(gld)
	if !res.Exists {
		return data.RSPError(res.Message)
	}
	elem, res := db.GetElement(int(elid))
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
		v, ok := el.(float64)
		if !ok {
			return data.RSPError("Bad request")
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
