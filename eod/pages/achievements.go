package pages

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/lib/pq"
)

var achievmentListSorts = []sevcord.Choice{
	sevcord.NewChoice("Name", "name"),
}

func (p *Pages) AchievementList(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	sort := "name"
	if opts[0] != nil {
		sort = opts[0].(string)
	}
	p.AchievementListHandler(c, "next|"+c.Author().User.ID+"|"+sort+"|-1")

}

// Format: prevnext|user|sort|page
func (p *Pages) AchievementListHandler(c sevcord.Ctx, params string) {
	parts := strings.Split(params, "|")
	var cnt int
	err := p.db.QueryRow(`SELECT COUNT(*) FROM achievements WHERE guild=$1`, c.Guild()).Scan(&cnt)
	if err != nil {
		p.base.Error(c, err)
		return
	}
	length := p.base.PageLength(c)
	pagecnt := int(math.Ceil(float64(cnt) / float64(length)))

	// Apply pages
	page, _ := strconv.Atoi(parts[3])
	page = ApplyPage(parts[0], page, pagecnt)

	var items []struct {
		Name string `db:"name"`
		Cont bool   `db:"cont"`
	}

	err = p.db.Select(&items, `SELECT name,id=ANY(SELECT UNNEST(achievements) FROM achievers WHERE guild=$1 AND "user"=$2) cont FROM achievements ORDER BY`+` name `+`LIMIT $3 OFFSET $4`, c.Guild(), parts[1], length, length*page)
	if err != nil {
		return
	}
	desc := &strings.Builder{}
	for _, v := range items {
		if v.Cont {
			fmt.Fprintf(desc, "%s %s\n", v.Name, types.Check)
		} else {
			fmt.Fprintf(desc, "%s %s\n", v.Name, types.NoCheck)
		}
	}
	emb := sevcord.NewEmbed().
		Title(fmt.Sprintf("All Achievements (%d)", cnt)).
		Description(desc.String()).
		Footer(fmt.Sprintf("Page %d/%d", page+1, pagecnt), "").
		Color(10181046) // Purple
	c.Respond(sevcord.NewMessage("").AddEmbed(emb).AddComponentRow(PageSwitchBtns("achievementlist", fmt.Sprintf("%s|%s|%d", parts[1], parts[2], page))...))

}
func (p *Pages) UserAchievments(c sevcord.Ctx, args []any) {

	c.Acknowledge()

	// Get params
	user := c.Author().User.ID
	if args[0] != nil {
		user = args[0].(*discordgo.User).ID
	}
	sort := "name"
	if args[1] != nil {
		sort = args[1].(string)
	}

	// Create embed
	p.UserAchievementsHandler(c, fmt.Sprintf("next|%s|%s|-1", user, sort))
}

// Params: prevnext|user|sort|page
func (p *Pages) UserAchievementsHandler(c sevcord.Ctx, params string) {
	parts := strings.Split(params, "|")

	var cntarr pq.Int32Array
	err := p.db.Get(&cntarr, `SELECT achievements FROM achievers WHERE guild=$1 AND "user"=$2`, c.Guild(), parts[1])
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

	length := p.base.PageLength(c)
	pagecnt := int(math.Ceil(float64(len(cntarr)) / float64(length)))

	// Apply page
	page, _ := strconv.Atoi(parts[3])
	page = ApplyPage(parts[0], page, pagecnt)

	// Get values
	var inv []struct {
		Name string `db:"name"`
		Cont bool   `db:"cont"`
	}
	err = p.db.Select(&inv, `SELECT name, id=ANY(SELECT UNNEST(achievements) FROM achievers WHERE guild=$1 AND "user"=$5) cont FROM achievements WHERE id=ANY(SELECT UNNEST(achievements) FROM achievers WHERE guild=$1 AND "user"=$2) AND guild=$1 ORDER BY `+`name`+` LIMIT $3 OFFSET $4`, c.Guild(), parts[1], length, length*page, c.Author().User.ID)
	if err != nil {
		p.base.Error(c, err)
		return
	}

	desc := &strings.Builder{}
	for _, v := range inv {
		if c.Author().User.ID != parts[1] {
			if v.Cont {
				fmt.Fprintf(desc, "%s %s\n", v.Name, types.Check)
			} else {
				fmt.Fprintf(desc, "%s %s\n", v.Name, types.NoCheck)
			}
		} else {
			fmt.Fprintf(desc, "%s\n", v.Name)
		}
	}

	embed := sevcord.NewEmbed().
		Title(fmt.Sprintf("%s's Achievements (%s)", name, humanize.Comma(int64(len(cntarr))))).
		Description(desc.String()).
		Footer(fmt.Sprintf("Page %d/%d", page+1, pagecnt), "").
		Color(15105570) // Orange

	c.Respond(sevcord.NewMessage("").
		AddEmbed(embed).
		AddComponentRow(PageSwitchBtns("userachievements", fmt.Sprintf("%s|%s|%d", parts[1], parts[2], page))...),
	)
}
func (p *Pages) AchievementFound(c sevcord.Ctx, args []any) {
	c.Acknowledge()

	id, _ := strconv.Atoi(args[0].(string))
	// Create embed
	p.AchievementFoundHandler(c, fmt.Sprintf("next|%d|-1", id))
}

// Params: prevnext|achievement|page
func (p *Pages) AchievementFoundHandler(c sevcord.Ctx, params string) {

	var cnt int

	parts := strings.Split(params, "|")
	length := p.base.PageLength(c)
	achievement, _ := strconv.Atoi(parts[1])
	err := p.db.Get(&cnt, `SELECT COUNT(*) FROM achievers WHERE $1= ANY(achievements) AND guild=$2`, achievement, c.Guild())
	if err != nil {
		p.base.Error(c, err)
		return
	}
	pagecnt := int(math.Ceil(float64(cnt) / float64(length)))

	// Apply page
	page, _ := strconv.Atoi(parts[2])
	page = ApplyPage(parts[0], page, pagecnt)

	// Get values
	var found []string
	err = p.db.Select(&found, `SELECT "user" FROM achievers WHERE $2=ANY(achievements) AND guild=$1 ORDER BY cardinality(achievements) DESC LIMIT $3 OFFSET $4`, c.Guild(), achievement, length, length*page)
	if err != nil {
		p.base.Error(c, err)
		return
	}

	// Make text
	desc := &strings.Builder{}
	for _, v := range found {
		fmt.Fprintf(desc, "<@%s>\n", v)
	}
	//get achievement name
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		p.base.Error(c, err)
		return
	}
	var name string
	err = p.db.Get(&name, `SELECT name from achievements WHERE id=$1 AND guild=$2`, id, c.Guild())
	if err != nil {
		p.base.Error(c, err)
		return
	}
	embed := sevcord.NewEmbed().
		Title(fmt.Sprintf("%s's Found (%s)", name, humanize.Comma(int64(cnt)))).
		Description(desc.String()).
		Footer(fmt.Sprintf("Page %d/%d", page+1, pagecnt), "").
		Color(15277667) // Pink

	c.Respond(sevcord.NewMessage("").
		AddEmbed(embed).
		AddComponentRow(PageSwitchBtns("achievementfound", fmt.Sprintf("%d|%d", parts[1], page))...),
	)
}
