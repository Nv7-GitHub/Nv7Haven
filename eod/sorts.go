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

var IDPrefixes = []string{
	"#",
	"{id}",
}

func (b *Bot) MsgSugElement(c sevcord.Ctx, val string) {
	ok, name := b.checkElementExists(c, val)
	if ok {
		val = name
	}
	b.elements.Suggest(c, []any{any(val), nil})
}
func (b *Bot) combineElements(c sevcord.Ctx, elements []string) {

	ids, ok := b.getElementIds(c, elements)
	if ok {
		b.elements.Combine(c, ids)
	}

}
func (b *Bot) ApplyMultiplier(c sevcord.Ctx, val string) (ok bool, multelements []string) {
	if strings.HasPrefix(val, "*") {
		parts := strings.SplitN(val[1:], " ", 2)
		cnt, err := strconv.Atoi(parts[0])
		if err != nil {
			c.Respond(sevcord.NewMessage("Invalid number of repeats! " + types.RedCircle))
			return false, []string{}
		}
		if cnt > types.MaxComboLength {
			c.Respond(sevcord.NewMessage(fmt.Sprintf("You can only combine up to %d elements! "+types.RedCircle, types.MaxComboLength)))
			return false, []string{}
		}
		if cnt < 2 {
			c.Respond(sevcord.NewMessage("You need to combine at least 2 elements! " + types.RedCircle))
			return false, []string{}
		}
		if len(parts) == 2 {
			inps := make([]string, 0, cnt)
			for i := 0; i < cnt; i++ {
				inps = append(inps, strings.TrimSpace(parts[1]))
			}
			return true, inps
		} else {
			comb, ok := b.base.GetCombCache(c)
			if !ok.Ok {
				c.Respond(ok.Response())
				return false, []string{}
			}
			if comb.Result == -1 {
				c.Respond(sevcord.NewMessage("You haven't combined anything! " + types.RedCircle))
				return false, []string{}
			}
			name, err := b.base.GetName(c.Guild(), comb.Result)
			if err != nil {
				b.base.Error(c, err)
				return false, []string{}
			}
			new := make([]string, 0, cnt)
			for i := 0; i < cnt; i++ {
				new = append(new, name)
			}
			return true, new
		}
	}
	return false, []string{}
}
func (b *Bot) checkElementExists(c sevcord.Ctx, val string) (bool, string) {

	val = convertVariableID(c, b, val)
	var err error
	_, ok := IsNumericID(val)
	val = convertName(val)
	var name string
	if ok {
		err = b.db.QueryRow("SELECT name FROM elements WHERE id=$1 AND guild=$2", strings.ToLower(strings.TrimLeft(strings.TrimSpace(val), "#")), c.Guild()).Scan(&name)
	} else {
		err = b.db.QueryRow("SELECT name FROM elements WHERE LOWER(name)=$1 AND guild=$2", strings.ToLower(strings.TrimLeft(strings.TrimSpace(val), "#")), c.Guild()).Scan(&name)
	}
	if err != nil {
		return false, ""
	} else {
		return true, name
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

	prefix := ""
	teststr := ""
	for _, pre := range IDPrefixes {
		teststr = strings.ToLower(strings.TrimSpace(strings.TrimPrefix(val, pre)))
		if teststr != strings.ToLower(strings.TrimSpace(val)) {
			prefix = pre
			break
		}
	}

	switch teststr {
	case "last":
		cache, _ := b.base.GetCombCache(c)
		if cache.Result != -1 {
			return prefix + fmt.Sprintf("%d", cache.Result)
		}
	case "rand":
		var id int
		err := b.db.QueryRow(`SELECT id FROM elements WHERE guild=$1 ORDER BY RANDOM()`, c.Guild()).Scan(&id)
		if err == nil {
			return prefix + fmt.Sprintf("%d", id)
		}
	case "randinv", "randininv":
		var id int
		err := b.db.QueryRow(`SELECT id FROM elements WHERE guild=$1 AND id=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$2) ORDER BY RANDOM()`, c.Guild(), c.Author().User.ID).Scan(&id)
		if err == nil {
			return prefix + fmt.Sprintf("%d", id)
		}
	}
	return val

}
func convertName(val string) string {
	parts := strings.SplitN(val, "}", 2)
	if strings.HasPrefix(val, "{") && len(parts) > 1 {
		prefix := strings.TrimPrefix(strings.TrimSpace(parts[0]), "{")
		switch strings.ToLower(prefix) {
		case "raw", "text", "name":
			return strings.TrimLeft(parts[1], " ")
		}
	}

	return val
}
func IsNumericID(val string) (int64, bool) {

	for _, prefix := range IDPrefixes {
		id, err := strconv.ParseInt(strings.TrimPrefix(strings.TrimSpace(val), prefix), 10, 64)
		if err == nil {
			return id, true
		}
	}
	return -1, false
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
			//convert all ids to "#" format
			vals[i] = "#" + fmt.Sprintf("%d", id)
			numericIDs = append(numericIDs, id)
		} else {

			vals[i] = convertName(vals[i])
			names = append(names, strings.TrimSpace(strings.ToLower(vals[i])))
		}
		convert[strings.ToLower(vals[i])] = vals[i]
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
		c.Respond(sevcord.NewMessage("Element **" + convert[strings.TrimPrefix(strings.TrimSuffix(invalid[0], "**"), "**")] + "** doesn't exist! " + types.RedCircle))
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
