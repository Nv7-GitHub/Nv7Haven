package elements

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *Elements) FoundCmd(elem string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}

	rsp.Acknowledge()

	el, res := db.GetElementByName(elem)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	items := make(map[string]types.Empty)
	db.RLock()
	for _, inv := range db.Invs() {
		if inv.Contains(el.ID) {
			items[inv.User] = types.Empty{}
		}
	}
	db.RUnlock()

	out := make([]string, len(items))
	i := 0
	for k := range items {
		out[i] = k
		i++
	}
	sort.Slice(out, func(i, j int) bool {
		int1, err1 := strconv.Atoi(out[i])
		int2, err2 := strconv.Atoi(out[j])
		if err1 != nil && err2 != nil {
			return int1 < int2
		}
		return out[i] < out[j]
	})
	for i, v := range out {
		out[i] = fmt.Sprintf("<@%s>", v)
	}

	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind: types.PageSwitchInv,
		Title: db.Config.LangProperty("ElemFound", map[string]interface{}{
			"Element": el.Name,
			"Count":   len(out),
		}),
		PageGetter: b.base.InvPageGetter,
		Items:      out,
		User:       m.Author.ID,
		Thumbnail:  el.Image,
	}, m, rsp)
}
