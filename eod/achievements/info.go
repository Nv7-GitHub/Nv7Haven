package achievements

import (
	"fmt"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
)

func (a *Achievements) Info(ctx sevcord.Ctx, opts []any) {

	ctx.Acknowledge()
	var ac types.Achievement

	err := a.db.Get(&ac, "SELECT * FROM achievements WHERE name=$1 AND guild=$2", opts[0].(string), ctx.Guild())
	if err != nil {
		ctx.Respond(sevcord.NewMessage("Achivement **" + opts[0].(string) + "** doesn't exist" + types.RedCircle))
		return
	}
	var description string
	var have bool
	err = a.db.QueryRow(`SELECT $1=ANY(achievements) FROM achievers WHERE guild=$2 AND "user"=$3`, ac.ID, ctx.Guild(), ctx.Author().User.ID).Scan(&have)
	if err != nil {
		return
	}
	if have {
		description = "📫 **You have this.**\n\n" + description
	} else {
		description = "📪 **You don't have this.**\n\n" + description
	}
	emb := sevcord.NewEmbed().
		Title(opts[0].(string) + " Info").
		Description(description)

	switch ac.Kind {
	case types.AchievementKindElement:
		emb = emb.AddField("Kind", "Element", true)
		name, err := a.base.GetName(ctx.Guild(), int(ac.Data["elem"].(float64)))
		if err != nil {
			a.base.Error(ctx, err)
			return
		}
		emb = emb.AddField("Element", name, true)
	case types.AchievementKindCatNum:
		emb = emb.AddField("Kind", "Category Number", true)
		emb = emb.AddField("Category", ac.Data["cat"].(string), true)
		emb = emb.AddField("Number", fmt.Sprintf("%f", ac.Data["num"]), true)
	}
	ctx.Respond(sevcord.NewMessage("").AddEmbed(emb))
}
