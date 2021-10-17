package elements

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *Elements) FoundCmd(elem string, m types.Msg, rsp types.Rsp) {
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

	items := make(map[string]types.Empty)
	for _, inv := range dat.Inventories {
		if inv.Elements.Contains(el.Name) {
			items[inv.User] = types.Empty{}
		}
	}

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
		Kind:       types.PageSwitchInv,
		Title:      fmt.Sprintf("%s Found (%d)", el.Name, len(out)),
		PageGetter: b.base.InvPageGetter,
		Items:      out,
		User:       m.Author.ID,
	}, m, rsp)
}
