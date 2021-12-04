package trees

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

// Tree calculator
type Tree struct {
	text   *strings.Builder
	rawTxt *strings.Builder
	calced map[int]types.Empty
	num    int

	db *eodb.DB
}

func (t *Tree) AddElem(elem int) (bool, string) {
	_, exists := t.calced[elem]
	if !exists {
		el, res := t.db.GetElement(elem)
		if !res.Exists {
			return false, res.Message
		}
		if len(el.Parents) == 1 {
			el.Parents = append(el.Parents, el.Parents[0])
		}
		for _, parent := range el.Parents {
			suc, msg := t.AddElem(parent)
			if !suc {
				return false, msg
			}
		}

		perf := &strings.Builder{}

		perf.WriteString("%d. ")
		params := make([]interface{}, len(el.Parents))
		for i, val := range el.Parents {
			if i == 0 {
				perf.WriteString("%s")
			} else {
				perf.WriteString(" + %s")
			}
			el, _ := t.db.GetElement(val)
			params[i] = interface{}(el.Name)
		}
		params = append([]interface{}{t.num}, params...)
		params = append(params, el.Name)
		if len(el.Parents) >= 2 {
			p := perf.String()
			fmt.Fprintf(t.text, p+" = **%s**\n", params...)
			fmt.Fprintf(t.rawTxt, p+" = %s\n", params...)
			t.num++
		}
		t.calced[elem] = types.Empty{}
	}
	return true, ""
}

// Tree calculation utilities
func CalcTree(db *eodb.DB, elem int) (string, bool, string) {
	// Commented out code is for profiling

	/*runtime.GC()
	cpuprof, _ := os.Create("cpuprof.pprof")
	pprof.StartCPUProfile(cpuprof)*/

	t := Tree{
		text:   &strings.Builder{},
		rawTxt: &strings.Builder{},
		calced: make(map[int]types.Empty),
		num:    1,
		db:     db,
	}
	suc, msg := t.AddElem(elem)

	/*pprof.StopCPUProfile()
	memprof, _ := os.Create("memprof.pprof")
	_ = pprof.WriteHeapProfile(memprof)*/

	text := t.text.String()
	if len(text) > 2000 {
		return t.rawTxt.String(), suc, msg
	}

	return text, suc, msg
}

func CalcTreeCat(db *eodb.DB, elems map[int]types.Empty) (string, bool, string) {
	// Commented out code is for profiling

	/*runtime.GC()
	cpuprof, _ := os.Create("cpuprof.pprof")
	pprof.StartCPUProfile(cpuprof)*/

	t := Tree{
		text:   &strings.Builder{},
		rawTxt: &strings.Builder{},
		calced: make(map[int]types.Empty),
		num:    1,
		db:     db,
	}
	for elem := range elems {
		suc, msg := t.AddElem(elem)
		if !suc {
			return "", false, msg
		}
	}

	/*pprof.StopCPUProfile()
	memprof, _ := os.Create("memprof.pprof")
	_ = pprof.WriteHeapProfile(memprof)*/

	text := t.text.String()
	if len(text) > 2000 {
		return t.rawTxt.String(), true, ""
	}

	return text, true, ""
}
