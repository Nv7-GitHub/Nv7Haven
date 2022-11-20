package polls

import (
	"fmt"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (p *Polls) elemImageSuccess(po *types.Poll, newsFunc func(string)) error {
	// Update image
	_, err := p.db.Exec(`UPDATE elements SET image=$1 WHERE id=$2 AND guild=$3`, po.Data["new"], int(po.Data["elem"].(float64)), po.Guild)
	if err != nil {
		return err
	}

	// Get name
	name, err := p.base.GetName(po.Guild, int(po.Data["elem"].(float64)))
	if err != nil {
		return err
	}

	// News
	word := "Changed"
	if po.Data["old"] == "" {
		word = "Added"
	}
	newsFunc(fmt.Sprintf("📸 %s Image - **%s** %s", word, name, p.pollContextMsg(po)))

	return nil
}
