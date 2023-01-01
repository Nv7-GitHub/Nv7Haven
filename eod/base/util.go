package base

import (
	"fmt"
	"log"

	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
)

func (b *Base) Error(ctx sevcord.Ctx, err error) {
	if err != nil {
		ctx.Acknowledge()
		ctx.Respond(sevcord.NewMessage("").AddEmbed(
			sevcord.NewEmbed().
				Title("Error").
				Color(15548997). // Red
				Description("```" + err.Error() + "```"),
		))
	}
}

func (b *Base) IsPlayChannel(c sevcord.Ctx) bool {
	// Check if play channel
	var cnt bool
	err := b.db.QueryRow(`SELECT $1=ANY(play) FROM config WHERE guild=$2`, c.Channel(), c.Guild()).Scan(&cnt)
	if err != nil {
		fmt.Println("Play channel error", err)
		return false
	}
	return cnt
}

func (b *Base) PageLength(ctx sevcord.Ctx) int {
	if b.IsPlayChannel(ctx) {
		return 30
	}
	return 10
}

func (b *Base) NameMap(items []int, guild string) (map[int]string, error) {
	var names []struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}
	err := b.db.Select(&names, "SELECT id, name FROM elements WHERE id = ANY($1) AND guild=$2", pq.Array(items), guild)
	if err != nil {
		return nil, err
	}
	nameMap := make(map[int]string)
	for _, v := range names {
		nameMap[v.ID] = v.Name
	}
	return nameMap, nil
}

func (b *Base) GetNames(items []int, guild string) ([]string, error) {
	m, err := b.NameMap(items, guild)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(items))
	for i, v := range items {
		names[i] = m[v]
	}
	return names, nil
}

func (b *Base) GetName(guild string, elem int) (string, error) {
	var name string
	err := b.db.QueryRow("SELECT name FROM elements WHERE id=$1 AND guild=$2", elem, guild).Scan(&name)
	return name, err
}

func (b *Base) IncrementCommandStat(c sevcord.Ctx, name string) {
	_, err := b.db.Exec("INSERT INTO command_stats (guild, command, count) VALUES ($1, $2, 1) ON CONFLICT (guild, command) DO UPDATE SET count = command_stats.count + 1", c.Guild(), name)
	if err != nil {
		log.Println("command stats error", err)
	}
}
