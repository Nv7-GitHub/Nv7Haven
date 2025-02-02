package achievements

import (
	"database/sql"
	"fmt"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
)

func (a *Achievements) createCmd(c sevcord.Ctx, name string, kind types.AchievementKind, data map[string]any) {

	var existsName string
	err := a.db.QueryRow("SELECT name FROM achievements WHERE data@>$1 AND data<@$1 AND kind=$3 AND guild=$2", types.PgData(data), c.Guild(), string(kind)).Scan(&existsName)
	if err != nil && err != sql.ErrNoRows {
		a.base.Error(c, err)
		return
	}
	if err == nil {
		c.Respond(sevcord.NewMessage(fmt.Sprintf("Achievement **%s** already exists with this data! "+types.RedCircle, existsName)))
		return
	}
	var ID int
	err = a.db.QueryRow("SELECT nextval('achievements_id_seq')").Scan(&ID)
	if err != nil {
		a.base.Error(c, err)
		return
	}
	_, err = a.db.Exec(`INSERT INTO achievements (guild,id,name,kind,icon,data) VALUES ($1,$2,$3,$4,$5,$6)`, c.Guild(), ID, name, string(kind), "🏅", types.PgData(data))
	if err != nil {
		a.base.Error(c, err)
		return
	}
	c.Respond(sevcord.NewMessage("Created achievement 🏅"))
}
func (a *Achievements) CreateElementCmd(c sevcord.Ctx, opts []any) {

	c.Acknowledge()
	// Check if element exists
	var exists bool
	err := a.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM elements WHERE id=$1 AND guild=$2)", opts[1].(int64), c.Guild())
	if err != nil {
		a.base.Error(c, err)
		return
	}
	if !exists {
		c.Respond(sevcord.NewMessage("Element does not exist! " + types.RedCircle))
		return
	}
	a.createCmd(c, opts[0].(string), types.AchievementKindElement, map[string]any{"elem": float64(opts[1].(int64))})
}
func (a *Achievements) CreateCatNumCmd(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	//Check if cat exists
	var exists bool

	err := a.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM categories WHERE name=$1 AND guild=$2 )", opts[1].(string), c.Guild())
	if err != nil {
		a.base.Error(c, err)
		return
	}
	if !exists {
		c.Respond(sevcord.NewMessage("Category does not exist! " + types.RedCircle))
		return
	}
	var elements pq.Int32Array
	err = a.db.QueryRow("SELECT elements FROM categories WHERE name=$1 AND guild=$2", opts[1].(string), c.Guild()).Scan(&elements)
	if err != nil {
		return
	}
	if opts[2].(int) <= 0 /*|| opts[2].(int64) > len(elements) */ {
		c.Respond(sevcord.NewMessage("Invalid number! " + types.RedCircle))
		return
	}

	a.createCmd(c, opts[0].(string), types.AchievementKindCatNum, types.PgData{
		"cat": opts[1].(string),
		"num": opts[2].(int64),
	})

}
func (a *Achievements) CreateCatPercentCmd(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	var exists bool

	err := a.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM categories WHERE name=$1)", opts[1].(string))
	if err != nil {
		a.base.Error(c, err)
		return
	}
	if !exists {
		c.Respond(sevcord.NewMessage("Category does not exist! " + types.RedCircle))
		return
	}

	a.createCmd(c, opts[0].(string), types.AchievementKindCatPercent, types.PgData{
		"cat":     opts[1].(string),
		"percent": opts[2].(float32),
	})

}
