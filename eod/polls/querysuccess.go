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
		_, err := p.db.Exec(`INSERT INTO queries (guild, name, creator, createdon, kind, data, image, comment, imager, colorer, commenter, color) VALUES ($1, $2, $3, $4, $5, $6, $7, $7, $7, $7, $7, $8)`, po.Guild, po.Data["query"], po.Creator, time.Now(), po.Data["kind"], types.PgData(po.Data["data"].(map[string]any)), "", 0)
		if err != nil {
			return err
		}
	}

	// News
	news(fmt.Sprintf("🧮 Created Query - **%s** %s", po.Data["query"], p.pollContextMsg(po)))
	return nil
}
