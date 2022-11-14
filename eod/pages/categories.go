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

var sortSql = map[string]string{
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
	err = p.db.Select(&cats, `SELECT name, array_length(elements, 1) length, COALESCE(array_length(elements & (SELECT inv FROM inventories WHERE guild=$1 AND "user"=$4), 1), 0) common FROM categories WHERE guild=$1 ORDER BY `+sortSql[parts[2]]+` LIMIT $2 OFFSET $3`, c.Guild(), length, length*page, parts[1])
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
		Footer(fmt.Sprintf("Page %d/%d", page+1, pagecnt), "")
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
