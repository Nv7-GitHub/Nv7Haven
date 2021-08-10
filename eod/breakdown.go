package eod

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

type breakDownTree struct {
	lock      *sync.RWMutex
	added     map[string]types.Empty
	elemCache map[string]types.Element
	breakdown map[string]int // map[userid]count
	total     int
	tree      bool
}

func (b *breakDownTree) addElem(elem string) (bool, string) {
	_, exists := b.added[strings.ToLower(elem)]
	if exists {
		return true, ""
	}

	b.lock.RLock()
	el, exists := b.elemCache[strings.ToLower(elem)]
	b.lock.RUnlock()
	if !exists {
		return false, fmt.Sprintf("Element **%s** doesn't exist!", elem)
	}

	if b.tree {
		for _, par := range el.Parents {
			suc, err := b.addElem(par)
			if !suc {
				return suc, err
			}
		}
	}

	b.breakdown[el.Creator]++
	b.total++

	b.added[strings.ToLower(elem)] = types.Empty{}
	return true, ""
}

type breakDownSort struct {
	Count int
	Text  string
}

func (b *breakDownTree) getStringArr() []string {
	sorts := make([]breakDownSort, len(b.breakdown))
	i := 0
	for k, v := range b.breakdown {
		sorts[i] = breakDownSort{
			Count: v,
			Text:  fmt.Sprintf("<@%s> - %d, %s%%", k, v, formatFloat(float32(v)/float32(b.total)*100, 2)),
		}
		i++
	}

	sort.Slice(sorts, func(i, j int) bool {
		return sorts[i].Count > sorts[j].Count
	})

	out := make([]string, len(sorts))
	for i, v := range sorts {
		out[i] = v.Text
	}
	return out
}

func (b *EoD) elemBreakdownCmd(elem string, calcTree bool, m types.Msg, rsp types.Rsp) {
	rsp.Acknowledge()

	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("Guild isn't setup yet!")
		return
	}

	dat.Lock.RLock()
	el, exists := dat.ElemCache[strings.ToLower(elem)]
	dat.Lock.RUnlock()
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", elem))
		return
	}

	tree := &breakDownTree{
		lock:      dat.Lock,
		elemCache: dat.ElemCache,
		breakdown: make(map[string]int),
		added:     make(map[string]types.Empty),
		tree:      calcTree,
		total:     0,
	}
	suc, err := tree.addElem(el.Name)
	if !suc {
		rsp.ErrorMessage(err)
		return
	}

	b.newPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      fmt.Sprintf("%s Breakdown (%d)", el.Name, tree.total),
		PageGetter: b.invPageGetter,
		Items:      tree.getStringArr(),
	}, m, rsp)
}

func (b *EoD) catBreakdownCmd(catName string, calcTree bool, m types.Msg, rsp types.Rsp) {
	rsp.Acknowledge()

	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("Guild isn't setup yet!")
		return
	}

	dat.Lock.RLock()
	cat, exists := dat.CatCache[strings.ToLower(catName)]
	dat.Lock.RUnlock()
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Category **%s** doesn't exist!", catName))
		return
	}

	tree := &breakDownTree{
		lock:      dat.Lock,
		elemCache: dat.ElemCache,
		breakdown: make(map[string]int),
		added:     make(map[string]types.Empty),
		total:     0,
		tree:      calcTree,
	}

	for elem := range cat.Elements {
		suc, err := tree.addElem(elem)
		if !suc {
			rsp.ErrorMessage(err)
			return
		}
	}

	b.newPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      fmt.Sprintf("%s Breakdown (%d)", cat.Name, tree.total),
		PageGetter: b.invPageGetter,
		Items:      tree.getStringArr(),
	}, m, rsp)
}
