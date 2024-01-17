package pages

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
)

// Params: prevnext|elem|sort|page
func (p *Pages) ProductsHandler(c sevcord.Ctx, params string) {
	parts := strings.Split(params, "|")
	elem, _ := strconv.Atoi(parts[1])

	// Get count
	var cnt int
	err := p.db.QueryRow(`SELECT COUNT(*) FROM combos WHERE guild=$1 AND $2=ANY(els)`, c.Guild(), parts[1]).Scan(&cnt)
	if err != nil {
		p.base.Error(c, err)
		return
	}
	length := p.base.PageLength(c)
	pagecnt := int(math.Ceil(float64(cnt) / float64(length)))

	// Apply page
	page, _ := strconv.Atoi(parts[3])
	page = ApplyPage(parts[0], page, pagecnt)

	// Get values
	var items []struct {
		Name string `db:"name"`
		Cont bool   `db:"cont"`
	}
	err = p.db.Select(&items, `SELECT name, id=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$5) cont FROM elements WHERE guild=$1 AND id IN (SELECT result FROM combos WHERE guild=$1 AND $2=ANY(els)) ORDER BY `+types.SortSql[parts[2]]+` LIMIT $3 OFFSET $4`, c.Guild(), elem, length, length*page, c.Author().User.ID)
	if err != nil {
		p.base.Error(c, err)
		return
	}

	// Make description
	desc := &strings.Builder{}
	for _, v := range items {
		if v.Cont {
			fmt.Fprintf(desc, "%s %s\n", v.Name, types.Check)
		} else {
			fmt.Fprintf(desc, "%s %s\n", v.Name, types.NoCheck)
		}
	}

	// Get elem
	name, err := p.base.GetName(c.Guild(), elem)
	if err != nil {
		p.base.Error(c, err)
		return
	}

	// Create
	embed := sevcord.NewEmbed().
		Title(fmt.Sprintf("Products of %s (%s)", name, humanize.Comma(int64(cnt)))).
		Description(desc.String()).
		Footer(fmt.Sprintf("Page %d/%d", page+1, pagecnt), "").
		Color(15548997) // Fuschia

	c.Respond(sevcord.NewMessage("").
		AddEmbed(embed).
		AddComponentRow(PageSwitchBtns("products", fmt.Sprintf("%s|%s|%d", parts[1], parts[2], page))...),
	)
}

func (p *Pages) Products(c sevcord.Ctx, args []any) {
	c.Acknowledge()

	// Get params
	id := args[0].(int64)
	sort := "id"
	if args[1] != nil {
		sort = args[1].(string)
	}

	// Create embed
	p.ProductsHandler(c, fmt.Sprintf("next|%d|%s|-1", id, sort))
}
