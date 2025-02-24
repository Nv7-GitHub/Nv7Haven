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
func (p *Pages) ElemCatHandler(c sevcord.Ctx, params string) {
	parts := strings.Split(params, "|")
	elem, _ := strconv.Atoi(parts[1])

	// Get count
	var cnt int
	err := p.db.QueryRow(`SELECT COUNT(*) FROM categories WHERE guild=$1 AND $2=ANY(elements)`, c.Guild(), elem).Scan(&cnt)
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
	var cats []string
	err = p.db.Select(&cats, `SELECT name FROM categories WHERE $2=ANY(elements) AND guild=$1 ORDER BY array_length(elements, 1) DESC LIMIT $3 OFFSET $4`, c.Guild(), elem, length, length*page)
	if err != nil {
		p.base.Error(c, err)
		return
	}

	// Get elem
	name, err := p.base.GetName(c.Guild(), elem)
	if err != nil {
		p.base.Error(c, err)
		return
	}

	// Create
	embed := sevcord.NewEmbed().
		Title(fmt.Sprintf("%s's Categories (%s)", name, humanize.Comma(int64(cnt)))).
		Description(strings.Join(cats, "\n")).
		Footer(fmt.Sprintf("Page %d/%d", page+1, pagecnt), "").
		Color(15277667) // Orange

	c.Respond(sevcord.NewMessage("").
		AddEmbed(embed).
		AddComponentRow(PageSwitchBtns("elemcats", fmt.Sprintf("%d|%d", elem, page))...),
	)
}

func (p *Pages) ElemCats(c sevcord.Ctx, args []any) {
	c.Acknowledge()

	// Create embed
	page := -1
	if len(args) > 3 && args[3] != nil {
		page = int(args[3].(int64)) - 2
	}
	p.ElemCatHandler(c, fmt.Sprintf("next|%d|%d", args[0].(int64), page))

}
