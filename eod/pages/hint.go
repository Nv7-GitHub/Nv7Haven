package pages

import (
	"database/sql"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
	"github.com/lib/pq"
)

var noObscure = map[rune]struct{}{
	' ': {},
	'.': {},
	'-': {},
	'_': {},
	'(': {},
	')': {},
	'[': {},
	']': {},
	'{': {},
	'}': {},
}

func Obscure(val string) string {
	out := make([]rune, len([]rune(val)))
	i := 0
	for _, char := range val {
		_, exists := noObscure[char]
		if exists {
			out[i] = char
		} else {
			out[i] = '?'
		}
		i++
	}
	return string(out)
}

const hintQueryRand = `SELECT id FROM elements 
LEFT JOIN (SELECT UNNEST(inv) el FROM inventories WHERE guild=$1 AND "user"=$2) s ON id=el
WHERE 
guild=$1 AND
el IS NULL
%s
LIMIT 1`

const hintQuery = `SELECT id FROM elements WHERE 
guild=$1 AND 
NOT (id=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$2))
%s
%s
LIMIT 1`

// Format: prevnext|user|elementid|query|page
func (p *Pages) HintHandler(c sevcord.Ctx, params string) {
	parts := strings.Split(params, "|")
	if c.Author().User.ID != parts[1] {
		c.Acknowledge()
		c.Respond(sevcord.NewMessage("You are not authorized! " + types.RedCircle))
		return
	}
	if len(parts) != 5 {
		return
	}
	elVal, err := strconv.Atoi(parts[2])
	if err != nil {
		p.base.Error(c, err)
		return
	}
	query := parts[3]
	// Get element
	var el int
	var elem types.Element
	if elVal != -1 {
		el = elVal
	} else {
		// Pick random element
		var err error
		if query == "" { // Not from a query
			err = p.db.QueryRow(fmt.Sprintf(hintQueryRand, "ORDER BY RANDOM()"), c.Guild(), c.Author().User.ID).Scan(&el)
		} else { // From a query
			var qu *types.Query
			var ok bool
			qu, ok = p.base.CalcQuery(c, query)
			if !ok {
				return
			}

			err = p.db.QueryRow(fmt.Sprintf(hintQuery, "AND id=ANY($3)", "AND RANDOM() < 0.01"), c.Guild(), c.Author().User.ID, pq.Array(qu.Elements)).Scan(&el)
			if err == sql.ErrNoRows {
				err = p.db.QueryRow(fmt.Sprintf(hintQuery, "AND id=ANY($3)", "ORDER BY RANDOM()"), c.Guild(), c.Author().User.ID, pq.Array(qu.Elements)).Scan(&el)
			}
		}
		// Get random element that the user can make
		if err != nil {
			if err == sql.ErrNoRows {
				c.Respond(sevcord.NewMessage("No hints found! Try again later. " + types.RedCircle))
			} else {
				p.base.Error(c, err)
			}
			return
		}
	}
	//Get element for thumbnail
	err = p.db.Get(&elem, "SELECT * FROM elements WHERE id=$1 AND guild=$2", el, c.Guild())
	if err != nil {
		p.base.Error(c, err)
		return
	}
	// Check if you have
	var has bool
	err = p.db.QueryRow(`SELECT $3=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$2)`, c.Guild(), c.Author().User.ID, el).Scan(&has)
	if err != nil {
		p.base.Error(c, err)
		return
	}

	// Get hint
	var items []struct {
		Els  pq.Int32Array `db:"els"`
		Cont bool          `db:"cont"` // Whether user can make it
	}
	maxHintEls := p.base.PageLength(c)
	cnt := 0
	p.db.QueryRow(`SELECT COUNT(*) FROM combos WHERE guild=$1 AND result=$2`, c.Guild(), el).Scan(&cnt)
	pagecnt := int(math.Ceil(float64(cnt) / float64(maxHintEls)))

	// Apply page

	page, _ := strconv.Atoi(parts[4])
	page = ApplyPage(parts[0], page, pagecnt)
	err = p.db.Select(&items, `SELECT els, els <@ (SELECT inv FROM inventories WHERE guild=$1 AND "user"=$2 LIMIT 1) cont FROM combos WHERE guild=$1 AND result=$3`, c.Guild(), c.Author().User.ID, el)
	if err != nil {
		p.base.Error(c, err)
		return
	}

	// Sort & limit
	sort.Slice(items, func(i, j int) bool {
		if items[i].Cont && !items[j].Cont {
			return true
		}
		return false
	})
	if len(items) > maxHintEls {
		max := math.Min(float64(maxHintEls+page*maxHintEls), float64(len(items)))
		items = items[page*maxHintEls : int(max)]
	}

	// Get names
	ids := []int32{int32(el)}
	for _, item := range items {
		ids = append(ids, item.Els...)
	}
	nameMap, err := p.base.NameMap(util.Map(ids, func(a int32) int { return int(a) }), c.Guild())
	if err != nil {
		p.base.Error(c, err)
		return
	}
	// Create message
	description := &strings.Builder{}
	for _, item := range items {
		//Emoji
		if item.Cont {
			description.WriteString(types.Check)
		} else {
			description.WriteString(types.NoCheck)
		}
		description.WriteString(" ")

		// Elements
		for i, el := range item.Els {
			if i > 0 {
				description.WriteString(" + ")
			}
			name := nameMap[int(el)]
			if i == len(item.Els)-1 {
				name = Obscure(name)
			}
			description.WriteString(name)
		}

		description.WriteRune('\n')
	}

	// Embed
	dontHave := ""
	if !has {
		dontHave = " don't"
	}
	pgtext := ""
	pgtext = fmt.Sprintf("Page %d/%d • ", page+1, pagecnt)

	emb := sevcord.NewEmbed().
		Title("Hints for "+nameMap[int(el)]).
		Description(description.String()).
		Color(elem.Color).
		Footer(fmt.Sprintf("%s%s Hints • You%s have this", pgtext, humanize.Comma(int64(cnt)), dontHave), "")

	if elem.Image != "" {
		emb = emb.Thumbnail(elem.Image)
	}
	comps := make([]sevcord.Component, 0)
	elemparam := parts[2]
	//the + is added so that switching to a different page preserves the option to get a new random hint
	if strings.HasPrefix(parts[2], "+") {
		elemparam = "-1"
	}
	params = fmt.Sprintf("next|%s|%s|%s|-1", parts[1], elemparam, parts[3])

	comps = append(comps, sevcord.NewButton("New Hint", sevcord.ButtonStylePrimary, "hint", params).WithEmoji(sevcord.ComponentEmojiCustom("hint", "932833472396025908", false)))
	comps = append(comps, PageSwitchBtns("hint", fmt.Sprintf("%s|+%d|%s|%d", parts[1], el, parts[3], page))...)
	c.Respond(sevcord.NewMessage("").
		AddEmbed(emb).
		AddComponentRow(comps...))

}

func (p *Pages) Hint(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	el := -1
	if opts[0] != nil {
		el = int(opts[0].(int64))
	}
	query := ""
	if opts[1] != nil {
		query = opts[1].(string)
	}
	page := -1
	if len(opts) > 2 && opts[2] != nil {
		page = int(opts[2].(int64) - 2)
	}

	p.HintHandler(c, fmt.Sprintf("next|%s|%d|%s|%d", c.Author().User.ID, el, query, page))
}
