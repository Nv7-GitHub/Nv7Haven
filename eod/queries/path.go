package queries

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
)

type parseableElement struct {
	Name      string    `json:"name"`
	Parents   []int32   `json:"parents"`
	Image     string    `json:"image"`
	Color     string    `json:"color"`
	Comment   string    `json:"comment"`
	Creator   string    `json:"creator"`
	Created   time.Time `json:"createdon"`
	Commenter string    `json:"commenter"`
	Colorer   string    `json:"colorer"`
	Imager    string    `json:"imager"`
	TreeSize  int32     `json:"treesize"`
}

func (q *Queries) PathCmd(c sevcord.Ctx, opts []any, text bool) {
	c.Acknowledge()

	// Get query
	qu, ok := q.base.CalcQuery(c, opts[0].(string))
	if !ok {
		return
	}

	// Check if every element intersects with the author's inv
	var has bool
	err := q.db.QueryRow(`SELECT $3 <@ inv FROM inventories WHERE guild=$1 AND "user"=$2`, c.Guild(), c.Author().User.ID, pq.Array(qu.Elements)).Scan(&has)
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
		ID        int32         `db:"id"`
		Name      string        `db:"name"`
		Parents   pq.Int32Array `db:"parents"`
		Image     string        `db:"image"`
		Color     string        `db:"color"`
		Comment   string        `db:"comment"`
		Creator   string        `db:"creator"`
		Created   time.Time     `db:"createdon"`
		Commenter string        `db:"commenter"`
		Colorer   string        `db:"colorer"`
		Imager    string        `db:"imager"`
		TreeSize  int32         `db:"treesize"`
	}
	extraOpts := ""
	if !text {
		extraOpts = ", image, color, comment, creator, createdon, commenter, colorer, imager, treesize"
	}
	err = q.db.Select(&els, fmt.Sprintf(`WITH RECURSIVE parents AS (
		(select parents, id from elements WHERE id=ANY($2) and guild=$1)
	UNION
		(SELECT b.parents, b.id FROM elements b INNER JOIN parents p ON b.id=ANY(p.parents) where guild=$1)
	) select id, name, parents%s FROM elements WHERE id=ANY(SELECT id FROM parents) AND guild=$1`, extraOpts), c.Guild(), pq.Array(qu.Elements))
	if err != nil {
		q.base.Error(c, err)
		return
	}

	// Create maps
	pars := make(map[int32][]int32, len(els))
	info := make(map[int32]parseableElement, len(els))
	for _, el := range els {
		pars[el.ID] = []int32(el.Parents)
		if text {
			info[el.ID] = parseableElement{
				Name: el.Name,
			}
		} else {
			info[el.ID] = parseableElement{
				Name:      el.Name,
				Parents:   []int32(el.Parents),
				Image:     el.Image,
				Color:     el.Color,
				Comment:   el.Comment,
				Creator:   el.Creator,
				Created:   el.Created,
				TreeSize:  el.TreeSize,
				Commenter: el.Commenter,
				Colorer:   el.Colorer,
			}
		}
	}

	// Calculate
	cnt := 1
	var out any
	if text {
		out = &strings.Builder{}
	} else {
		out = make(map[int32]parseableElement)
	}
	for _, v := range qu.Elements {
		addTree(out, int32(v), pars, info, &cnt, text)
	}

	// Make reader
	var outreader io.Reader
	var outlen int
	var name string
	var typ string
	if text {
		outreader = strings.NewReader(out.(*strings.Builder).String())
		name = "path.txt"
		typ = "text/plain"
	} else {
		dat, err := json.Marshal(out)
		if err != nil {
			q.base.Error(c, err)
			return
		}
		outlen = len(dat)
		outreader = bytes.NewReader(dat)
		name = "path.json"
		typ = "application/json"
	}

	// Send DM
	dm, err := c.Dg().UserChannelCreate(c.Author().User.ID)
	if err != nil {
		q.base.Error(c, err)
		return
	}
	msg := sevcord.NewMessage(fmt.Sprintf("📄 Path for **%s**:", qu.Name)).
		AddFile(name, typ, outreader, outlen)
	_, err = c.Dg().ChannelMessageSendComplex(dm.ID, msg.Dg())
	if err != nil {
		q.base.Error(c, err)
		return
	}

	// Respond
	c.Respond(sevcord.NewMessage("Sent path in DMs! 📄"))
}

func addTree(val any, id int32, parsMap map[int32][]int32, info map[int32]parseableElement, cnt *int, parseable bool) {
	pars, exists := parsMap[id]
	if !exists {
		return
	}
	if len(pars) == 0 {
		return
	}
	for _, par := range pars {
		addTree(val, par, parsMap, info, cnt, parseable)
	}

	// Add elem
	if parseable {
		combo := ""
		for i, v := range pars {
			if i > 0 {
				combo += " + "
			}
			combo += info[v].Name
		}
		fmt.Fprintf(val.(*strings.Builder), "%d. %s = %s\n", *cnt, combo, info[id].Name)
		*cnt++
	} else {
		val.(map[int32]parseableElement)[id] = info[id]
	}

	// Remove from map
	delete(parsMap, id)
}
