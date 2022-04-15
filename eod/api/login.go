package api

import (
	"fmt"

	"github.com/Nv7-Github/Nv7Haven/eod/api/data"
)

func (a *API) MethodGuild(vals map[string]any, id string) data.Response {
	// Process params
	gldV, ok := vals["gld"]
	if !ok {
		return data.RSPBadRequest
	}
	gld, ok := gldV.(string)
	if !ok {
		return data.RSPBadRequest
	}

	// Check
	db, res := a.GetDB(gld)
	if !res.Exists {
		return data.RSPError(res.Message)
	}
	db.RLock()
	_, exists := db.Invs()[id]
	db.RUnlock()
	if !exists {
		return data.RSPError(fmt.Sprintf("User %s doesn't have an inventory in guild %s!", id, gld))
	}
	return data.RSPSuccess(map[string]any{})
}
