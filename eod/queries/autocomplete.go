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

func (q *Queries) queryParents(c sevcord.Ctx, name string, res map[string]struct{}) error {
	_, exists := res[name]
	if exists {
		return nil
	}
	qu, err := q.base.CalcQuery(c, name)
	if err != nil {
		return err
	}
	res[qu.Name] = struct{}{}
	switch qu.Kind {
	case types.QueryKindProducts, types.QueryKindParents, types.QueryKindRegex, types.QueryKindComparison:
		return q.queryParents(c, qu.Data["query"].(string), res)

	case types.QueryKindOperation:
		err = q.queryParents(c, qu.Data["left"].(string), res)
		if err != nil {
			return err
		}
		return q.queryParents(c, qu.Data["right"].(string), res)
	}

	return nil
}
