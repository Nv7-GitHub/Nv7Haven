package elements

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Format: user|query index|comb index|query|(s for skip or d for delete)
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

	combind, err := strconv.Atoi(parts[2])
	if err != nil {
		e.base.Error(c, err)
		return
	}

	query := parts[3]

	// Calculate query
	qu, ok := e.base.CalcQuery(c, query)
	if !ok {
		return
	}
	//check if done
	if ind >= len(qu.Elements) {
		c.Respond(sevcord.NewMessage("Finished uncheesing!"))
		return
	}

	var combos []struct {
		Elements pq.Int32Array `db:"els"`
	}
	var tx *sqlx.Tx
	tx, err = e.db.Beginx()
	if err != nil {
		tx.Rollback()
		e.base.Error(c, err)
		return
	}
	err = tx.Select(&combos, "SELECT els FROM combos WHERE result=$1 AND guild=$2", qu.Elements[ind], c.Guild())
	if err != nil {
		tx.Rollback()
		e.base.Error(c, err)
		return
	}
	// Print combo out
	names, err := e.base.GetNames(util.Map(combos[combind].Elements, func(a int32) int { return int(a) }), c.Guild())
	if err != nil {
		e.base.Error(c, err)
		return
	}
	res, _ := e.base.GetName(c.Guild(), qu.Elements[ind])
	combo := strings.Join(names, " + ") + " = " + res
	// Delete combo for element
	if parts[4] == "d" {

		var prevParents pq.Int32Array
		err = tx.QueryRow("SELECT parents FROM elements WHERE guild=$1 AND id=$2", c.Guild(), qu.Elements[ind]).Scan(&prevParents)
		if err != nil {
			tx.Rollback()
			e.base.Error(c, err)
			return
		}
		IsMinComb := true
		for i := 0; i < len(combos[combind].Elements); i++ {

			if i >= len(prevParents) {
				IsMinComb = false
				break
			}
			if combos[combind].Elements[i] != prevParents[i] {
				IsMinComb = false
				break

			}
		}

		_, err = tx.Exec("DELETE FROM combos WHERE result=$1 AND guild=$2 AND els=$3", qu.Elements[ind], c.Guild(), combos[combind].Elements)
		if err != nil {
			tx.Rollback()
			e.base.Error(c, err)
			return
		}

		// Update combo if it is lowest tree size combo
		if IsMinComb {
			// Find minimum tree size and combo
			min := -1
			var minind int
			for i, combo := range combos {
				treesize, loop, err := e.base.TreeSize(tx, qu.Elements[ind], util.Map(combo.Elements, func(a int32) int { return int(a) }), c.Guild())
				if err != nil {
					tx.Rollback()
					e.base.Error(c, err)
					return
				}
				if (min == -1 || treesize < min) && !loop {
					min = treesize
					minind = i
				}
			}

			if min == -1 {
				tx.Rollback()
				c.Respond(sevcord.NewMessage("Cannot uncheese! " + types.RedCircle))
				return
			}
			_, err = tx.Exec("UPDATE elements SET parents=$1, treesize=$2 WHERE id=$3 AND guild=$4", combos[minind].Elements, min, qu.Elements[ind], c.Guild())

			ind++
			combind = 0
		}

		if err != nil {
			e.base.Error(c, err)
			return
		}

		// Commit
		err = tx.Commit()
		if err != nil {
			e.base.Error(c, err)
			return
		}
	}

	// Send embed
	var nextindex int
	var nextcombindex int
	if combind+1 < len(combos) {
		nextindex = ind
		nextcombindex = combind + 1
	} else {
		nextindex = ind + 1
		nextcombindex = 0
	}

	comps := make([]sevcord.Component, 0)
	if len(combos) > 1 {
		comps = append(comps, sevcord.NewButton("Delete Combo", sevcord.ButtonStyleDanger, "uncheese", fmt.Sprintf("%s|%d|%d|%s|d", parts[0], ind, combind, parts[3])).WithEmoji(sevcord.ComponentEmojiDefault([]rune("🗑️")[0])))
	}
	comps = append(comps, sevcord.NewButton("Skip", sevcord.ButtonStylePrimary, "uncheese", fmt.Sprintf("%s|%d|%d|%s|s", parts[0], nextindex, nextcombindex, parts[3])).WithEmoji(sevcord.ComponentEmojiDefault('⏭')))
	emb := sevcord.NewEmbed().
		Title("Uncheese elements in "+qu.Name).
		Description(combo).
		Color(15548997). // Red
		Footer(fmt.Sprintf("%s elements of %s • %s combos of %s", humanize.Comma(int64(ind+1)), humanize.Comma(int64(len(qu.Elements))), humanize.Comma(int64(combind+1)), humanize.Comma(int64(len(combos)))), "")
	c.Respond(sevcord.NewMessage("").
		AddEmbed(emb).
		AddComponentRow(comps...))
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
		return
	}

	// Respond
	e.UncheeseHandler(c, fmt.Sprintf("%s|0|0|%s|s", c.Author().User.ID, opts[0].(string)))
}
