package base

import (
	"fmt"

	"github.com/Nv7-Github/sevcord/v2"
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
