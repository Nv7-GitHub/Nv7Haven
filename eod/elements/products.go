package elements

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
)

type productsItem struct {
	Result int  `db:"result"`
	Cont   bool `db:"cont"`
}

func (e *Elements) Products(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	e.base.IncrementCommandStat(c, "products")

	// Get items
	var items []productsItem
	err := e.db.Select(&items, `SELECT result, result=ANY((SELECT inv FROM inventories WHERE guild=$1 AND "user"=$3 LIMIT 1)::integer[]) cont FROM combos WHERE guild=$1 AND $2=ANY(els)`, c.Guild(), opts[0].(int64), c.Author().User.ID)
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Sort & limit
	sort.Slice(items, func(i, j int) bool {
		if items[i].Cont && !items[j].Cont {
			return true
		}
		return false
	})

	// Get names
	names, err := e.base.GetNames(append(util.Map(items, func(a productsItem) int { return a.Result }), int(opts[0].(int64))), c.Guild())
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Description
	desc := &strings.Builder{}
	for i, item := range items {
		if i > 0 {
			desc.WriteString("\n")
		}
		if item.Cont {
			desc.WriteString(types.Check)
		} else {
			desc.WriteString(types.NoCheck)
		}
		desc.WriteString(" ")
		desc.WriteString(names[i])
	}

	// Respond
	emb := sevcord.NewEmbed().
		Title("Products of "+names[len(names)-1]).
		Description(desc.String()).
		Footer(fmt.Sprintf("%s Products", humanize.Comma(int64(len(items)))), "").
		Color(15548997) // Fuschia
	c.Respond(sevcord.NewMessage("").AddEmbed(emb))
}
