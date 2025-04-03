package base

import (
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
)

func pgArrToIntArr(arr pq.Int32Array) []int {
	return util.Map([]int32(arr), func(v int32) int { return int(v) })
}

var comparisonFieldMap = map[string]string{
	"id":        "id",
	"name":      "name",
	"image":     "image",
	"color":     "color",
	"comment":   "comment",
	"creator":   "creator",
	"commenter": "commenter",
	"colorer":   "colorer",
	"imager":    "imager",
	"treesize":  "treesize",
	"length":    "LENGTH(name)",
}

func (b *Base) CalcQuery(ctx sevcord.Ctx, name string) (*types.Query, bool) {
	// Get
	var query = &types.Query{}
	err := b.db.Get(query, "SELECT * FROM queries WHERE LOWER(name)=$1 AND guild=$2", strings.ToLower(name), ctx.Guild())
	if err != nil {
		b.Error(ctx, err, "Query **"+name+"** doesn't exist!")
		return nil, false
	}

	// Calc based on type
	switch query.Kind {
	case types.QueryKindElement:
		query.Elements = []int{int(query.Data["elem"].(float64))}

	case types.QueryKindCategory:
		var els pq.Int32Array
		err = b.db.QueryRow(`SELECT elements FROM categories WHERE name=$1 AND guild=$2`, query.Data["cat"].(string), ctx.Guild()).Scan(&els)
		if err != nil {
			b.Error(ctx, err, "Category **"+query.Data["cat"].(string)+"** doesn't exist!")
			return nil, false
		}
		query.Elements = pgArrToIntArr(els)

	case types.QueryKindProducts:
		// Get query elems
		parent, ok := b.CalcQuery(ctx, query.Data["query"].(string))
		if !ok {
			return nil, false
		}
		// Calc
		err = b.db.Select(&query.Elements, `SELECT DISTINCT(result) FROM combos WHERE guild=$1 AND array_length($2 & els, 1)>=1`, ctx.Guild(), pq.Array(parent.Elements))
		if err != nil {
			b.Error(ctx, err)
			return nil, false
		}

	case types.QueryKindParents:
		// Get query elems
		parent, ok := b.CalcQuery(ctx, query.Data["query"].(string))
		if !ok {
			return nil, false
		}
		// Calc
		err = b.db.Select(&query.Elements, `WITH RECURSIVE parents AS (
			(select parents, id from elements where id=ANY($2) and guild=$1)
		UNION
			(SELECT b.parents, b.id FROM elements b INNER JOIN parents p ON b.id=ANY(p.parents) where guild=$1)
		) select id FROM parents`, ctx.Guild(), pq.Array(parent.Elements))
		if err != nil {
			b.Error(ctx, err)
			return nil, false
		}

	case types.QueryKindInventory:
		var els pq.Int32Array
		err = b.db.QueryRow(`SELECT inv FROM inventories WHERE "user"=$1 AND guild=$2`, query.Data["user"].(string), ctx.Guild()).Scan(&els)
		if err != nil {
			b.Error(ctx, err)
			return nil, false
		}
		query.Elements = pgArrToIntArr(els)

	case types.QueryKindElements:
		var max int
		err = b.db.QueryRow(`SELECT MAX(id) FROM elements WHERE guild=$1`, ctx.Guild()).Scan(&max)
		if err != nil {
			b.Error(ctx, err)
			return nil, false
		}
		query.Elements = make([]int, max)
		for i := range query.Elements {
			query.Elements[i] = i + 1
		}

	case types.QueryKindRegex:
		var parent = &types.Query{}
		err = b.db.Get(parent, "SELECT * FROM queries WHERE LOWER(name)=$1 AND guild=$2", strings.ToLower(query.Data["query"].(string)), ctx.Guild())

		if err != nil || parent.Kind == types.QueryKindElements {
			err = b.db.Select(&query.Elements, `SELECT id FROM elements WHERE guild=$1 AND name ~ $2`, ctx.Guild(), query.Data["regex"].(string))
		} else {
			parent, _ = b.CalcQuery(ctx, parent.Name)
			err = b.db.Select(&query.Elements, `SELECT id FROM elements WHERE guild=$1 AND name ~ $2 AND id=ANY($3)`, ctx.Guild(), query.Data["regex"].(string), pq.Array(parent.Elements))
		}

		if err != nil {
			b.Error(ctx, err)
			return nil, false
		}

	case types.QueryKindComparison:
		// Calc
		var op string
		switch query.Data["typ"].(string) {
		case "equal":
			op = "="
		case "notequal":
			op = "!="
		case "greater":
			op = ">"
		case "less":
			op = "<"
		}
		err = b.db.Select(&query.Elements, `SELECT id FROM elements WHERE guild=$1 AND `+comparisonFieldMap[query.Data["field"].(string)]+op+`$2`, ctx.Guild(), query.Data["value"])
		if err != nil {
			b.Error(ctx, err)
			return nil, false
		}

	case types.QueryKindOperation:
		// Get elems
		left, ok := b.CalcQuery(ctx, query.Data["left"].(string))
		if !ok {
			return nil, false
		}
		right, ok := b.CalcQuery(ctx, query.Data["right"].(string))
		if !ok {
			return nil, false
		}

		// Operate
		var out map[int]struct{}
		switch query.Data["op"].(string) {
		case "union":
			out = make(map[int]struct{}, len(left.Elements)+len(right.Elements))
			for _, elem := range left.Elements {
				out[elem] = struct{}{}
			}
			for _, elem := range right.Elements {
				out[elem] = struct{}{}
			}
		case "difference":
			out = make(map[int]struct{}, len(left.Elements))
			for _, elem := range left.Elements {
				out[elem] = struct{}{}
			}
			for _, elem := range right.Elements {
				delete(out, elem)
			}

		case "intersection":
			rightV := make(map[int]struct{}, len(right.Elements))
			for _, elem := range right.Elements {
				rightV[elem] = struct{}{}
			}
			leftV := make(map[int]struct{}, len(left.Elements))
			for _, elem := range left.Elements {
				leftV[elem] = struct{}{}
			}
			out = make(map[int]struct{}, len(left.Elements))
			for _, elem := range left.Elements {
				if _, ok := rightV[elem]; ok {
					out[elem] = struct{}{}
				}
			}
			for _, elem := range right.Elements {
				if _, ok := leftV[elem]; ok {
					out[elem] = struct{}{}
				}
			}
		}

		// Save
		query.Elements = make([]int, 0, len(out))
		for elem := range out {
			query.Elements = append(query.Elements, elem)
		}
	}

	// Return
	return query, true
}
