package elements

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
	"github.com/lib/pq"
)

// Format: user|query|offset
func (e *Elements) NextHandler(c sevcord.Ctx, params string) {
	parts := strings.Split(params, "|")
	if c.Author().User.ID != parts[0] {
		c.Acknowledge()
		c.Respond(sevcord.NewMessage("You are not authorized! " + types.RedCircle))
		return
	}
	offset, _ := strconv.Atoi(parts[2])

	// Get element to make
	qu := parts[1]
	var err error
	var res int
	var elem types.Element
	if qu == "" { // No query
		err = e.db.QueryRow(`SELECT c.result FROM combos c JOIN inventories i ON c.els <@ i.inv 
				    WHERE i."user"=$2 AND i.guild=$1 AND c.guild=$1 AND NOT (c.result = ANY(i.inv)) LIMIT 1 OFFSET $3`, c.Guild(), c.Author().User.ID, offset).Scan(&res)
	} else { // With query
		query, ok := e.base.CalcQuery(c, qu)
		if !ok {
			return
		}
		err = e.db.QueryRow(`SELECT c.result FROM combos c JOIN inventories i ON c.els <@ i.inv 
				    WHERE i."user"=$2 AND i.guild=$1 AND c.guild=$1 AND NOT (c.result = ANY(i.inv)) AND c.result=ANY($4) LIMIT 1 OFFSET $3`, c.Guild(), c.Author().User.ID, offset, pq.Array(query.Elements)).Scan(&res)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			c.Respond(sevcord.NewMessage("Nothing to do next found! Try again later. " + types.RedCircle))
			return
		} else {
			e.base.Error(c, err)
			return
		}
	}

	//Get element for thumbnail
	err = e.db.Get(&elem, "SELECT * FROM elements WHERE id=$1 AND guild=$2", res, c.Guild())
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Get combos
	var items []struct {
		Els pq.Int32Array `db:"els"`
	}
	err = e.db.Select(&items, `WITH els as MATERIALIZED(SELECT els FROM combos WHERE result=$1 AND guild=$2)
	SELECT els FROM els WHERE els <@ (SELECT inv FROM inventories WHERE "user"=$3 AND guild=$2)`, res, c.Guild(), c.Author().User.ID) // TODO: Figure out why MATERIALIZED is needed - makes really bad query plan if not included
	if err != nil {
		e.base.Error(c, err)
		return
	}
	itemCnt := len(items)
	if len(items) > maxHintEls {
		items = items[:maxHintEls]
	}

	// Get names
	ids := []int32{int32(res)}
	for _, item := range items {
		ids = append(ids, item.Els...)
	}
	nameMap, err := e.base.NameMap(util.Map(ids, func(a int32) int { return int(a) }), c.Guild())
	if err != nil {
		e.base.Error(c, err)
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

	params = fmt.Sprintf("%s|%s|%d", parts[0], parts[1], offset+1)
	emb := sevcord.NewEmbed().
		Title("Your next element is "+nameMap[int(res)]).
		Description(desc.String()).
		Color(elem.Color).
		Footer(fmt.Sprintf("%s Combos • Element #%d", humanize.Comma(int64(itemCnt)), elem.ID), "")
	if elem.Image != "" {
		emb = emb.Thumbnail(elem.Image)
	}
	c.Respond(sevcord.NewMessage("").
		AddEmbed(emb).
		AddComponentRow(sevcord.NewButton("Next Element", sevcord.ButtonStylePrimary, "next", params).
			WithEmoji(sevcord.ComponentEmojiCustom("next", "1133079167043375204", false))))
}

func (e *Elements) Next(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	query := ""
	if opts[0] != nil {
		query = opts[0].(string)
	}
	e.NextHandler(c, fmt.Sprintf("%s|%s|0", c.Author().User.ID, query))
}
