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

// Format: user|query|count|distinct|cmdtype
func (e *Elements) IdeaHandler(c sevcord.Ctx, params string) {
	parts := strings.Split(params, "|")
	if c.Author().User.ID != parts[0] {
		c.Acknowledge()
		c.Respond(sevcord.NewMessage("You are not authorized! " + types.RedCircle))
		return
	}
	cnt, _ := strconv.Atoi(parts[2])
	distinct := 0
	if len(parts) > 3 {
		distinct, _ = strconv.Atoi(parts[3])
	}

	// Get combo
	var els pq.Int32Array = nil
	foundEl := false
	var query *types.Query
	var ok bool

	if parts[1] != "" {
		query, ok = e.base.CalcQuery(c, parts[1])
		if !ok {
			return
		}
	}
	cmdtype := "idea"
	if len(parts) > 4 {
		cmdtype = parts[4]
	}
	max := MaxIdeaReqs
	if cmdtype == "randomcombo" {
		max = 1
	}
	for i := 0; i < max; i++ {
		els = make(pq.Int32Array, 0)
		// Get elements
		var err error
		if distinct > 0 {
			if parts[1] == "" {
				err = e.db.Select(&els, `SELECT id FROM elements WHERE guild=$1 AND id=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$3) ORDER BY RANDOM() LIMIT $2`, c.Guild(), cnt, c.Author().User.ID)
			} else {

				err = e.db.Select(&els, `SELECT id FROM elements WHERE guild=$1 AND id=ANY($2) AND id=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$4) ORDER BY RANDOM() LIMIT $3`, c.Guild(), pq.Array(query.Elements), cnt, c.Author().User.ID)
			}
		} else {
			if parts[1] == "" {

				for i := 0; i < cnt; i++ {
					var el int32
					err = e.db.Get(&el, `SELECT id FROM elements WHERE guild=$1 AND id=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$2) ORDER BY RANDOM()`, c.Guild(), c.Author().User.ID)
					if err != nil {
						e.base.Error(c, err)
						return
					}
					els = append(els, el)
				}
			} else {
				for i := 0; i < cnt; i++ {
					var el int32
					err = e.db.Get(&el, `SELECT id FROM elements WHERE guild=$1 AND id=ANY($2) AND id=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$3) ORDER BY RANDOM()`, c.Guild(), pq.Array(query.Elements), c.Author().User.ID)
					if err != nil {
						e.base.Error(c, err)
						return
					}
					els = append(els, el)
				}
			}

		}

		if err != nil {
			e.base.Error(c, err)
			return
		}
		if len(els) < cnt {
			continue
		}
		sort.Slice(els, func(i, j int) bool {
			return els[i] < els[j]
		})

		// Check combo
		var exists bool
		err = e.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM combos WHERE guild=$1 AND els=$2)`, c.Guild(), els).Scan(&exists)
		if err != nil {
			e.base.Error(c, err)
		}
		if !exists {
			foundEl = true
			break
		}

	}
	if !foundEl && cmdtype != "randomcombo" {
		c.Respond(sevcord.NewMessage("No ideas found! Try again later. " + types.RedCircle))
		return
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

	// Respond
	if cmdtype == "randomcombo" {
		var msgtext string
		if foundEl {
			msgtext = fmt.Sprintf("Your random combination is... **%s**\n\tSuggest a result using </suggest:%s>", elDesc.String(), suggestCmdId)
			e.base.SaveCombCache(c, types.CombCache{Elements: util.Map(els, func(a int32) int { return int(a) }), Result: -1})
		} else {
			have := false
			res := 0
			err = e.db.QueryRow(`SELECT result FROM combos WHERE guild=$1 AND result <@ $2 AND result @> $2`, c.Guild(), pq.Array(els)).Scan(&res)
			if err != nil {
				return
			}
			resname, _ := e.base.GetName(c.Guild(), res)
			e.db.QueryRow(`SELECT EXISTS(SELECT id FROM elements WHERE guild=$1 AND id=$2 AND id=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$3))`, c.Guild(), res, c.Author().User.ID).Scan(&have)
			if !have {
				//add element to inv
				_, err := e.db.Exec(`UPDATE inventories SET inv=array_append(inv, $3) WHERE guild=$1 AND "user"=$2`, c.Guild(), c.Author().User.ID, res)
				if err != nil {
					e.base.Error(c, err)
					return
				}

				msgtext = fmt.Sprintf("Your random combination is... **%s**\n\tYou made **%s** 🆕", elDesc.String(), resname)
			} else {
				msgtext = fmt.Sprintf("Your random combination is... **%s**\n\tYou made **%s**, but already have it. 🔵", elDesc.String(), resname)
			}
			e.base.SaveCombCache(c, types.CombCache{Elements: util.Map(els, func(a int32) int { return int(a) }), Result: res})
		}
		c.Respond(sevcord.NewMessage(msgtext).
			AddComponentRow(sevcord.NewButton("New Random Combination", sevcord.ButtonStylePrimary, "randcombo", params).
				WithEmoji(sevcord.ComponentEmojiDefault('🎲')),
			))
	} else {
		// Update comb cache
		e.base.SaveCombCache(c, types.CombCache{Elements: util.Map(els, func(a int32) int { return int(a) }), Result: -1})

		c.Respond(
			sevcord.NewMessage(fmt.Sprintf("Your random unused combination is... **%s**\n\tSuggest a result using </suggest:%s>", elDesc.String(), suggestCmdId)).
				AddComponentRow(sevcord.NewButton("New Idea", sevcord.ButtonStylePrimary, "idea", params).
					WithEmoji(sevcord.ComponentEmojiCustom("idea", "932832178847502386", false))),
		)
	}

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
	distinctval := 1
	if len(opts) > 1 && opts[2] != nil && !opts[2].(bool) {
		distinctval = 0
	}
	e.IdeaHandler(c, fmt.Sprintf("%s|%s|%d|%d|idea", c.Author().User.ID, query, cnt, distinctval))
}
func (e *Elements) RandomCombo(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	query := ""
	if opts[0] != nil {
		query = opts[0].(string)
	}
	cnt := 2
	if opts[1] != nil {
		cnt = int(opts[1].(int64))
	}
	distinctval := 0
	if len(opts) > 1 && opts[2] != nil && opts[2].(bool) {
		distinctval = 1
	}
	e.IdeaHandler(c, fmt.Sprintf("%s|%s|%d|%d|randomcombo", c.Author().User.ID, query, cnt, distinctval))
}
