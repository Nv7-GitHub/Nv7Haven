package eod

import (
	"fmt"
	"sort"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *EoD) foundCmd(elem string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
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

	b.newPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      fmt.Sprintf("%s Found (%d)", el.Name, len(out)),
		PageGetter: b.invPageGetter,
		Items:      out,
		User:       m.Author.ID,
	}, m, rsp)
}

func (b *EoD) categoriesCmd(elem string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	rsp.Acknowledge()

	el, res := dat.GetElement(elem)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	// Get Categories
	catsMap := make(map[catSortInfo]types.Empty)
	dat.Lock.RLock()
	for _, cat := range dat.Categories {
		_, exists := cat.Elements[el.Name]
		if exists {
			catsMap[catSortInfo{
				Name: cat.Name,
				Cnt:  len(cat.Elements),
			}] = types.Empty{}
		}
	}
	dat.Lock.RUnlock()
	cats := make([]catSortInfo, len(catsMap))
	i := 0
	for k := range catsMap {
		cats[i] = k
		i++
	}

	// Sort categories by count
	sort.Slice(cats, func(i, j int) bool {
		return cats[i].Cnt > cats[j].Cnt
	})

	// Convert to array
	out := make([]string, len(cats))
	for i, cat := range cats {
		out[i] = cat.Name
	}

	b.newPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      fmt.Sprintf("%s Categories (%d)", el.Name, len(out)),
		PageGetter: b.invPageGetter,
		Items:      out,
		User:       m.Author.ID,
	}, m, rsp)
}
