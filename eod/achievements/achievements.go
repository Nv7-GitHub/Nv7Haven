package achievements

import (
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Achievements struct {
	db   *sqlx.DB
	base *base.Base
	s    *sevcord.Sevcord
}

func (a *Achievements) CheckRequirements(ctx sevcord.Ctx) {

	var achievements []types.Achievement
	var userachievements pq.Int32Array
	err := a.db.QueryRow(`SELECT achievements FROM achievers WHERE guild=$1 AND "user"=$2`, ctx.Guild(), ctx.Author().User.ID).Scan(&userachievements)
	if err != nil {

	}
	if len(userachievements) > 0 {

		err = a.db.Select(&achievements, `SELECT * FROM achievements WHERE guild=$1 AND NOT id = ANY($2)`, ctx.Guild(), userachievements)
	} else {

		err = a.db.Select(&achievements, `SELECT * FROM achievements WHERE guild=$1 `, ctx.Guild())
	}
	if err != nil {
		return
	}

	var inv pq.Int32Array
	err = a.db.QueryRow(`SELECT inv FROM inventories WHERE guild=$1 AND "user"=$2`, ctx.Guild(), ctx.Author().User.ID).Scan(&inv)
	if err != nil {
		return
	}
	for i := 0; i < len(achievements); i++ {
		switch achievements[i].Kind {

		case types.AchievementKindElement:
			found := false
			for j := 0; j < len(inv); j++ {
				if inv[j] == int32(achievements[i].Data["elem"].(float64)) {
					found = true
					break
				} else {
					found = false
				}

			}
			if found {
				a.EarnAchievement(ctx, achievements[i])
			}
		case types.AchievementKindCatNum:
			var common float64
			err = a.db.QueryRow(`SELECT COALESCE(array_length(elements & (SELECT inv FROM inventories WHERE guild=$1 AND "user"=$3), 1), 0) FROM categories WHERE guild=$1 AND LOWER(name)=$2`, ctx.Guild(), strings.ToLower(achievements[i].Data["cat"].(string)), ctx.Author().User.ID).Scan(&common)

			if err != nil {
				continue
			}
			if common >= achievements[i].Data["num"].(float64) {
				a.EarnAchievement(ctx, achievements[i])
			}
		case types.AchievementKindCatPercent:
			var cnt int
			var common int
			err := a.db.QueryRow(`SELECT array_length(elements, 1), COALESCE(array_length(elements & (SELECT inv FROM inventories WHERE guild=$1 AND "user"=$3), 1), 0) FROM categories WHERE guild=$1 AND LOWER(name)=$2`, ctx.Guild(), strings.ToLower(achievements[i].Data["cat"].(string)), ctx.Author().User.ID).Scan(&cnt, &common)
			if err != nil {
				a.base.Error(ctx, err)
				return
			}
			percent := float64(common) / float64(cnt) * 100
			if percent >= achievements[i].Data["percent"].(float64) {
				a.EarnAchievement(ctx, achievements[i])
			}
		case types.AchievementKindInvCnt:
			if float64(len(inv)) >= achievements[i].Data["num"].(float64) {
				a.EarnAchievement(ctx, achievements[i])
			}
		}
	}
}
func (a *Achievements) EarnAchievement(ctx sevcord.Ctx, ac types.Achievement) {

	a.db.Exec(`UPDATE achievers SET achievements =array_append(achievements,$3) WHERE guild=$1 AND "user"=$2`, ctx.Guild(), ctx.Author().User.ID, ac.ID)
	ctx.Respond(sevcord.NewMessage(ctx.Author().Mention() + " New achievement earned: **" + ac.Name + "** " + ac.Icon))
	var news string
	err := a.db.QueryRow("SELECT news FROM config WHERE guild=$1", ctx.Guild()).Scan(&news)
	if err != nil {
		return
	}
	a.s.Dg().ChannelMessageSend(news, "🏆 Earned Achievement - **"+ac.Name+" "+ac.Icon+"** (By "+ctx.Author().Mention()+")")
}

func (a *Achievements) Init() {

}
func NewAchievements(db *sqlx.DB, base *base.Base, s *sevcord.Sevcord) *Achievements {
	a := &Achievements{
		base: base,
		db:   db,
		s:    s,
	}
	a.Init()
	return a
}
