package polls

import (
	"fmt"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (p *Polls) queryCreateSuccess(po *types.Poll, news func(string)) error {
	// Update DB
	if po.Data["edit"].(bool) {
		_, err := p.db.Exec(`UPDATE queries SET data=$1, creator=$2, createdon=$3, kind=$6 WHERE name=$4 AND guild=$5`, types.PgData(po.Data["data"].(map[string]any)), po.Creator, time.Now(), po.Data["query"], po.Guild, po.Data["kind"])
		if err != nil {
			return err
		}
	} else {
		_, err := p.db.Exec(`INSERT INTO queries (guild, name, creator, createdon, kind, data, image, comment, imager, colorer, commenter, color) VALUES ($1, $2, $3, $4, $5, $6, $7, $9, $7, $7, $7, $8)`, po.Guild, po.Data["query"], po.Creator, time.Now(), po.Data["kind"], types.PgData(po.Data["data"].(map[string]any)), "", 0, types.DefaultMark)
		if err != nil {
			return err
		}
	}

	// News
	word := "Created"
	if po.Data["edit"].(bool) {
		word = "Edited"
	}
	news(fmt.Sprintf("🧮 %s Query - **%s** %s", word, po.Data["query"], p.pollContextMsg(po)))
	return nil
}

func (p *Polls) queryDeleteSuccess(po *types.Poll, news func(string)) error {
	// Delete
	_, err := p.db.Exec(`DELETE FROM queries WHERE name=$1 AND guild=$2`, po.Data["query"], po.Guild)
	if err != nil {
		return err
	}

	// News
	news(fmt.Sprintf("🗑️ Deleted Query - **%s** %s", po.Data["query"], p.pollContextMsg(po)))
	return nil
}

func (p *Polls) queryImageSuccess(po *types.Poll, newsFunc func(string)) error {
	// Update image
	_, err := p.db.Exec(`UPDATE queries SET image=$1, imager=$2 WHERE name=$3 AND guild=$4`, po.Data["new"], po.Creator, po.Data["query"], po.Guild)
	if err != nil {
		return err
	}

	// News
	word := "Changed"
	if po.Data["old"] == "" {
		word = "Added"
	}
	newsFunc(fmt.Sprintf("📸 %s Query Image - **%s** %s (%s)", word, po.Data["query"].(string), p.pollContextMsg(po), po.Data["new"]))

	return nil
}

func (p *Polls) queryMarkSuccess(po *types.Poll, newsFunc func(string)) error {
	// Update image
	_, err := p.db.Exec(`UPDATE queries SET comment=$1, commenter=$2 WHERE name=$3 AND guild=$4`, po.Data["new"], po.Creator, po.Data["query"], po.Guild)
	if err != nil {
		return err
	}

	// News
	newsFunc(fmt.Sprintf("📝 Signed Query - **%s** %s", po.Data["query"].(string), p.pollContextMsg(po)))

	return nil
}

func (p *Polls) queryColorSuccess(po *types.Poll, newsFunc func(string)) error {
	// Update image
	_, err := p.db.Exec(`UPDATE queries SET color=$1, colorer=$2 WHERE name=$3 AND guild=$4`, po.Data["new"], po.Creator, po.Data["query"], po.Guild)
	if err != nil {
		return err
	}

	// News
	newsFunc(fmt.Sprintf("🎨 Colored Query - **%s** %s", po.Data["query"].(string), p.pollContextMsg(po)))

	return nil
}
func (p *Polls) queryRenameSuccess(po *types.Poll, newsFunc func(string)) error {
	//Update name

	_, err := p.db.Exec(`UPDATE queries SET name=$1 WHERE name=$2 AND guild=$3`, po.Data["new"], po.Data["query"], po.Guild)
	if err != nil {
		return err
	}

	//Update dependent queries
	var queries []types.Query
	var queriesdata []struct {
		name string       `db:name`
		data types.PgData `db:data`
	}

	err = p.db.Select(&queriesdata, `SELECT name,data FROM queries WHERE data->>'query'=$1 AND guild =$2`, po.Data["query"], po.Guild)
	if err == nil && len(queries) > 0 {
		//_, err = p.db.Exec("UPDATE queries set data=arraytojson")
	}

	//News
	newsFunc(fmt.Sprintf("📝 Renamed Query - **%s** %s", po.Data["new"].(string), p.pollContextMsg(po)))
	return nil
}
