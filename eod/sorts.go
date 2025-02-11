package eod

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
)

func (b *Bot) combineElements(c sevcord.Ctx, elements []string) {

	ids, ok := b.getElementIds(c, elements)
	if ok {
		b.elements.Combine(c, ids)
	}

}
func (b *Bot) getElementId(c sevcord.Ctx, val string) (int64, bool) {
	ids, ok := b.getElementIds(c, []string{val})
	if len(ids) == 1 {
		return ids[0], ok
	} else {
		return 0, ok
	}

}
func convertVariableID(c sevcord.Ctx, b *Bot, val string) string {

	teststr := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(val, "#")))
	switch teststr {
	case "last":
		cache, _ := b.base.GetCombCache(c)
		if cache.Result != -1 {
			return "#" + fmt.Sprintf("%d", cache.Result)
		}
	case "rand":
		var id int
		err := b.db.QueryRow(`SELECT id FROM elements WHERE guild=$1 ORDER BY RANDOM()`, c.Guild()).Scan(&id)
		if err == nil {
			return "#" + fmt.Sprintf("%d", id)
		}
	case "randinv", "randininv":
		var id int
		err := b.db.QueryRow(`SELECT id FROM elements WHERE guild=$1 AND id=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$2) ORDER BY RANDOM()`, c.Guild(), c.Author().User.ID).Scan(&id)
		if err == nil {
			return "#" + fmt.Sprintf("%d", id)
		}
	}
	return val

}
func (b *Bot) getElementIds(c sevcord.Ctx, vals []string) ([]int64, bool) {

	var ids []int64
	numericIDs := make([]int64, 0)
	names := make([]string, 0)
	var invalid []string

	idmap := make(map[int64]string)
	namemap := make(map[string]int64)
	convert := make(map[string]string)
	for i := 0; i < len(vals); i++ {
		vals[i] = convertVariableID(c, b, vals[i])
		id, ok := IsNumericID(strings.TrimSpace(vals[i]))
		if ok {
			numericIDs = append(numericIDs, id)
		} else {

			names = append(names, strings.TrimSpace(strings.ToLower(vals[i])))
		}
	}
	type nameres struct {
		ID   int64
		Name string
	}
	//get named elements
	if len(names) > 0 {
		var datares []nameres

		err := b.db.Select(&datares, "SELECT id,name FROM elements WHERE LOWER(name)=ANY($1) AND guild=$2", pq.Array(names), c.Guild())
		if err != nil {

		}
		for i := 0; i < len(datares); i++ {

			idmap[datares[i].ID] = datares[i].Name
			namemap[datares[i].Name] = datares[i].ID
			convert[strings.ToLower(datares[i].Name)] = datares[i].Name

		}
		for i := 0; i < len(names); i++ {
			_, ok := namemap[convert[names[i]]]

			if !ok && !slices.Contains(invalid, fmt.Sprintf("**%s**", names[i])) {
				invalid = append(invalid, fmt.Sprintf("**%s**", names[i]))
			}

		}

	}
	//get numeric IDs
	if len(numericIDs) > 0 {
		var datares []nameres
		b.db.Select(&datares, "SELECT id,name FROM elements WHERE id=ANY($1) AND guild=$2 ", pq.Array(numericIDs), c.Guild())

		for i := 0; i < len(datares); i++ {
			idmap[datares[i].ID] = datares[i].Name
			namemap[datares[i].Name] = datares[i].ID
			convert[fmt.Sprintf("#%d", datares[i].ID)] = datares[i].Name
		}
		for i := 0; i < len(numericIDs); i++ {
			_, ok := idmap[numericIDs[i]]

			if !ok && !slices.Contains(invalid, fmt.Sprintf("**#%d**", numericIDs[i])) {
				invalid = append(invalid, fmt.Sprintf("**#%d**", numericIDs[i]))
			}
		}
	}

	if len(invalid) == 0 {
		for i := 0; i < len(vals); i++ {
			id, ok := namemap[convert[strings.ToLower(strings.TrimSpace(vals[i]))]]
			if ok {
				ids = append(ids, id)
			}
		}
		return ids, true
	}
	if len(ids) == 0 && len(invalid) == 0 {
		c.Respond(sevcord.NewMessage("Invalid format! " + types.RedCircle))
		return nil, false
	}
	if len(invalid) == 1 {
		c.Respond(sevcord.NewMessage("Element **" + invalid[0] + "** doesn't exist! " + types.RedCircle))
		return nil, false
	} else {
		var orderedinvalid []string
		for i := 0; i < len(vals); i++ {
			if slices.Contains(invalid, fmt.Sprintf("**%s**", strings.ToLower(vals[i]))) {
				orderedinvalid = append(orderedinvalid, fmt.Sprintf("**%s**", vals[i]))
			}
		}

		output := makeListResp("Elements", "and", " don't exist!", orderedinvalid)
		c.Respond(sevcord.NewMessage(output))
		return nil, false
	}

}

func IsNumericID(val string) (int64, bool) {
	id, err := strconv.ParseInt(strings.TrimPrefix(strings.TrimSpace(val), "#"), 10, 64)
	if err == nil && strings.HasPrefix(val, "#") {
		return id, true
	} else {
		return -1, false
	}
}
func makeListResp(start, join, end string, vals []string) string {

	if len(vals) > 1 {
		var strbuilder strings.Builder
		endstr := false
		strbuilder.WriteString(start + " ")
		for i := 0; i < len(vals); i++ {
			checkstr := join + " " + vals[i] + end + " " + types.RedCircle
			if strbuilder.Len() >= 1850 {
				strbuilder.WriteString(checkstr)
				endstr = true
			} else if i == len(vals)-1 {
				strbuilder.WriteString(join + " " + vals[i])
			} else {
				strbuilder.WriteString(vals[i])
				if len(vals) > 2 && i < len(vals)-2 {
					strbuilder.WriteString(", ")
				} else {
					strbuilder.WriteString(" ")
				}

			}
			//discord caps at 2000 so stop before that
			if endstr {
				return strbuilder.String()
			}

		}

		strbuilder.WriteString(end + " " + types.RedCircle)
		return strbuilder.String()
	}
	return ""
}
func getSort(input string) string {

	switch strings.TrimSpace(strings.ToLower(input)) {
	case "creator":
		return "creator"
	case "name":
		return "name"
	case "id":
		return "id"
	case "created on", "createdon":
		return "createdon"
	case "treesize", "tree size":
		return "treesize"
	case "length":
		return "length"
	case "found":
		return "found"
	default:
		return "id"
	}
}
func getLbSort(input string) string {
	switch input {
	case "made":
		return "made"
	case "imaged", "img":
		return "img"
	case "votes":
		return "voted"
	case "signed":
		return "signed"
	case "colored", "colour", "coloured", "color":
		return "color"
	case "catsigned":
		return "catsigned"
	case "catimg", "catimage":
		return "catimg"
	case "catcolor":
		return "catcolor"
	case "querysigned":
		return "querysigned"
	case "queryimage", "queryimg":
		return "queryimg"
	case "querycolor":
		return "querycolor"
	case "found", "":
		return "found"
	default:
		return "found"
	}
}
