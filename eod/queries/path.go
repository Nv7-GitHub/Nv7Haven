package queries

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
)

func (q *Queries) PathCmd(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	// Get query
	qu, err := q.base.CalcQuery(c, opts[0].(string))
	if err != nil {
		q.base.Error(c, err)
		return
	}

	// Check if every element intersects with the author's inv
	var has bool
	err = q.db.QueryRow(`SELECT COALESCE(array_length($3 & inv, 1), 0) = array_length($3, 1) FROM inventories WHERE guild=$1 AND "user"=$2`, c.Guild(), c.Author().User.ID, pq.Array(qu.Elements)).Scan(&has)
	if err != nil {
		q.base.Error(c, err)
		return
	}
	if !has {
		c.Respond(sevcord.NewMessage("You don't have all the elements in this query! " + types.RedCircle))
		return
	}

	// Get vals
	var els []struct {
		ID      int32         `db:"id"`
		Name    string        `db:"name"`
		Parents pq.Int32Array `db:"parents"`
	}
	err = q.db.Select(&els, `WITH RECURSIVE parents AS (
		(select parents, id from elements where id=ANY($2) and guild=$1)
	UNION
		(SELECT b.parents, b.id FROM elements b INNER JOIN parents p ON b.id=ANY(p.parents) where guild=$1)
	) select id, name, parents FROM elements WHERE id=ANY(SELECT id FROM parents) AND guild=$1`, c.Guild(), pq.Array(qu.Elements))
	if err != nil {
		q.base.Error(c, err)
		return
	}

	// Create maps
	pars := make(map[int32][]int32, len(els))
	names := make(map[int32]string, len(els))
	for _, el := range els {
		pars[el.ID] = []int32(el.Parents)
		names[el.ID] = el.Name
	}

	// Calculate
	cnt := 1
	out := &strings.Builder{}
	for v := range qu.Elements {
		addTree(out, int32(v), pars, names, &cnt)
	}

	// Send DM
	dm, err := c.Dg().UserChannelCreate(c.Author().User.ID)
	if err != nil {
		q.base.Error(c, err)
		return
	}
	msg := sevcord.NewMessage(fmt.Sprintf("📄 Path for **%s**:", qu.Name)).
		AddFile("path.txt", "text/plain", strings.NewReader(out.String()))
	_, err = c.Dg().ChannelMessageSendComplex(dm.ID, msg.Dg())
	if err != nil {
		q.base.Error(c, err)
		return
	}

	// Respond
	c.Respond(sevcord.NewMessage("Sent path in DMs! 📄"))
}

func addTree(val *strings.Builder, id int32, parsMap map[int32][]int32, names map[int32]string, cnt *int) {
	pars, exists := parsMap[id]
	if !exists {
		return
	}
	if len(pars) == 0 {
		return
	}
	for _, par := range pars {
		addTree(val, par, parsMap, names, cnt)
	}

	// Add elem
	combo := ""
	for i, v := range pars {
		if i > 0 {
			combo += " + "
		}
		combo += names[v]
	}
	fmt.Fprintf(val, "%d. %s = %s\n", *cnt, combo, names[id])
	*cnt++

	// Remove from map
	delete(parsMap, id)
}
