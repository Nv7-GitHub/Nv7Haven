package elements

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
	"github.com/lib/pq"
)

const maxHintEls = 30

var noObscure = map[rune]struct{}{
	' ': {},
	'.': {},
	'-': {},
	'_': {},
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

// Format: user|elementid|query

func (e *Elements) hintHandler(c sevcord.Ctx, params string) {
	parts := strings.Split(params, "|")
	if c.Author().User.ID != parts[0] {
		c.Respond(sevcord.NewMessage("You are not authorized! " + types.RedCircle))
		return
	}
	elVal, err := strconv.Atoi(parts[1])
	if err != nil {
		e.base.Error(c, err)
		return
	}
	query := parts[2]

	// Get element
	var el int
	if elVal != -1 {
		el = elVal
	} else {
		// Pick random element
		var err error
		if query == "" { // Not from a query
			err = e.db.QueryRow(`SELECT result FROM combos WHERE 
		guild=$1 AND 
		NOT (result=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$2))
		ORDER BY RANDOM()
		LIMIT 1`, c.Guild(), c.Author().User.ID).Scan(&el)
		} else { // From a query
			var qu *types.Query
			var ok bool
			qu, ok = e.base.CalcQuery(c, query)
			if !ok {
				return
			}
			err = e.db.QueryRow(`SELECT result FROM combos WHERE 
		guild=$1 AND 
		NOT (result=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$2)) AND
		result=ANY($3)
		ORDER BY RANDOM()
		LIMIT 1`, c.Guild(), c.Author().User.ID, pq.Array(qu.Elements)).Scan(&el)
		}

		// Get random element that the user can make
		if err != nil {
			if err == sql.ErrNoRows {
				c.Respond(sevcord.NewMessage("No hints found! Try again later. " + types.RedCircle))
			} else {
				e.base.Error(c, err)
			}
			return
		}
	}

	// Check if you have
	var has bool
	err = e.db.QueryRow(`SELECT $3=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$2)`, c.Guild(), c.Author().User.ID, el).Scan(&has)
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Get hint
	var items []struct {
		Els  pq.Int32Array `db:"els"`
		Cont bool          `db:"cont"` // Whether user can make it
	}
	err = e.db.Select(&items, `SELECT els, els <@ (SELECT inv FROM inventories WHERE guild=$1 AND "user"=$2 LIMIT 1) cont FROM combos WHERE guild=$1 AND result=$3`, c.Guild(), c.Author().User.ID, el)
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Sort & limit
	sort.Slice(items, func(i, j int) bool {
		if items[i].Cont && !items[j].Cont {
			return true
		}
		return false
	})
	itemCnt := len(items)
	if len(items) > maxHintEls {
		items = items[:maxHintEls]
	}

	// Get names
	ids := []int32{int32(el)}
	for _, item := range items {
		ids = append(ids, item.Els...)
	}
	nameMap, err := e.base.NameMap(util.Map(ids, func(a int32) int { return int(a) }), c.Guild())
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Create message
	description := &strings.Builder{}
	for _, item := range items {
		// Emoji
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
	emb := sevcord.NewEmbed().
		Title("Hints for "+nameMap[int(el)]).
		Description(description.String()).
		Color(3447003). // Blue
		Footer(fmt.Sprintf("%s Hints • You%s have this", humanize.Comma(int64(itemCnt)), dontHave), "")
	err = c.Respond(sevcord.NewMessage("").
		AddEmbed(emb).
		AddComponentRow(sevcord.NewButton("New Hint", sevcord.ButtonStylePrimary, "hint", params).
			WithEmoji(sevcord.ComponentEmojiCustom("hint", "932833472396025908", false))))
	fmt.Println(err)
}

func (e *Elements) Hint(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	el := -1
	if opts[0] != nil {
		el = int(opts[0].(int64))
	}
	query := ""
	if opts[1] != nil {
		query = opts[1].(string)
	}
	e.hintHandler(c, fmt.Sprintf("%s|%d|%s", c.Author().User.ID, el, query))
}
