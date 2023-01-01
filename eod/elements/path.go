package elements

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
)

func (e *Elements) PathCmd(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	// Get element
	var name string
	err := e.db.QueryRow("SELECT name FROM elements WHERE id=$1 AND guild=$2", opts[0].(int64), c.Guild()).Scan(&name)
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Get vals
	var els []struct {
		ID      int32         `db:"id"`
		Name    string        `db:"name"`
		Parents pq.Int32Array `db:"parents"`
	}
	err = e.db.Select(&els, `WITH RECURSIVE parents AS (
		(select parents, id from elements where id=$2 and guild=$1)
	UNION
		(SELECT b.parents, b.id FROM elements b INNER JOIN parents p ON b.id=ANY(p.parents) where guild=$1)
	) select id, name, parents FROM elements WHERE id=ANY(SELECT id FROM parents) AND guild=$1`, c.Guild(), opts[0].(int64))
	if err != nil {
		e.base.Error(c, err)
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
	addTree(out, int32(opts[0].(int64)), pars, names, &cnt)

	// Send DM
	dm, err := c.Dg().UserChannelCreate(c.Author().User.ID)
	if err != nil {
		e.base.Error(c, err)
		return
	}
	msg := sevcord.NewMessage(fmt.Sprintf("📄 Path for **%s**:", name)).
		AddFile("path.txt", "text/plain", strings.NewReader(out.String()))
	_, err = c.Dg().ChannelMessageSendComplex(dm.ID, msg.Dg())
	if err != nil {
		e.base.Error(c, err)
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
