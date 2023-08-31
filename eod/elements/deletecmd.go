package elements

import (
	"fmt"
	"log"

	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
)

func (e *Elements) deleteNewsMessage(c sevcord.Ctx, message string) {
	var news string
	err := e.db.QueryRow(`SELECT news FROM config WHERE guild=$1`, c.Guild()).Scan(&news)
	if err != nil {
		log.Println("news err", err)
		return
	}
	_, err = c.Dg().ChannelMessageSend(news, fmt.Sprintf("🗑️ "+message))
	if err != nil {
		log.Println("news err", err)
	}
}

func (e *Elements) DeleteComboCmd(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	_, err := e.db.Exec("DELETE FROM combos WHERE guild=$2 AND result=$1 AND els NOT IN (SELECT els FROM combos WHERE guild=$2 AND result=$1 ORDER BY createdon ASC LIMIT 1);", opts[0].(int64), c.Guild())
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Get remaining combo elements
	var els pq.Int64Array
	err = e.db.QueryRow(`SELECT els FROM combos WHERE guild=$2 AND result=$1`, opts[0].(int64), c.Guild()).Scan(&els)
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Get new tree size
	var treesize int
	err = e.db.QueryRow(`WITH RECURSIVE parents(els, id) AS (
		VALUES($2::integer[], 0)
	UNION
		(SELECT b.parents els, b.id id FROM elements b INNER JOIN parents p ON b.id=ANY(p.els) where guild=$1)
	) SELECT COUNT(*)  FROM parents WHERE id>0`, c.Guild(), pq.Array(els), opts[0].(int64)).Scan(&treesize, &loop)

	// Update parents
	_, err = e.db.Exec(`UPDATE elements SET parents=$1, treesize=$2 WHERE id=$3 AND guild=$4`, pq.Array(els), treesize, opts[0].(int64), c.Guild())
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Get element name
	nameE, err := e.base.GetName(c.Guild(), int(opts[0].(int64)))
	if err != nil {
		e.base.Error(c, err)
		return
	}
	// Send message in news
	c.Respond(sevcord.NewMessage("Successfully deleted all extra combos! ✅"))
	e.deleteNewsMessage(c, fmt.Sprintf("Deleted Combos - **%s**", nameE))
}
