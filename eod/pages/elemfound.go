package pages

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
)

// Params: prevnext|elem|page
func (p *Pages) ElemFoundHandler(c sevcord.Ctx, params string) {
	parts := strings.Split(params, "|")
	elem, _ := strconv.Atoi(parts[1])

	// Get count
	var cnt int
	err := p.db.QueryRow(`SELECT COUNT(*) FROM inventories WHERE guild=$1 AND $2=ANY(inv)`, c.Guild(), elem).Scan(&cnt)
	if err != nil {
		p.base.Error(c, err)
		return
	}
	length := p.base.PageLength(c)
	pagecnt := int(math.Ceil(float64(cnt) / float64(length)))

	// Apply page
	page, _ := strconv.Atoi(parts[2])
	page = ApplyPage(parts[0], page, pagecnt)

	// Get values
	var found []string
	err = p.db.Select(&found, `SELECT "user" FROM inventories WHERE $2=ANY(inv) AND guild=$1 ORDER BY array_length(inv, 1) DESC LIMIT $3 OFFSET $4`, c.Guild(), elem, length, length*page)
	if err != nil {
		p.base.Error(c, err)
		return
	}

	// Make text
	desc := &strings.Builder{}
	for _, v := range found {
		fmt.Fprintf(desc, "<@%s>\n", v)
	}

	// Get elem
	name, err := p.base.GetName(c.Guild(), elem)
	if err != nil {
		p.base.Error(c, err)
		return
	}

	// Create
	embed := sevcord.NewEmbed().
		Title(fmt.Sprintf("%s's Found (%s)", name, humanize.Comma(int64(cnt)))).
		Description(desc.String()).
		Footer(fmt.Sprintf("Page %d/%d", page+1, pagecnt), "").
		Color(15277667) // Pink

	c.Respond(sevcord.NewMessage("").
		AddEmbed(embed).
		AddComponentRow(PageSwitchBtns("elemfound", fmt.Sprintf("%d|%d", elem, page))...),
	)
}

func (p *Pages) ElemFound(c sevcord.Ctx, args []any) {
	c.Acknowledge()

	page := -1
	if len(args) > 1 && args[1] != nil {
		page = int(args[1].(int64)) - 2
	}
	// Create embed
	p.ElemFoundHandler(c, fmt.Sprintf("next|%d|%d", args[0].(int64), page))
}
