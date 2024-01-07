package polls

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
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
		// Check if exists
		err := p.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM categories WHERE guild=$1 AND LOWER(name)=$2)`, po.Guild, strings.ToLower(po.Data["cat"].(string))).Scan(&exists)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("cat: category already exists")
		}
		// Check if valid name
		name, ok := base.CheckName(po.Data["cat"].(string))
		if !ok.Ok {
			return ok.Error()
		}
		po.Data["cat"] = name
		// Make
		_, err = p.db.Exec(`INSERT INTO categories (guild, name, elements, comment, image, color, commenter, imager, colorer) VALUES ($1, $2, $3, $6, $4, $5, $4, $4, $4)`, po.Guild, name, pq.Array(elems), "", 0, "None")
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

	// Check if empty
	_, err = p.db.Exec(`DELETE FROM categories WHERE guild=$1 AND name=$2 AND elements=$3`, po.Guild, po.Data["cat"].(string), pq.Array([]int32{}))
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

func (p *Polls) catImageSuccess(po *types.Poll, newsFunc func(string)) error {
	// Update image
	_, err := p.db.Exec(`UPDATE categories SET image=$1, imager=$2 WHERE name=$3 AND guild=$4`, po.Data["new"], po.Creator, po.Data["cat"], po.Guild)
	if err != nil {
		return err
	}

	// News
	word := "Changed"
	if po.Data["old"] == "" {
		word = "Added"
	}
	newsFunc(fmt.Sprintf("📸 %s Category Image - **%s** %s (%s)", word, po.Data["cat"].(string), p.pollContextMsg(po), po.Data["new"]))

	return nil
}

func (p *Polls) catMarkSuccess(po *types.Poll, newsFunc func(string)) error {
	// Update image
	_, err := p.db.Exec(`UPDATE categories SET comment=$1, commenter=$2 WHERE name=$3 AND guild=$4`, po.Data["new"], po.Creator, po.Data["cat"], po.Guild)
	if err != nil {
		return err
	}

	// News
	newsFunc(fmt.Sprintf("📝 Signed Category - **%s** %s", po.Data["cat"].(string), p.pollContextMsg(po)))

	return nil
}

func (p *Polls) catColorSuccess(po *types.Poll, newsFunc func(string)) error {
	// Update image
	_, err := p.db.Exec(`UPDATE categories SET color=$1, colorer=$2 WHERE name=$3 AND guild=$4`, po.Data["new"], po.Creator, po.Data["cat"], po.Guild)
	if err != nil {
		return err
	}

	// News
	newsFunc(fmt.Sprintf("🎨 Colored Category - **%s** %s", po.Data["cat"].(string), p.pollContextMsg(po)))

	return nil
}
