package elements

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
func (e *Elements) InvHandler(c sevcord.Ctx, params string) {
	parts := strings.Split(params, "|")

	// Get count
	var cnt int
	err := e.db.QueryRow(`SELECT array_length(inv, 1) FROM inventories WHERE guild=$1 AND "user"=$2`, c.Guild(), parts[1]).Scan(&cnt)
	if err != nil {
		e.base.Error(c, err)
		return
	}
	length := e.base.PageLength(c)
	pagecnt := int(math.Ceil(float64(cnt) / float64(length)))

	// Apply page
	page, _ := strconv.Atoi(parts[3])
	switch parts[0] {
	case "prev":
		page--

	case "next":
		page++
	}
	if page < 0 {
		page = pagecnt - 1
	}
	if page >= pagecnt {
		page = 0
	}

	// Get values
	var inv []string
	err = e.db.Select(&inv, `SELECT name FROM elements WHERE id=ANY(SELECT UNNEST(inv) FROM inventories WHERE guild=$1 AND "user"=$2) AND guild=$1 ORDER BY $3 LIMIT $4 OFFSET $5`, c.Guild(), parts[1], parts[2], length, length*page)
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Get user
	m, err := c.Dg().GuildMember(c.Guild(), parts[1])
	if err != nil {
		u, err := c.Dg().User(parts[1])
		if err != nil {
			e.base.Error(c, err)
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
		AddComponentRow(types.PageSwitchBtns("inv", fmt.Sprintf("%s|%s|%d", parts[1], parts[2], page))...),
	)
}

func (e *Elements) Inv(c sevcord.Ctx, args []any) {
	c.Acknowledge()

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
	e.InvHandler(c, fmt.Sprintf("next|%s|%s|-1", user, sort))
}
