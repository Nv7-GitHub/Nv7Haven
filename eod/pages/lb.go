package pages

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/lib/pq"
)

var lbSorts = []sevcord.Choice{
	sevcord.NewChoice("Found", "found"),
	sevcord.NewChoice("Invented", "made"),
	sevcord.NewChoice("Votes", "votes"),
	sevcord.NewChoice("Signed", "signed"),
	sevcord.NewChoice("Imaged", "img"),
	sevcord.NewChoice("Colored", "color"),
	sevcord.NewChoice("Categories Signed", "catsigned"),
	sevcord.NewChoice("Categories Imaged", "catimg"),
	sevcord.NewChoice("Categories Colored", "catcolor"),
	sevcord.NewChoice("Queries Signed", "querysigned"),
	sevcord.NewChoice("Queries Imaged", "queryimg"),
	sevcord.NewChoice("Queries Colored", "querycolor"),
}

// Params: [guild]
var lbSortCode = map[string]string{
	"found":       "array_length(inv, 1)",
	"made":        `(SELECT COUNT(*) FROM elements WHERE guild=$1 AND creator="user")`,
	"votes":       "votecnt",
	"signed":      `(SELECT COUNT(*) FROM elements WHERE guild=$1 AND commenter="user")`,
	"img":         `(SELECT COUNT(*) FROM elements WHERE guild=$1 AND imager="user")`,
	"color":       `(SELECT COUNT(*) FROM elements WHERE guild=$1 AND colorer="user")`,
	"catsigned":   `(SELECT COUNT(*) FROM categories WHERE guild=$1 AND commenter="user")`,
	"catimg":      `(SELECT COUNT(*) FROM categories WHERE guild=$1 AND imager="user")`,
	"catcolor":    `(SELECT COUNT(*) FROM categories WHERE guild=$1 AND colorer="user")`,
	"querysigned": `(SELECT COUNT(*) FROM queries WHERE guild=$1 AND commenter="user")`,
	"queryimg":    `(SELECT COUNT(*) FROM queries WHERE guild=$1 AND imager="user")`,
	"querycolor":  `(SELECT COUNT(*) FROM queries WHERE guild=$1 AND colorer="user")`,
}
var lbQuerySortCode = map[string]string{
	"found":  "COALESCE(array_length(inv & $2, 1), 0)",
	"made":   `(SELECT COUNT(*) FROM elements WHERE guild=$1 AND creator="user" AND id=ANY($2))`,
	"signed": `(SELECT COUNT(*) FROM elements WHERE guild=$1 AND commenter="user" AND id=ANY($2))`,
	"img":    `(SELECT COUNT(*) FROM elements WHERE guild=$1 AND imager="user" AND id=ANY($2))`,
	"color":  `(SELECT COUNT(*) FROM elements WHERE guild=$1 AND colorer="user" AND id=ANY($2))`,
}

// Format: prevnext|user|sort|page|query
func (p *Pages) LbHandler(c sevcord.Ctx, params string) {
	parts := strings.Split(params, "|")

	// Query
	var qu *types.Query
	var ok bool
	if parts[4] != "" {
		qu, ok = p.base.CalcQuery(c, parts[4])
		if !ok {
			return
		}
	}

	// Get values
	var err error
	var vals []struct {
		User string `db:"user"`
		Cnt  int    `db:"cnt"`
	}
	if qu == nil {
		err = p.db.Select(&vals, `SELECT "user", `+lbSortCode[parts[2]]+` cnt FROM inventories WHERE guild=$1 ORDER BY cnt DESC`, c.Guild())
	} else {
		sort, ok := lbQuerySortCode[parts[2]]
		if !ok {
			c.Respond(sevcord.NewMessage("This sort doesn't support queries! " + types.RedCircle))
			return
		}
		err = p.db.Select(&vals, `SELECT "user", `+sort+` cnt FROM inventories WHERE guild=$1 ORDER BY cnt DESC`, c.Guild(), pq.Array(qu.Elements))
	}
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
		fmt.Fprintf(description, "%d\\. <@%s>%s - %s\n", i+1+length*page, v.User, youTxt, humanize.Comma(int64(v.Cnt)))
	}
	if !contains {
		fmt.Fprintf(description, "\n%d\\. <@%s> *You* - %s", pos+1, parts[1], humanize.Comma(int64(usercnt)))
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
	if qu != nil {
		title += " (" + qu.Name + ")"
	}
	emb := sevcord.NewEmbed().
		Title("Top Most "+title).
		Description(description.String()).
		Color(1752220). // Aqua
		Footer(fmt.Sprintf("Page %d/%d", page+1, pagecnt), "")
	c.Respond(sevcord.NewMessage("").
		AddEmbed(emb).
		AddComponentRow(PageSwitchBtns("lb", fmt.Sprintf("%s|%s|%d|%s", parts[1], parts[2], page, parts[4]))...))
}

func (p *Pages) Lb(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	// Params
	sort := "found"
	if opts[0] != nil {
		sort = opts[0].(string)
	}
	user := c.Author().User.ID
	if opts[1] != nil {
		user = opts[1].(*discordgo.User).ID
	}
	query := ""
	if opts[2] != nil {
		query = opts[2].(string)
	}

	// Handler
	p.LbHandler(c, "next|"+user+"|"+sort+"|-1|"+query)
}
