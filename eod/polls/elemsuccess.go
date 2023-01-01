package polls

import (
	"fmt"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (p *Polls) elemImageSuccess(po *types.Poll, newsFunc func(string)) error {
	// Update image
	_, err := p.db.Exec(`UPDATE elements SET image=$1, imager=$2 WHERE id=$3 AND guild=$4`, po.Data["new"], po.Creator, int(po.Data["elem"].(float64)), po.Guild)
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

func (p *Polls) elemMarkSuccess(po *types.Poll, newsFunc func(string)) error {
	// Update image
	_, err := p.db.Exec(`UPDATE elements SET comment=$1, commenter=$2 WHERE id=$3 AND guild=$4`, po.Data["new"], po.Creator, int(po.Data["elem"].(float64)), po.Guild)
	if err != nil {
		return err
	}

	// Get name
	name, err := p.base.GetName(po.Guild, int(po.Data["elem"].(float64)))
	if err != nil {
		return err
	}

	// News
	newsFunc(fmt.Sprintf("📝 Signed - **%s** %s", name, p.pollContextMsg(po)))

	return nil
}

func (p *Polls) elemColorSuccess(po *types.Poll, newsFunc func(string)) error {
	// Update image
	_, err := p.db.Exec(`UPDATE elements SET color=$1, colorer=$2 WHERE id=$3 AND guild=$4`, po.Data["new"], po.Creator, int(po.Data["elem"].(float64)), po.Guild)
	if err != nil {
		return err
	}

	// Get name
	name, err := p.base.GetName(po.Guild, int(po.Data["elem"].(float64)))
	if err != nil {
		return err
	}

	// News
	newsFunc(fmt.Sprintf("🎨 Colored - **%s** %s", name, p.pollContextMsg(po)))

	return nil
}
