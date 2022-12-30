package queries

import (
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord"
	"github.com/lib/pq"
)

func (q *Queries) CalcQuery(ctx sevcord.Ctx, name string) (*types.Query, error) {
	// Get
	var query *types.Query
	err := q.db.Get(query, "SELECT * FROM queries WHERE LOWER(name)=$1 AND guild=$2", name, ctx.Guild())
	if err != nil {
		return nil, err
	}

	// Calc based on type
	switch query.Kind {
	case types.QueryKindElement:
		query.Elements = []int{int(query.Data["elem"].(float64))}

	case types.QueryKindCategory:
		err = q.db.QueryRow(`SELECT elements FROM categories WHERE name=$1 AND guild=$2`, query.Data["cat"].(string), ctx.Guild()).Scan(pq.Array(&query.Elements))
		if err != nil {
			return nil, err
		}

	case types.QueryKindProducts:
		// Get query elems
		parent, err := q.CalcQuery(ctx, query.Data["query"].(string))
		if err != nil {
			return nil, err
		}
		// Calc
		err = q.db.QueryRow(`SELECT result FROM combos WHERE guild=$1 AND array_length($2 & els)>=1`, ctx.Guild(), pq.Array(parent.Elements)).Scan(pq.Array(&query.Elements))
		if err != nil {
			return nil, err
		}

	case types.QueryKindParents:
		// Get query elems
		parent, err := q.CalcQuery(ctx, query.Data["query"].(string))
		if err != nil {
			return nil, err
		}
		// Calc
		err = q.db.Select(&query.Elements, `WITH RECURSIVE parents AS (
			(select parents, id from elements where id=ANY($2) and guild=$1)
		UNION
			(SELECT b.parents, b.id FROM elements b INNER JOIN parents p ON b.id=ANY(p.parents) where guild=$1)
		) select id FROM parents`, ctx.Guild(), pq.Array(parent.Elements))
		if err != nil {
			return nil, err
		}

	case types.QueryKindInventory:
		err = q.db.QueryRow(`SELECT elements FROM inventories WHERE user=$1 AND guild=$2`, query.Data["user"].(string), ctx.Guild()).Scan(pq.Array(&query.Elements))
		if err != nil {
			return nil, err
		}

	case types.QueryKindElements:
		err = q.db.Select(&query.Elements, `SELECT id FROM elements WHERE guild=$1`, ctx.Guild())
		if err != nil {
			return nil, err
		}

	case types.QueryKindRegex:
		err = q.db.Select(`SELECT id FROM elements WHERE guild=$1 AND name ~ $2`, ctx.Guild(), query.Data["regex"].(string))
		if err != nil {
			return nil, err
		}

	case types.QueryKindComparison:
		// Get query elems
		parent, err := q.CalcQuery(ctx, query.Data["query"].(string))
		if err != nil {
			return nil, err
		}
		// Calc
		var op string
		switch query.Data["op"].(string) {
		case "equal":
			op = "="
		case "notequal":
			op = "!="
		case "greater":
			op = ">"
		case "less":
			op = "<"
		}
		err = q.db.Select(&query.Elements, `SELECT id FROM elements WHERE guild=$1 AND id=ANY($2) AND value `+op+` $3`, ctx.Guild(), pq.Array(parent.Elements), query.Data["value"].(float64))
		if err != nil {
			return nil, err
		}
	}

	// Return
	return query, nil
}
