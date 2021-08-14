package trees

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

// Tree calculator
type Tree struct {
	text   *strings.Builder
	rawTxt *strings.Builder
	calced map[string]types.Empty
	num    int

	dat types.ServerData
}

func (t *Tree) AddElem(elem string) (bool, string) {
	_, exists := t.calced[strings.ToLower(elem)]
	if !exists {
		el, res := t.dat.GetElement(elem)
		if !res.Exists {
			return false, elem
		}
		if len(el.Parents) == 1 {
			el.Parents = append(el.Parents, el.Parents[0])
		}
		for _, parent := range el.Parents {
			if len(strings.TrimSpace(parent)) == 0 {
				continue
			}
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
			el, _ := t.dat.GetElement(val)
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
		t.calced[strings.ToLower(elem)] = types.Empty{}
	}
	return true, ""
}

// Tree calculation utilities
func CalcTree(dat types.ServerData, elem string) (string, bool, string) {
	// Commented out code is for profiling

	/*runtime.GC()
	cpuprof, _ := os.Create("cpuprof.pprof")
	pprof.StartCPUProfile(cpuprof)*/

	t := Tree{
		text:   &strings.Builder{},
		rawTxt: &strings.Builder{},
		calced: make(map[string]types.Empty),
		num:    1,
		dat:    dat,
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

func CalcTreeCat(dat types.ServerData, elems types.Container) (string, bool, string) {
	// Commented out code is for profiling

	/*runtime.GC()
	cpuprof, _ := os.Create("cpuprof.pprof")
	pprof.StartCPUProfile(cpuprof)*/

	t := Tree{
		text:   &strings.Builder{},
		rawTxt: &strings.Builder{},
		calced: make(map[string]types.Empty),
		num:    1,
		dat:    dat,
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
