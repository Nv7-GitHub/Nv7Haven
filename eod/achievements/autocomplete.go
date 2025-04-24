package achievements

import (
	"log"
	"strconv"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
)

func (e *Achievements) Autocomplete(c sevcord.Ctx, val any) []sevcord.Choice {
	var res []types.Achievement
	err := e.db.Select(&res, "SELECT id, name FROM achievements WHERE guild=$1 AND name ILIKE $2 || '%' ORDER BY similarity(name, $2) DESC, id LIMIT 25", c.Guild(), val.(string))
	if err != nil {
		log.Println("autocomplete error", err)
		return nil
	}
	choices := make([]sevcord.Choice, len(res))
	for i, v := range res {
		choices[i] = sevcord.NewChoice(v.Name, strconv.Itoa(v.ID))
	}
	return choices
}
func (a *Achievements) AutocompleteName(ctx sevcord.Ctx, val any) []sevcord.Choice {
	var res []types.Achievement
	err := a.db.Select(&res, "SELECT name FROM achievements WHERE guild=$1 AND name ILIKE $2 || '%' ORDER BY similarity(name, $2) DESC, name LIMIT 25", ctx.Guild(), val.(string))
	if err != nil {
		log.Println("achievement autocomplete error", err)
		return nil
	}
	choices := make([]sevcord.Choice, len(res))
	for i, v := range res {
		choices[i] = sevcord.NewChoice(v.Name, v.Name)
	}
	return choices
}
