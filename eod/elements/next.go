package elements

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
	"github.com/lib/pq"
)

func (e *Elements) Next(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	// Get element to make
	var res int
	err := e.db.QueryRow(`WITH inv as (SELECT inv FROM inventories WHERE "user"=$2 AND guild=$1)
	SELECT result FROM combos WHERE els <@ (SELECT inv FROM inv) AND NOT (result=ANY(SELECT UNNEST(inv) FROM inv)) AND guild=$1 LIMIT 1`, c.Guild(), c.Author().User.ID).Scan(&res)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Respond(sevcord.NewMessage("Nothing to do next found! Try again later. " + types.RedCircle))
		} else {
			e.base.Error(c, err)
			return
		}
	}

	// Get combos
	var items []struct {
		Els pq.Int32Array `db:"els"`
	}
	err = e.db.Select(&items, `WITH els as MATERIALIZED(SELECT els FROM combos WHERE result=$1 AND guild=$2)
	SELECT els FROM els WHERE els <@ (SELECT inv FROM inventories WHERE "user"=$3 AND guild=$2)`, res, c.Guild(), c.Author().User.ID) // TODO: Figure out why MATERIALIZED is needed - makes really bad query plan if not included
	if err != nil {
		e.base.Error(c, err)
		return
	}
	itemCnt := len(items)
	if len(items) > maxHintEls {
		items = items[:maxHintEls]
	}

	// Get names
	ids := []int32{int32(res)}
	for _, item := range items {
		ids = append(ids, item.Els...)
	}
	nameMap, err := e.base.NameMap(util.Map(ids, func(a int32) int { return int(a) }), c.Guild())
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Result
	desc := &strings.Builder{}
	for _, item := range items {
		for i, el := range item.Els {
			if i > 0 {
				desc.WriteString(" + ")
			}
			name := nameMap[int(el)]
			if i == len(item.Els)-1 {
				name = Obscure(name)
			}
			desc.WriteString(name)
		}
		desc.WriteRune('\n')
	}

	emb := sevcord.NewEmbed().
		Title("Your next element is "+nameMap[int(res)]).
		Description(desc.String()).
		Color(15158332). // Red
		Footer(fmt.Sprintf("%s Combos", humanize.Comma(int64(itemCnt))), "")
	c.Respond(sevcord.NewMessage("").
		AddEmbed(emb))
}
