package eod

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
)

func (b *Bot) getElementId(c sevcord.Ctx, val string, showerr bool) (int64, bool) {
	var id int64
	var err error
	id, err = strconv.ParseInt(strings.TrimPrefix(strings.TrimSpace(val), "#"), 10, 64)
	if err == nil && strings.HasPrefix(val, "#") {
		err = b.db.QueryRow("SELECT id FROM elements WHERE id=$1 AND guild=$2", strings.ToLower(strings.TrimLeft(strings.TrimSpace(val), "#")), c.Guild()).Scan(&id)
	} else {
		err = b.db.QueryRow("SELECT id FROM elements WHERE LOWER(name)=$1 AND guild=$2", strings.ToLower(strings.TrimSpace(val)), c.Guild()).Scan(&id)
	}
	if err != nil && showerr {
		b.base.Error(c, err, "Element **"+val+"** doesn't exist!")
		return 0, false
	} else if err != nil && !showerr {
		return 0, false
	}

	return id, true
}
func makeListResp(start, join, end string, vals []string) string {
	if len(vals) == 2 {
		return fmt.Sprintf("%s %s %s %s%s %s", start, vals[0], join, vals[1], end, types.RedCircle)
	} else if len(vals) > 2 {
		return fmt.Sprintf("%s %s, %s %s%s %s", start, strings.Join(vals[:len(vals)-1], ", "), join, vals[len(vals)-1], end, types.RedCircle)
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
