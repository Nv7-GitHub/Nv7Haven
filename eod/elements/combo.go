package elements

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
)

const suggestCmdId = "1041173178912878662"

func makeListResp(start, join, end string, vals []string) string {
	if len(vals) == 2 {
		return fmt.Sprintf("%s %s %s %s%s %s", start, vals[0], join, vals[1], end, types.RedCircle)
	} else if len(vals) > 2 {
		return fmt.Sprintf("%s %s, %s %s%s %s", start, strings.Join(vals[:len(vals)-1], ", "), join, vals[len(vals)-1], end, types.RedCircle)
	}
	return ""
}

type comboRes struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Cont bool   `db:"cont"`
}

func (e *Elements) Combine(c sevcord.Ctx, ids []int64) {
	c.Acknowledge()
	e.base.IncrementCommandStat(c, "combine")

	if len(ids) > types.MaxComboLength {
		c.Respond(sevcord.NewMessage(fmt.Sprintf("You can only combine up to %d elements! "+types.RedCircle, types.MaxComboLength)))
		return
	}
	if len(ids) < 2 {
		c.Respond(sevcord.NewMessage("You need to combine at least 2 elements! " + types.RedCircle))
		return
	}
	var res []comboRes
	err := e.db.Select(&res, `SELECT id,name,id=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$2) cont FROM elements WHERE guild=$1 AND id=ANY($3)`, c.Guild(), c.Author().User.ID, pq.Array(ids))
	if err != nil {
		e.base.Error(c, err)
		return
	}
	// See what elements you don't have
	dontHave := make([]string, 0)
	for _, v := range res {
		if !v.Cont {
			dontHave = append(dontHave, "**"+v.Name+"**")
		}
	}
	if len(dontHave) == 1 {
		c.Respond(sevcord.NewMessage(fmt.Sprintf("You don't have **%s**! %s", dontHave[0], types.RedCircle)))
		return
	} else if len(dontHave) > 1 {
		c.Respond(sevcord.NewMessage(makeListResp("You don't have", "or", "!", dontHave)))
		return
	}

	// Get items
	nameMap := make(map[string]int, len(res))
	for _, v := range res {
		nameMap[strings.ToLower(v.Name)] = v.ID
	}
	items := make([]int, len(ids))
	for i := range ids {
		items[i] = int(ids[i])
	}
	sort.Ints(items)

	// Save combcache
	e.base.SaveCombCache(c, types.CombCache{Elements: items, Result: -1})

	// Query
	var result int
	err = e.db.QueryRow(`SELECT result FROM combos WHERE guild=$1 AND els=$2`, c.Guild(), pq.Array(items)).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Respond(sevcord.NewMessage(fmt.Sprintf("That combination doesn't exist! %s\n\tSuggest a result using </suggest:%s>", types.RedCircle, suggestCmdId)))
			return
		}

		e.base.Error(c, err)
		return
	}
	e.base.SaveCombCache(c, types.CombCache{Elements: items, Result: result})

	// Check if in inv & get element
	var cont bool
	err = e.db.QueryRow(`SELECT $3=ANY(inv) cont FROM inventories WHERE guild=$1 AND "user"=$2`, c.Guild(), c.Author().User.ID, result).Scan(&cont)
	if err != nil {
		e.base.Error(c, err)
		return
	}
	var name string
	err = e.db.QueryRow(`SELECT name FROM elements WHERE guild=$1 AND id=$2`, c.Guild(), result).Scan(&name)
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Show result
	if cont {
		c.Respond(sevcord.NewMessage(fmt.Sprintf("You made **%s**, but already have it. 🔵", name)))
	} else {
		// Add to inv
		_, err := e.db.Exec(`UPDATE inventories SET inv=array_append(inv, $3) WHERE guild=$1 AND "user"=$2`, c.Guild(), c.Author().User.ID, result)
		if err != nil {
			e.base.Error(c, err)
			return
		}
		c.Respond(sevcord.NewMessage(fmt.Sprintf("You made **%s** 🆕", name)))
	}
}
