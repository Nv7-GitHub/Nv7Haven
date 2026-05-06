package polls

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var createlock = &sync.Mutex{}

func (e *Polls) elemCreate(p *types.Poll, news func(string)) (err error) {
	createlock.Lock()
	defer createlock.Unlock()

	els := util.Map(p.Data["els"].([]any), func(a any) int { return int(a.(float64)) })
	_, exists := p.Data["result"].(float64)

	// Check if combo has result
	var cmbexists bool
	err = e.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM combos WHERE els=$1 AND guild=$2)`, pq.Array(els), p.Guild).Scan(&cmbexists)
	if err != nil {
		return err
	}
	if cmbexists {
		return errors.New("already has result")
	}

	// Check if name already exists
	if !exists {
		var nameid int
		err = e.db.QueryRow(`SELECT id FROM elements WHERE LOWER(name)=$1 AND guild=$2`, strings.ToLower(p.Data["result"].(string)), p.Guild).Scan(&nameid)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if err == nil {
			exists = true
			p.Data["result"] = float64(nameid)
		}
	}

	// Make tx
	var tx *sqlx.Tx
	tx, err = e.db.Beginx()
	if err != nil {
		return
	}
	defer func() {
		if err != nil && tx != nil {
			err = tx.Rollback()
			return
		}
	}()

	// Create elem if not exists
	var id int
	var combid int // ID of combo if elem exists for news
	var name string
	if !exists {
		// Get id
		err = tx.QueryRow(`SELECT MAX(id) FROM elements WHERE guild=$1`, p.Guild).Scan(&id)
		if err != nil {
			return
		}
		id++

		// Get parents
		var parents []types.Element
		err = tx.Select(&parents, `SELECT id, color FROM elements WHERE id=ANY($1) AND guild=$2`, pq.Array(els), p.Guild)
		if err != nil {
			return
		}
		col := 0
		for _, parent := range parents {
			col += parent.Color
		}
		col /= len(parents)

		// Calc treesize
		var treeSize, tier int
		err = tx.QueryRow(`WITH RECURSIVE parents(els, id,depth) AS (
			VALUES($2::integer[], 1,0)
	 	UNION
			(SELECT b.parents els, b.id id,depth+1 FROM elements b INNER JOIN parents p ON b.id=ANY(p.els) where guild=$1)
	 	) SELECT COUNT(*), MAX(depth) FROM parents WHERE id>0`, p.Guild, pq.Array(els)).Scan(&treeSize, &tier)
		if err != nil {
			return
		}

		// Create element
		el := types.Element{
			ID:        id,
			Guild:     p.Guild,
			Name:      p.Data["result"].(string),
			Color:     col,
			Creator:   p.Creator,
			Comment:   types.DefaultMark,
			CreatedOn: time.Now(),
			Parents:   pq.Int32Array(util.Map(els, func(a int) int32 { return int32(a) })),
			TreeSize:  treeSize,
			MadeWith:  1,
			UsedIn:    0,
			FoundBy:   1,
			Tier:      tier,
		}
		name = el.Name

		// Insert element
		_, err = tx.NamedExec(`INSERT INTO elements (id, guild, name, image, color, comment, creator, createdon, commenter, colorer, imager, parents, treesize,madewith,usedin,foundby,tier) VALUES (:id, :guild, :name, :image, :color, :comment, :creator, :createdon, :commenter, :colorer, :imager, :parents, :treesize,:madewith, :usedin,:foundby,:tier)`, el)
		if err != nil {
			return
		}
	} else {
		id = int(p.Data["result"].(float64))

		// Combo ID
		err = tx.QueryRow(`SELECT COUNT(*) FROM combos WHERE guild=$1`, p.Guild).Scan(&combid)
		if err != nil {
			return
		}
		combid++

		// Get name
		var currtreesize int
		var madewith int
		err = tx.QueryRow(`SELECT name, treesize,madewith FROM elements WHERE id=$1 AND guild=$2`, id, p.Guild).Scan(&name, &currtreesize)
		if err != nil {
			return
		}
		madewith++
		_, err = tx.Exec(`UPDATE elements SET madewith=$1 WHERE id=$2 AND guild=$3`, madewith, id, p.Guild)
		if err != nil {
			return
		}
		// Check if need to update parents
		var treesize int
		var loop bool
		treesize, loop, err = e.base.TreeSize(tx, id, els, p.Guild)
		if err != nil {
			return
		}
		if !loop && treesize < currtreesize {
			// Update parents
			_, err = tx.Exec(`UPDATE elements SET parents=$1, treesize=$2 WHERE id=$3 AND guild=$4`, pq.Array(els), treesize, id, p.Guild)
			if err != nil {
				return
			}
		}

	}
	//update used in

	var parIDorder []int
	var parusedin []int

	err = tx.Select(&parIDorder, `SELECT id FROM elements WHERE id=ANY($1) AND guild=$2`, pq.Array(els), p.Guild)
	err = tx.Select(&parusedin, `SELECT usedin FROM elements WHERE id=ANY($1) AND guild=$2`, pq.Array(els), p.Guild)
	if err != nil {

	}
	for i, _ := range parIDorder {
		parusedin[i]++
		_, err = tx.Exec(`UPDATE elements SET usedin = $3 WHERE id=$1 AND guild=$2 `, parIDorder[i], p.Guild, parusedin[i])
		if err != nil {
			return
		}
	}

	// Create combo
	_, err = tx.Exec(`INSERT INTO combos (guild, els, result, createdon) VALUES ($1, $2, $3, $4)`, p.Guild, pq.Array(els), id, time.Now())
	if err != nil {
		return
	}

	// Commit
	err = tx.Commit()
	tx = nil
	if err != nil {
		return
	}

	// Add to creator's inv if not already in it
	var cont bool
	err = e.db.QueryRow(`SELECT $3=ANY(inv) FROM inventories WHERE guild=$1 AND "user"=$2`, p.Guild, p.Creator, id).Scan(&cont)
	if err != nil {
		return
	}
	if !cont {
		_, err = e.db.Exec(`UPDATE inventories SET inv=array_append(inv, $3) WHERE guild=$1 AND "user"=$2`, p.Guild, p.Creator, id)
		if err != nil {
			return
		}
	}

	// Message
	if exists {
		news(fmt.Sprintf("🆕 Combination - **%s** %s - Combination **#%d**", name, e.pollContextMsg(p), combid))
	} else {
		news(fmt.Sprintf("🆕 Element - **%s** %s - Element **#%d**", name, e.pollContextMsg(p), id))
	}
	return
}
