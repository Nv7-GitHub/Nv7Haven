package polls

import (
	"fmt"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/lib/pq"
)

func (p *Polls) categorizeSuccess(po *types.Poll, news func(string)) error {
	elems := util.Map(po.Data["elems"].([]any), func(v any) int32 { return int32(v.(float64)) })

	// Check if exists
	var exists bool
	err := p.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM categories WHERE guild=$1 AND name=$2)`, po.Guild, po.Data["cat"].(string)).Scan(&exists)
	if err != nil {
		return err
	}
	if exists { // Append
		_, err := p.db.Exec(`UPDATE categories SET elements=uniq(sort(elements || $1)) WHERE guild=$2 AND name=$3`, pq.Array(elems), po.Guild, po.Data["cat"].(string))
		if err != nil {
			return err
		}
	} else { // Create
		_, err := p.db.Exec(`INSERT INTO categories (guild, name, elements, comment, image, color, commenter, imager, colorer) VALUES ($1, $2, $3, $4, $4, $5, $4, $4, $4)`, po.Guild, po.Data["cat"].(string), elems, "", 0)
		if err != nil {
			return err
		}
	}

	// News
	var name string
	if len(elems) == 1 {
		name, err = p.base.GetName(po.Guild, int(elems[0]))
		if err != nil {
			return err
		}
	} else {
		name = fmt.Sprintf("%d elements", len(elems))
	}
	news(fmt.Sprintf("🗃️ Added **%s** to **%s** %s", name, po.Data["cat"].(string), p.pollContextMsg(po)))

	return nil
}

func (p *Polls) unCategorizeSuccess(po *types.Poll, news func(string)) error {
	elems := util.Map(po.Data["elems"].([]any), func(v any) int32 { return int32(v.(float64)) })

	// Remove
	_, err := p.db.Exec(`WITH elems (el) AS (
	SELECT UNNEST(elements) FROM categories WHERE guild=$1 AND name=$2
) UPDATE categories SET elements=ARRAY(
  SELECT el FROM elems WHERE NOT(el=ANY($3))
) WHERE guild=$1 AND name=$2`, po.Guild, po.Data["cat"].(string), pq.Array(elems))
	if err != nil {
		return err
	}

	// News
	var name string
	if len(elems) == 1 {
		name, err = p.base.GetName(po.Guild, int(elems[0]))
		if err != nil {
			return err
		}
	} else {
		name = fmt.Sprintf("%d elements", len(elems))
	}

	news(fmt.Sprintf("🗃️ Removed **%s** from **%s** %s", name, po.Data["cat"].(string), p.pollContextMsg(po)))
	return nil
}
