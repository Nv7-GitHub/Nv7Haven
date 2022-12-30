package queries

import (
	"log"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
)

func (q *Queries) Autocomplete(ctx sevcord.Ctx, val any) []sevcord.Choice {
	var res []types.Element
	err := q.db.Select(&res, "SELECT name FROM queries WHERE guild=$1 AND name ILIKE $2 || '%' ORDER BY similarity(name, $2) DESC, name LIMIT 25", ctx.Guild(), val.(string))
	if err != nil {
		log.Println("query autocomplete error", err)
		return nil
	}
	choices := make([]sevcord.Choice, len(res))
	for i, v := range res {
		choices[i] = sevcord.NewChoice(v.Name, v.Name)
	}
	return choices
}
