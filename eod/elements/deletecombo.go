package elements

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
	"github.com/lib/pq"
)

// Format: user|query index|query
func (e *Elements) UncheeseHandler(c sevcord.Ctx, params string) {
	parts := strings.Split(params, "|")
	if c.Author().User.ID != parts[0] {
		c.Acknowledge()
		c.Respond(sevcord.NewMessage("You are not authorized! " + types.RedCircle))
		return
	}
	ind, err := strconv.Atoi(parts[1])
	if err != nil {
		e.base.Error(c, err)
		return
	}
	query := parts[2]

	// Calculate query
	qu, ok := e.base.CalcQuery(c, query)
	if !ok {
		return
	}

	// Delete combo for previous element
	if ind > 0 {
		var prevParents pq.Int32Array
		err = e.db.QueryRow("SELECT parents FROM elements WHERE guild=$1 AND id=$2", c.Guild(), qu.Elements[ind-1]).Scan(&prevParents)
		if err != nil {
			e.base.Error(c, err)
			return
		}
		_, err = e.db.Exec("DELETE FROM combos WHERE result=$1 AND guild=$2 AND els=$3", qu.Elements[ind-1], c.Guild(), prevParents)
		if err != nil {
			e.base.Error(c, err)
			return
		}

		// Get combo with lowest tree size
		var combos []struct {
			Elements pq.Int32Array `db:"els"`
		}
		err = e.db.Select(&combos, "SELECT els FROM combos WHERE result=$1 AND guild=$2", qu.Elements[ind-1], c.Guild())
		if err != nil {
			e.base.Error(c, err)
			return
		}

		// Find minimum tree size and combo
		var min int
		var minind int
		for i, combo := range combos {
			var treesize int
			err = e.db.QueryRow(`WITH RECURSIVE parents(els, id) AS (
			VALUES($2::integer[], 0)
	 	UNION
			(SELECT b.parents els, b.id id FROM elements b INNER JOIN parents p ON b.id=ANY(p.els) where guild=$1)
	 	) SELECT COUNT(*) FROM parents WHERE id>0`, c.Guild(), combo.Elements).Scan(&treesize)
			if err != nil {
				e.base.Error(c, err)
				return
			}
			if i == 0 || treesize < min {
				min = treesize
				minind = i
			}
		}

		// Update combo
		_, err = e.db.Exec("UPDATE elements SET parents=$1, treesize=$2 WHERE id=$3 AND guild=$4", combos[minind].Elements, min, qu.Elements[ind-1], c.Guild())
		if err != nil {
			e.base.Error(c, err)
			return
		}
	}

	if ind == len(qu.Elements) {
		c.Respond(sevcord.NewMessage("Finished uncheesing!"))
		return
	}

	// Get combo for element and see if its ok
	var parents pq.Int32Array
	err = e.db.QueryRow("SELECT parents FROM elements WHERE guild=$1 AND id=$2", c.Guild(), qu.Elements[ind]).Scan(&parents)
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Print combo out
	names, err := e.base.GetNames(util.Map(parents, func(a int32) int { return int(a) }), c.Guild())
	if err != nil {
		e.base.Error(c, err)
		return
	}
	combo := strings.Join(names, " + ")

	// Send embed
	params = fmt.Sprintf("%s|%d|%s", parts[0], ind+1, parts[2])
	emb := sevcord.NewEmbed().
		Title("Uncheese elements in "+qu.Name).
		Description(combo).
		Color(3447003). // Blue
		Footer(fmt.Sprintf("%s deleted out of %s", humanize.Comma(int64(ind)), humanize.Comma(int64(len(qu.Elements)))), "")
	c.Respond(sevcord.NewMessage("").
		AddEmbed(emb).
		AddComponentRow(sevcord.NewButton("Delete Combo", sevcord.ButtonStyleDanger, "uncheese", params).WithEmoji(sevcord.ComponentEmojiDefault([]rune("🗑️")[0]))))
}

func (e *Elements) Uncheese(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	// Calculate query, make sure that each element in the query has more than 1 combo
	qu, ok := e.base.CalcQuery(c, opts[0].(string))
	if !ok {
		return
	}
	var cnt int
	var minid int
	err := e.db.QueryRow("SELECT id, val FROM (SELECT (SELECT COUNT(*) FROM combos WHERE result=id) val, id FROM elements WHERE id=ANY($2) AND guild=$1) AS abc ORDER BY val ASC LIMIT 1", c.Guild(), pq.Int32Array(util.Map(qu.Elements, func(a int) int32 { return int32(a) }))).Scan(&minid, &cnt)
	if err != nil {
		e.base.Error(c, err)
		return
	}
	if cnt <= 1 {
		name, err := e.base.GetName(c.Guild(), minid)
		if err != nil {
			e.base.Error(c, err)
			return
		}
		c.Respond(sevcord.NewMessage("Element " + name + " has only 1 combo!"))
	}

	// Respond
	e.UncheeseHandler(c, fmt.Sprintf("%s|0|%s", c.Author().User.ID, opts[0].(string)))
}
