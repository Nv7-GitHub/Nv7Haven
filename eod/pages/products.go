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

// Params: prevnext|elem|sort|postfix|page
func (p *Pages) ProductsHandler(c sevcord.Ctx, params string) {
	parts := strings.Split(params, "|")
	elem, _ := strconv.Atoi(parts[1])

	if len(parts) != 5 {
		return
	}
	// Get count
	var cnt int
	err := p.db.QueryRow(`SELECT COUNT(DISTINCT(result)) FROM combos WHERE guild=$1 AND $2=ANY(els)`, c.Guild(), parts[1]).Scan(&cnt)
	if err != nil {
		p.base.Error(c, err)
		return
	}
	length := p.base.PageLength(c)
	pagecnt := int(math.Ceil(float64(cnt) / float64(length)))

	// Apply page
	page, _ := strconv.Atoi(parts[4])
	page = ApplyPage(parts[0], page, pagecnt)

	// Get values
	var items []struct {
		Name    string `db:"name"`
		Cont    bool   `db:"cont"`
		Postfix string `db:"postfix"`
	}
	postfix := false
	if parts[3] == "1" {
		postfix = true
	} else {
		postfix = false
	}
	postfixable := parts[2] != "found" && parts[2] != "length"
	if postfixable && postfix {
		err = p.db.Select(&items, `WITH els AS MATERIALIZED(SELECT *, id=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$5) cont FROM elements WHERE guild=$1 AND id IN (SELECT result FROM combos WHERE guild=$1 AND $2=ANY(els))) SELECT name, cont, `+parts[2]+` postfix FROM els ORDER BY `+types.SortSql[parts[2]]+` LIMIT $3 OFFSET $4`, c.Guild(), elem, length, length*page, c.Author().User.ID)
	} else {
		err = p.db.Select(&items, `WITH els AS MATERIALIZED(SELECT *, id=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$5) cont FROM elements WHERE guild=$1 AND id IN (SELECT result FROM combos WHERE guild=$1 AND $2=ANY(els))) SELECT name, cont FROM els ORDER BY `+types.SortSql[parts[2]]+` LIMIT $3 OFFSET $4`, c.Guild(), elem, length, length*page, c.Author().User.ID)
	}

	if err != nil {
		p.base.Error(c, err)
		return
	}

	// Make description
	desc := &strings.Builder{}
	for _, v := range items {

		if v.Cont {
			fmt.Fprintf(desc, "%s %s", v.Name, types.Check)
		} else {
			fmt.Fprintf(desc, "%s %s", v.Name, types.NoCheck)

		}
		if postfix && parts[2] != "found" {
			desc.WriteString(p.PrintPostfix(parts[2], v.Name, v.Postfix))
		}
		desc.WriteString("\n")

	}

	// Get elem
	name, err := p.base.GetName(c.Guild(), elem)
	if err != nil {
		p.base.Error(c, err)
		return
	}
	var color int
	var img string
	err = p.db.QueryRow("SELECT color,image FROM elements WHERE id=$1 AND guild=$2", elem, c.Guild()).Scan(&color, &img)
	if err != nil {
		p.base.Error(c, err)
		return
	}
	// Create
	embed := sevcord.NewEmbed().
		Title(fmt.Sprintf("Products of %s (%s)", name, humanize.Comma(int64(cnt)))).
		Description(desc.String()).
		Footer(fmt.Sprintf("Page %d/%d", page+1, pagecnt), "").
		Color(color)
	if img != "" {
		embed = embed.Thumbnail(img)
	}

	c.Respond(sevcord.NewMessage("").
		AddEmbed(embed).
		AddComponentRow(PageSwitchBtns("products", fmt.Sprintf("%s|%s|%s|%d", parts[1], parts[2], parts[3], page))...),
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
	postfix := false
	postfixval := 0
	if len(args) < 3 {
		postfixval = 0
	} else if args[2] != nil {
		postfix = args[2].(bool)
	}
	if postfix {
		postfixval = 1
	} else {
		postfixval = 0
	}

	// Create embed
	p.ProductsHandler(c, fmt.Sprintf("next|%d|%s|%d|-1", id, sort, postfixval))
}
