package pages

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
)

// Format: prevnext|page
func (p *Pages) CommandLbHandler(c sevcord.Ctx, params string) {
	parts := strings.Split(params, "|")

	// Get count
	var cnt int
	err := p.db.QueryRow(`SELECT COUNT(*) FROM command_stats WHERE guild=$1`, c.Guild()).Scan(&cnt)
	if err != nil {
		p.base.Error(c, err)
		return
	}
	length := p.base.PageLength(c)
	pagecnt := int(math.Ceil(float64(cnt) / float64(length)))

	// Apply page
	page, _ := strconv.Atoi(parts[1])
	page = ApplyPage(parts[0], page, pagecnt)

	// Get values
	var vals []struct {
		Command string `db:"command"`
		Count   int    `db:"count"`
	}
	err = p.db.Select(&vals, `SELECT command, count FROM command_stats WHERE guild=$1 ORDER BY count DESC LIMIT $2 OFFSET $3`, c.Guild(), length, length*page)
	if err != nil {
		p.base.Error(c, err)
		return
	}

	// Description
	desc := &strings.Builder{}
	for i, val := range vals {
		fmt.Fprintf(desc, "%d\\. **%s** - %s\n", i+1+length*page, val.Command, humanize.Comma(int64(val.Count)))
	}

	// Create
	embed := sevcord.NewEmbed().
		Title("Command Usage").
		Description(desc.String()).
		Color(1146986). // Dark aqua
		Footer(fmt.Sprintf("Page %d/%d", page+1, pagecnt), "")
	c.Respond(sevcord.NewMessage("").
		AddEmbed(embed).
		AddComponentRow(PageSwitchBtns("cmdlb", fmt.Sprintf("%d", page))...))
}

func (p *Pages) CommandLb(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	p.CommandLbHandler(c, "next|-1")
}
