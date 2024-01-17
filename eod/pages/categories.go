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

var catListSorts = []sevcord.Choice{
	sevcord.NewChoice("Name", "name"),
	sevcord.NewChoice("Element Count", "count"),
	sevcord.NewChoice("Found", "found"),
}

var catListSortSql = map[string]string{
	"name":  "name",
	"count": "array_length(elements, 1) DESC",
	"found": "common DESC",
}

// Format: prevnext|user|sort|page
func (p *Pages) CatListHandler(c sevcord.Ctx, params string) {
	parts := strings.Split(params, "|")

	// Get count
	var cnt int
	err := p.db.QueryRow(`SELECT COUNT(*) FROM categories WHERE guild=$1`, c.Guild()).Scan(&cnt)
	if err != nil {
		p.base.Error(c, err)
		return
	}
	length := p.base.PageLength(c)
	pagecnt := int(math.Ceil(float64(cnt) / float64(length)))

	// Apply pages
	page, _ := strconv.Atoi(parts[3])
	page = ApplyPage(parts[0], page, pagecnt)

	// Get values
	var cats []struct {
		Name         string `db:"name"`
		Length       int    `db:"length"`
		InvIntersect int    `db:"common"` // # of elements both in inv and cat
	}
	err = p.db.Select(&cats, `SELECT name, array_length(elements, 1) length, COALESCE(array_length(elements & (SELECT inv FROM inventories WHERE guild=$1 AND "user"=$4), 1), 0) common FROM categories WHERE guild=$1 ORDER BY `+catListSortSql[parts[2]]+` LIMIT $2 OFFSET $3`, c.Guild(), length, length*page, parts[1])
	if err != nil {
		p.base.Error(c, err)
		return
	}

	// Description
	desc := &strings.Builder{}
	for _, v := range cats {
		if v.Length == v.InvIntersect {
			fmt.Fprintf(desc, "%s %s\n", v.Name, types.Check)
		} else {
			fmt.Fprintf(desc, "%s (%s%%)\n", v.Name, humanize.FormatFloat("", float64(v.InvIntersect)/float64(v.Length)*100))
		}
	}

	// Respond
	emb := sevcord.NewEmbed().
		Title(fmt.Sprintf("All Categories (%d)", cnt)).
		Description(desc.String()).
		Footer(fmt.Sprintf("Page %d/%d", page+1, pagecnt), "").
		Color(10181046) // Purple
	c.Respond(sevcord.NewMessage("").AddEmbed(emb).AddComponentRow(PageSwitchBtns("catlist", fmt.Sprintf("%s|%s|%d", parts[1], parts[2], page))...))
}

func (p *Pages) CatList(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	// Params
	sort := "name"
	if opts[0] != nil {
		sort = opts[0].(string)
	}

	// Respond
	p.CatListHandler(c, "next|"+c.Author().User.ID+"|"+sort+"|-1")
}

// Params: prevnext|user|sort|page|cat
func (p *Pages) CatHandler(c sevcord.Ctx, params string) {
	parts := strings.SplitN(params, "|", 5)

	// Get count
	var cnt int
	var common int
	err := p.db.QueryRow(`SELECT array_length(elements, 1), COALESCE(array_length(elements & (SELECT inv FROM inventories WHERE guild=$1 AND "user"=$3), 1), 0) FROM categories WHERE guild=$1 AND LOWER(name)=$2`, c.Guild(), strings.ToLower(parts[4]), parts[1]).Scan(&cnt, &common)
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
	err = p.db.Select(&items, `SELECT name, id=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$5) cont FROM elements WHERE id=ANY(SELECT UNNEST(elements) FROM categories WHERE guild=$1 AND name=$2) AND guild=$1 ORDER BY `+types.SortSql[parts[2]]+` LIMIT $3 OFFSET $4`, c.Guild(), parts[4], length, length*page, parts[1])
	if err != nil {
		p.base.Error(c, err)
		return
	}

	// Description
	desc := &strings.Builder{}
	for _, v := range items {
		if v.Cont {
			fmt.Fprintf(desc, "%s %s\n", v.Name, types.Check)
		} else {
			fmt.Fprintf(desc, "%s %s\n", v.Name, types.NoCheck)
		}
	}

	// Create
	embed := sevcord.NewEmbed().
		Title(fmt.Sprintf("%s (%s, %s%%)", parts[4], humanize.Comma(int64(cnt)), humanize.FormatFloat("", float64(common)/float64(cnt)*100))).
		Description(desc.String()).
		Footer(fmt.Sprintf("Page %d/%d", page+1, pagecnt), "").
		Color(10181046) // Purple

	c.Respond(sevcord.NewMessage("").
		AddEmbed(embed).
		AddComponentRow(PageSwitchBtns("cat", fmt.Sprintf("%s|%s|%d|%s", parts[1], parts[2], page, parts[4]))...),
	)
}

func (p *Pages) Cat(c sevcord.Ctx, args []any) {
	c.Acknowledge()

	// Get params
	sort := "id"
	if args[1] != nil {
		sort = args[1].(string)
	}

	// Get name
	var name string
	err := p.db.QueryRow(`SELECT name FROM categories WHERE guild=$1 AND LOWER(name)=$2`, c.Guild(), strings.ToLower(args[0].(string))).Scan(&name)
	if err != nil {
		p.base.Error(c, err, "Category **"+args[0].(string)+"** doesn't exist!")
		return
	}

	// Create embed
	p.CatHandler(c, fmt.Sprintf("next|%s|%s|-1|%s", c.Author().User.ID, sort, name))
}
