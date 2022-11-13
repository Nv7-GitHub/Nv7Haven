package base

import (
	"fmt"

	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
)

func (b *Base) Error(ctx sevcord.Ctx, err error) {
	if err != nil {
		ctx.Acknowledge()
		ctx.Respond(sevcord.NewMessage("").AddEmbed(
			sevcord.NewEmbed().
				Title("Error").
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

func (b *Base) NameMap(items []int) (map[int]string, error) {
	var names []struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}
	err := b.db.Select(&names, "SELECT id, name FROM elements WHERE id = ANY($1)", pq.Array(items))
	if err != nil {
		return nil, err
	}
	nameMap := make(map[int]string)
	for _, v := range names {
		nameMap[v.ID] = v.Name
	}
	return nameMap, nil
}

func (b *Base) GetNames(items []int) ([]string, error) {
	m, err := b.NameMap(items)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(items))
	for i, v := range items {
		names[i] = m[v]
	}
	return names, nil
}
