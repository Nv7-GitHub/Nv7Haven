package pages

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
	"github.com/lib/pq"
)

var queryListSorts = []sevcord.Choice{
	sevcord.NewChoice("Name", "name"),
}

var queryListSortSql = map[string]string{
	"name": "name",
}

// Format: prevnext|sort|page
func (p *Pages) QueryListHandler(c sevcord.Ctx, params string) {
	parts := strings.Split(params, "|")

	// Get count
	var cnt int
	err := p.db.QueryRow(`SELECT COUNT(*) FROM queries WHERE guild=$1`, c.Guild()).Scan(&cnt)
	if err != nil {
		p.base.Error(c, err)
		return
	}
	length := p.base.PageLength(c)
	pagecnt := int(math.Ceil(float64(cnt) / float64(length)))

	// Apply pages
	page, _ := strconv.Atoi(parts[2])
	page = ApplyPage(parts[0], page, pagecnt)

	// Get values
	var cats []struct {
		Name string `db:"name"`
	}
	err = p.db.Select(&cats, `SELECT name FROM queries WHERE guild=$1 ORDER BY `+queryListSortSql[parts[1]]+` LIMIT $2 OFFSET $3`, c.Guild(), length, length*page)
	if err != nil {
		p.base.Error(c, err)
		return
	}

	// Description
	desc := &strings.Builder{}
	for _, v := range cats {
		desc.WriteString(v.Name + "\n")
	}

	// Respond
	emb := sevcord.NewEmbed().
		Title(fmt.Sprintf("All Queries (%d)", cnt)).
		Description(desc.String()).
		Footer(fmt.Sprintf("Page %d/%d", page+1, pagecnt), "").
		Color(10181046) // Purple
	c.Respond(sevcord.NewMessage("").AddEmbed(emb).AddComponentRow(PageSwitchBtns("querylist", fmt.Sprintf("%s|%d", parts[1], page))...))
}

func (p *Pages) QueryList(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	// Params
	sort := "name"
	if opts[0] != nil {
		sort = opts[0].(string)
	}

	// Respond
	p.QueryListHandler(c, "next|"+sort+"|-1")
}

// Params: prevnext|user|sort|page|query
func (p *Pages) QueryHandler(c sevcord.Ctx, params string) {
	parts := strings.SplitN(params, "|", 5)

	// Get query
	query, ok := p.base.CalcQuery(c, parts[4])
	if !ok {
		return
	}

	// Get count
	cnt := len(query.Elements)
	var common int
	err := p.db.QueryRow(`SELECT COALESCE(array_length($2 & (SELECT inv FROM inventories WHERE guild=$1 AND "user"=$3), 1), 0)`, c.Guild(), pq.Array(query.Elements), parts[1]).Scan(&common)
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
	err = p.db.Select(&items, `SELECT name, id=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$5) cont FROM elements WHERE id=ANY($2) AND guild=$1 ORDER BY `+types.SortSql[parts[2]]+` LIMIT $3 OFFSET $4`, c.Guild(), pq.Array(query.Elements), length, length*page, parts[1])
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
		AddComponentRow(PageSwitchBtns("query", fmt.Sprintf("%s|%s|%d|%s", parts[1], parts[2], page, parts[4]))...),
	)
}

func (p *Pages) Query(c sevcord.Ctx, args []any) {
	c.Acknowledge()

	// Get params
	sort := "id"
	if args[1] != nil {
		sort = args[1].(string)
	}

	// Get name
	var name string
	err := p.db.QueryRow(`SELECT name FROM queries WHERE guild=$1 AND LOWER(name)=$2`, c.Guild(), strings.ToLower(args[0].(string))).Scan(&name)
	if err != nil {
		p.base.Error(c, err)
		return
	}

	// Create embed
	p.QueryHandler(c, fmt.Sprintf("next|%s|%s|-1|%s", c.Author().User.ID, sort, name))
}
