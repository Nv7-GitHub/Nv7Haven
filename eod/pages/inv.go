package pages

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/bwmarrin/discordgo"
)

// Params: prevnext|user|sort|page
func (p *Pages) InvHandler(c sevcord.Ctx, params string) {
	parts := strings.Split(params, "|")

	// Get count
	var cnt int
	err := p.db.QueryRow(`SELECT array_length(inv, 1) FROM inventories WHERE guild=$1 AND "user"=$2`, c.Guild(), parts[1]).Scan(&cnt)
	if err != nil {
		p.base.Error(c, err)
		return
	}
	length := p.base.PageLength(c)
	pagecnt := int(math.Ceil(float64(cnt) / float64(length)))

	// Apply page
	page, _ := strconv.Atoi(parts[3])
	page = ApplyPage(parts[0], page, pagecnt)

	// Get values
	var inv []string
	err = p.db.Select(&inv, `SELECT name FROM elements WHERE id=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$2) AND guild=$1 ORDER BY `+types.SortSql[parts[2]]+` LIMIT $3 OFFSET $4`, c.Guild(), parts[1], length, length*page)
	if err != nil {
		p.base.Error(c, err)
		return
	}

	// Get user
	m, err := c.Dg().GuildMember(c.Guild(), parts[1])
	if err != nil {
		u, err := c.Dg().User(parts[1])
		if err != nil {
			p.base.Error(c, err)
			return
		}
		m = &discordgo.Member{User: u, Nick: ""}
		return
	}
	name := m.User.Username
	if m.Nick != "" {
		name = m.Nick
	}

	// Create
	embed := sevcord.NewEmbed().
		Title(fmt.Sprintf("%s's Inventory", name)).
		Description(strings.Join([]string(inv), "\n")).
		Footer(fmt.Sprintf("Page %d/%d", page+1, pagecnt), "").
		Color(15105570) // Orange

	c.Respond(sevcord.NewMessage("").
		AddEmbed(embed).
		AddComponentRow(PageSwitchBtns("inv", fmt.Sprintf("%s|%s|%d", parts[1], parts[2], page))...),
	)
}

func (p *Pages) Inv(c sevcord.Ctx, args []any) {
	c.Acknowledge()
	p.base.IncrementCommandStat(c, "inv")

	// Get params
	user := c.Author().User.ID
	if args[0] != nil {
		user = args[0].(*discordgo.User).ID
	}
	sort := "id"
	if args[1] != nil {
		sort = args[1].(string)
	}

	// Create embed
	p.InvHandler(c, fmt.Sprintf("next|%s|%s|-1", user, sort))
}
