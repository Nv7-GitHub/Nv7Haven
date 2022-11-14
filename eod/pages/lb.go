package pages

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
)

var lbSorts = []sevcord.Choice{
	sevcord.NewChoice("Found", "found"),
	sevcord.NewChoice("Made", "made"),
	sevcord.NewChoice("Votes", "votes"),
	sevcord.NewChoice("Signed", "signed"),
	sevcord.NewChoice("Imaged", "img"),
	sevcord.NewChoice("Colored", "color"),
	sevcord.NewChoice("Categories Signed", "catsigned"),
	sevcord.NewChoice("Categories Imaged", "catimg"),
	sevcord.NewChoice("Categories Colored", "catcolor"),
}

// Params: [guild]
var lbSortCode = map[string]string{
	"found":     "array_length(inv, 1)",
	"made":      `(SELECT COUNT(*) FROM elements WHERE guild=$1 AND creator="user")`,
	"votes":     "votecnt",
	"signed":    `(SELECT COUNT(*) FROM elements WHERE guild=$1 AND commenter="user")`,
	"img":       `(SELECT COUNT(*) FROM elements WHERE guild=$1 AND imager="user")`,
	"color":     `(SELECT COUNT(*) FROM elements WHERE guild=$1 AND colorer="user")`,
	"catsigned": `(SELECT COUNT(*) FROM categories WHERE guild=$1 AND commenter="user")`,
	"catimg":    `(SELECT COUNT(*) FROM categories WHERE guild=$1 AND imager="user")`,
	"catcolor":  `(SELECT COUNT(*) FROM categories WHERE guild=$1 AND colorer="user")`,
}

// TODO: Support queries

// Format: prevnext|user|sort|page
func (p *Pages) LbHandler(c sevcord.Ctx, params string) {
	parts := strings.Split(params, "|")

	// Get values
	var vals []struct {
		User string `db:"user"`
		Cnt  int    `db:"cnt"`
	}
	err := p.db.Select(&vals, `SELECT "user", `+lbSortCode[parts[2]]+` cnt FROM inventories WHERE guild=$1 ORDER BY cnt DESC`, c.Guild())
	if err != nil {
		p.base.Error(c, err)
		return
	}

	// Get count
	cnt := len(vals)
	length := p.base.PageLength(c)
	pagecnt := int(math.Ceil(float64(cnt) / float64(length)))

	// Apply page
	page, _ := strconv.Atoi(parts[3])
	page = ApplyPage(parts[0], page, pagecnt)

	// Get res
	res := vals[page*length : util.Min((page+1)*length, cnt)]

	// Check if contains
	contains := false
	for _, v := range res {
		if v.User == parts[1] {
			contains = true
			break
		}
	}
	var pos int
	var usercnt int
	if !contains {
		for i, v := range vals {
			if v.User == parts[1] {
				pos = i
				usercnt = v.Cnt
				break
			}
		}
	}

	// Create description
	description := &strings.Builder{}
	for i, v := range res {
		youTxt := ""
		if v.User == parts[1] {
			youTxt = " *You*"
		}
		fmt.Fprintf(description, "%d. <@%s>%s - %s\n", i+1+length*page, v.User, youTxt, humanize.Comma(int64(v.Cnt)))
	}
	if !contains {
		fmt.Fprintf(description, "\n%d. <@%s> *You* - %s", pos+1, parts[1], humanize.Comma(int64(usercnt)))
	}

	// Respond
	// Get title name
	title := ""
	for _, opt := range lbSorts {
		if opt.Value == parts[2] {
			title = opt.Name
			break
		}
	}
	emb := sevcord.NewEmbed().
		Title("Top Most "+title).
		Description(description.String()).
		Color(1752220). // Aqua
		Footer(fmt.Sprintf("Page %d/%d", page+1, pagecnt), "")
	c.Respond(sevcord.NewMessage("").
		AddEmbed(emb).
		AddComponentRow(PageSwitchBtns("lb", fmt.Sprintf("%s|%s|%d", parts[1], parts[2], page))...))
}

func (p *Pages) Lb(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	p.base.IncrementCommandStat(c, "lb")

	// Params
	sort := "found"
	if opts[0] != nil {
		sort = opts[0].(string)
	}
	user := c.Author().User.ID
	if opts[1] != nil {
		user = opts[1].(string)
	}

	// Handler
	p.LbHandler(c, "next|"+user+"|"+sort+"|-1")
}
