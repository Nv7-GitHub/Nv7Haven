package elements

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
)

const MaxIdeaReqs = 7

// Format: user|query|count
func (e *Elements) IdeaHandler(c sevcord.Ctx, params string) {
	parts := strings.Split(params, "|")
	if c.Author().User.ID != parts[0] {
		c.Acknowledge()
		c.Respond(sevcord.NewMessage("You are not authorized! " + types.RedCircle))
		return
	}
	cnt, _ := strconv.Atoi(parts[2])

	// Get combo
	var els pq.Int32Array = nil
	foundEl := false
	for i := 0; i < MaxIdeaReqs; i++ {
		// Get elements
		var err error
		if parts[1] == "" {
			err = e.db.Select(&els, `SELECT id FROM elements WHERE guild=$1 ORDER BY RANDOM() LIMIT 2`, c.Guild())
		} else {
			query, ok := e.base.CalcQuery(c, parts[1]) // TODO: Move this out of the loop
			if !ok {
				return
			}
			err = e.db.Select(&els, `SELECT id FROM elements WHERE guild=$1 AND id=ANY($2) ORDER BY RANDOM() LIMIT 2`, c.Guild(), pq.Array(query.Elements))
		}
		if err != nil {
			e.base.Error(c, err)
		}
		if len(els) < cnt {
			continue
		}
		sort.Slice(els, func(i, j int) bool {
			return els[i] < els[j]
		})

		// Check combo
		var res bool
		err = e.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM combos WHERE guild=$1 AND els=$2)`, c.Guild(), els).Scan(&res)
		if err != nil {
			e.base.Error(c, err)
		}
		if !res {
			foundEl = true
			break
		}
	}
	if !foundEl {
		c.Respond(sevcord.NewMessage("No ideas found! Try again later. " + types.RedCircle))
	}

	// Format response
	nameMap, err := e.base.NameMap(util.Map(els, func(a int32) int { return int(a) }), c.Guild())
	if err != nil {
		e.base.Error(c, err)
		return
	}
	elDesc := &strings.Builder{}
	for i, el := range els {
		if i > 0 {
			elDesc.WriteString(" + ")
		}
		elDesc.WriteString(nameMap[int(el)])
	}

	// Update comb cache
	e.base.SaveCombCache(c, types.CombCache{Elements: util.Map(els, func(a int32) int { return int(a) }), Result: -1})

	// Respond
	c.Respond(
		sevcord.NewMessage(fmt.Sprintf("Your random unused combination is... **%s**\n\tSuggest a result using </suggest:%s>", elDesc.String(), suggestCmdId)).
			AddComponentRow(sevcord.NewButton("New Idea", sevcord.ButtonStylePrimary, "idea", params).
				WithEmoji(sevcord.ComponentEmojiCustom("idea", "932832178847502386", false))),
	)
}

func (e *Elements) Idea(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	query := ""
	if opts[0] != nil {
		query = opts[0].(string)
	}
	cnt := 2
	if opts[1] != nil {
		cnt = int(opts[1].(int64))
	}
	e.IdeaHandler(c, fmt.Sprintf("%s|%s|%d", c.Author().User.ID, query, cnt))
}
