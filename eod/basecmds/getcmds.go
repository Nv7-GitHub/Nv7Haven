package basecmds

import (
	"fmt"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *BaseCmds) FoundCmd(elem string, m types.Msg, rsp types.Rsp) {
	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		return
	}

	rsp.Acknowledge()

	el, res := dat.GetElement(elem)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	var foundCnt int
	err := b.db.QueryRow(`SELECT COUNT(1) as cnt FROM eod_inv WHERE guild=? AND (JSON_EXTRACT(inv, CONCAT('$."', LOWER(?), '"')) IS NOT NULL)`, m.GuildID, el.Name).Scan(&foundCnt)
	if rsp.Error(err) {
		return
	}

	found, err := b.db.Query(`SELECT user as cnt FROM eod_inv WHERE guild=? AND (JSON_EXTRACT(inv, CONCAT('$."', LOWER(?), '"')) IS NOT NULL)`, m.GuildID, el.Name)
	if rsp.Error(err) {
		return
	}
	defer found.Close()

	out := make([]string, foundCnt)
	i := 0

	var user string
	for found.Next() {
		err = found.Scan(&user)
		if rsp.Error(err) {
			return
		}

		out[i] = fmt.Sprintf("<@%s>", user)
		i++
	}

	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      fmt.Sprintf("%s Found (%d)", el.Name, len(out)),
		PageGetter: b.base.InvPageGetter,
		Items:      out,
		User:       m.Author.ID,
	}, m, rsp)
}
