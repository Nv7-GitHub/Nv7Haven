package categories

import (
	"log"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
)

func (c *Categories) Autocomplete(ctx sevcord.Ctx, val any) []sevcord.Choice {
	var res []types.Element
	err := c.db.Select(&res, "SELECT name FROM categories WHERE guild=$1 AND name ILIKE $2 || '%' ORDER BY similarity(name, $2) DESC, id LIMIT 25", ctx.Guild(), val.(string))
	if err != nil {
		log.Println("cat autocomplete error", err)
		return nil
	}
	choices := make([]sevcord.Choice, len(res))
	for i, v := range res {
		choices[i] = sevcord.NewChoice(v.Name, v.Name)
	}
	return choices
}
