package achievements

import (
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

func (a *Achievements) CheckRequirements(ctx sevcord.Ctx, achievementkind types.AchievementKind) {

	var achievements []types.Achievement
	var achievement types.Achievement
	var userachievements pq.Int32Array
	err := a.db.QueryRow(`SELECT achievements FROM achievers WHERE guild=$1 AND "user"=$2`, ctx.Guild(), ctx.Author().User.ID).Scan(&userachievements)
	if err != nil {

	}
	if len(userachievements) > 0 {

		err = a.db.Get(&achievements, `SELECT * FROM achievements WHERE guild=$1 AND kind=$2 AND NOT id = ANY($3)`, ctx.Guild(), string(achievementkind), userachievements)
	} else {

		err = a.db.Get(&achievements, `SELECT * FROM achievements WHERE guild=$1 AND kind=$2`, ctx.Guild(), string(achievementkind))
	}
	//only 1 achievement in achievements
	if err != nil {
		if len(userachievements) > 0 {

			err = a.db.Get(&achievement, `SELECT * FROM achievements WHERE guild=$1 AND kind=$2 AND NOT id = ANY($3)`, ctx.Guild(), string(achievementkind), userachievements)
		} else {

			err = a.db.Get(&achievement, `SELECT * FROM achievements WHERE guild=$1 AND kind=$2`, ctx.Guild(), string(achievementkind))
		}
	}
	if err != nil {
		return
	}
	if len(achievements) == 0 {
		achievements = append(achievements, achievement)
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
				}
				found = false

			}
			if found {
				a.EarnAchievement(ctx, achievements[i])
			}
		case types.AchievementKindCatNum:
			var common float64
			err = a.db.QueryRow(`SELECT COALESCE(array_length(elements & (SELECT inv FROM inventories WHERE guild=$1 AND "user"=$3), 1), 0) FROM categories WHERE guild=$1 AND name=$2`, ctx.Guild(), achievements[i].Data["cat"], ctx.Author().User.ID).Scan(&common)

			if err != nil {
				continue
			}
			if common >= achievements[i].Data["num"].(float64) {
				a.EarnAchievement(ctx, achievements[i])
			}

		}

	}

}
func (a *Achievements) EarnAchievement(ctx sevcord.Ctx, ac types.Achievement) {

	a.db.Exec(`UPDATE achievers SET achievements =array_append(achievements,$3) WHERE guild=$1 AND "user"=$2`, ctx.Guild(), ctx.Author().User.ID, ac.ID)
	ctx.Respond(sevcord.NewMessage(ctx.Author().Mention() + " New achievement earned: **" + ac.Name + "** " + ac.Icon))
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
