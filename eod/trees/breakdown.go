package trees

import (
	"fmt"
	"sort"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

type BreakDownTree struct {
	Added     map[int]types.Empty
	DB        *eodb.DB
	Breakdown map[string]int // map[userid]count
	Total     int
	Tree      bool
}

func (b *BreakDownTree) AddElem(elem int, noerror ...bool) (bool, string) {
	_, exists := b.Added[elem]
	if exists {
		return true, ""
	}

	el, res := b.DB.GetElement(elem)
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

	b.Added[elem] = types.Empty{}
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
