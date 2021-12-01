package trees

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

type BreakDownTree struct {
	Added     map[string]types.Empty
	Dat       types.ServerDat
	Breakdown map[string]int // map[userid]count
	Total     int
	Tree      bool
}

func (b *BreakDownTree) AddElem(elem string, noerror ...bool) (bool, string) {
	_, exists := b.Added[strings.ToLower(elem)]
	if exists {
		return true, ""
	}

	el, res := b.Dat.GetElement(elem)
	if !res.Exists {
		return false, res.Message
	}

	if b.Tree {
		for _, par := range el.Parents {
			suc, err := b.AddElem(par)
			if !suc && len(noerror) == 0 {
				return suc, err
			}
		}
	}

	b.Breakdown[el.Creator]++
	b.Total++

	b.Added[strings.ToLower(elem)] = types.Empty{}
	return true, ""
}

type breakDownSort struct {
	Count int
	Text  string
}

func (b *BreakDownTree) GetStringArr() []string {
	sorts := make([]breakDownSort, len(b.Breakdown))
	i := 0
	for k, v := range b.Breakdown {
		sorts[i] = breakDownSort{
			Count: v,
			Text:  fmt.Sprintf("<@%s> - %d, %s%%", k, v, util.FormatFloat(float32(v)/float32(b.Total)*100, 2)),
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
