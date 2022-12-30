package pages

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Nv7-Github/sevcord/v2"
)

var queryListSorts = []sevcord.Choice{
	sevcord.NewChoice("Name", "name"),
	sevcord.NewChoice("Element Count", "count"),
	sevcord.NewChoice("Found", "found"),
}

var queryListSortSql = map[string]string{
	"name":  "name",
	"count": "array_length(elements, 1) DESC",
	"found": "common DESC",
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
