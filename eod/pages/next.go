package pages

import (
	"database/sql"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
	"github.com/lib/pq"
)

// Format: prevnext|user|query|offset|page
func (p *Pages) NextHandler(c sevcord.Ctx, params string) {
	parts := strings.Split(params, "|")
	if c.Author().User.ID != parts[1] {
		c.Acknowledge()
		c.Respond(sevcord.NewMessage("You are not authorized! " + types.RedCircle))
		return
	}
	if len(parts) != 5 {
		return
	}
	offset, _ := strconv.Atoi(parts[3])

	// Get element to make
	qu := parts[2]
	var err error
	var res int
	var elem types.Element
	if qu == "" { // No query
		err = p.db.QueryRow(`SELECT c.result FROM combos c JOIN inventories i ON c.els <@ i.inv 
				    WHERE i."user"=$2 AND i.guild=$1 AND c.guild=$1 AND NOT (c.result = ANY(i.inv)) LIMIT 1 OFFSET $3`, c.Guild(), c.Author().User.ID, offset).Scan(&res)
	} else { // With query
		query, ok := p.base.CalcQuery(c, qu)
		if !ok {
			return
		}
		err = p.db.QueryRow(`SELECT c.result FROM combos c JOIN inventories i ON c.els <@ i.inv 
				    WHERE i."user"=$2 AND i.guild=$1 AND c.guild=$1 AND NOT (c.result = ANY(i.inv)) AND c.result=ANY($4) LIMIT 1 OFFSET $3`, c.Guild(), c.Author().User.ID, offset, pq.Array(query.Elements)).Scan(&res)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			c.Respond(sevcord.NewMessage("Nothing to do next found! Try again later. " + types.RedCircle))
			return
		} else {
			p.base.Error(c, err)
			return
		}
	}

	//Get element for thumbnail
	err = p.db.Get(&elem, "SELECT * FROM elements WHERE id=$1 AND guild=$2", res, c.Guild())
	if err != nil {
		p.base.Error(c, err)
		return
	}

	// Get combos
	var items []struct {
		Els pq.Int32Array `db:"els"`
	}
	err = p.db.Select(&items, `WITH els as MATERIALIZED(SELECT els FROM combos WHERE result=$1 AND guild=$2)
	SELECT els FROM els WHERE els <@ (SELECT inv FROM inventories WHERE "user"=$3 AND guild=$2)`, res, c.Guild(), c.Author().User.ID) // TODO: Figure out why MATERIALIZED is needed - makes really bad query plan if not included
	if err != nil {
		p.base.Error(c, err)
		return
	}
	maxHintEls := p.base.PageLength(c)
	pagecnt := int(math.Ceil(float64(len(items)) / float64(maxHintEls)))
	page, _ := strconv.Atoi(parts[4])
	page = ApplyPage(parts[0], page, pagecnt)
	itemCnt := len(items)
	if len(items) > maxHintEls {
		max := math.Min(float64(maxHintEls+page*maxHintEls), float64(len(items)))
		items = items[page*maxHintEls : int(max)]
	}

	// Get names
	ids := []int32{int32(res)}
	for _, item := range items {
		ids = append(ids, item.Els...)
	}
	nameMap, err := p.base.NameMap(util.Map(ids, func(a int32) int { return int(a) }), c.Guild())
	if err != nil {
		p.base.Error(c, err)
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
	pgtext := ""
	pgtext = fmt.Sprintf("Page %d/%d • ", page+1, pagecnt)

	params = fmt.Sprintf("next|%s|%s|%d|-1", parts[1], parts[2], offset+1)
	emb := sevcord.NewEmbed().
		Title("Your next element is "+nameMap[int(res)]).
		Description(desc.String()).
		Color(elem.Color).
		Footer(fmt.Sprintf("%s%s Combos", pgtext, humanize.Comma(int64(itemCnt))), "")
	if elem.Image != "" {
		emb = emb.Thumbnail(elem.Image)
	}
	comps := make([]sevcord.Component, 0)
	comps = append(comps, sevcord.NewButton("Next Element", sevcord.ButtonStylePrimary, "next", params).
		WithEmoji(sevcord.ComponentEmojiCustom("next", "1133079167043375204", false)))
	comps = append(comps, PageSwitchBtns("next", fmt.Sprintf("%s|%s|%s|%d", parts[1], parts[2], parts[3], page))...)
	c.Respond(sevcord.NewMessage("").
		AddEmbed(emb).
		AddComponentRow(comps...))
}

func (p *Pages) Next(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	query := ""
	if opts[0] != nil {
		query = opts[0].(string)
	}
	page := -1
	if len(opts) > 2 && opts[1] != nil {
		page = (int)(opts[1].(int64) - 1)
	}
	p.NextHandler(c, fmt.Sprintf("next|%s|%s|0|%d", c.Author().User.ID, query, page))
}
