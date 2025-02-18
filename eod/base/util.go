package base

import (
	"database/sql"
	"log"
	"slices"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
)

func (b *Base) Error(ctx sevcord.Ctx, err error, config ...string) {
	if err != nil {
		ctx.Acknowledge()

		if err == sql.ErrNoRows && len(config) >= 1 {
			ctx.Respond(sevcord.NewMessage(config[0] + " " + types.RedCircle))
			return
		}

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
	config, resp := b.GetConfigCache(c)

	if !resp.Ok {
		err := b.db.Get(&config, `SELECT * FROM config WHERE guild=$1`, c.Guild())
		if err != nil {
			log.Println("Config error", err)
			return false
		}

	}
	var playchannels []string
	playchannels = append(playchannels, config.PlayChannels...)
	if len(playchannels) == 0 {
		log.Println("Play channel error")
		return false
	}
	if slices.Contains(playchannels, c.Channel()) {
		return true
	}
	return false

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
