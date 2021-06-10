package eod

import (
	"fmt"
	"strings"
)

// Tree calculator
type tree struct {
	text      *strings.Builder
	rawTxt    *strings.Builder
	elemCache map[string]element
	calced    map[string]empty
	num       int
}

func (t *tree) addElem(elem string) (bool, string) {
	_, exists := t.calced[strings.ToLower(elem)]
	if !exists {
		el, exists := t.elemCache[strings.ToLower(elem)]
		if !exists {
			return false, elem
		}
		if len(el.Parents) == 1 {
			el.Parents = append(el.Parents, el.Parents[0])
		}
		for _, parent := range el.Parents {
			if len(strings.TrimSpace(parent)) == 0 {
				continue
			}
			suc, msg := t.addElem(parent)
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
			params[i] = interface{}(t.elemCache[strings.ToLower(val)].Name)
		}
		params = append([]interface{}{t.num}, params...)
		params = append(params, el.Name)
		if len(el.Parents) >= 2 {
			p := perf.String()
			fmt.Fprintf(t.text, p+" = **%s**\n", params...)
			fmt.Fprintf(t.rawTxt, p+" = %s\n", params...)
			t.num++
		}
		t.calced[strings.ToLower(elem)] = empty{}
	}
	return true, ""
}

// Tree calculation utilities
func calcTree(elemCache map[string]element, elem string) (string, bool, string) {
	// Commented out code is for profiling

	/*runtime.GC()
	cpuprof, _ := os.Create("cpuprof.pprof")
	pprof.StartCPUProfile(cpuprof)*/

	t := tree{
		text:      &strings.Builder{},
		rawTxt:    &strings.Builder{},
		elemCache: elemCache,
		calced:    make(map[string]empty),
		num:       1,
	}
	suc, msg := t.addElem(elem)

	/*pprof.StopCPUProfile()
	memprof, _ := os.Create("memprof.pprof")
	_ = pprof.WriteHeapProfile(memprof)*/

	text := t.text.String()
	if len(text) > 2000 {
		return t.rawTxt.String(), suc, msg
	}

	return text, suc, msg
}

func calcTreeCat(elemCache map[string]element, elems map[string]empty) (string, bool, string) {
	// Commented out code is for profiling

	/*runtime.GC()
	cpuprof, _ := os.Create("cpuprof.pprof")
	pprof.StartCPUProfile(cpuprof)*/

	t := tree{
		text:      &strings.Builder{},
		rawTxt:    &strings.Builder{},
		elemCache: elemCache,
		calced:    make(map[string]empty),
		num:       1,
	}
	for elem := range elems {
		suc, msg := t.addElem(elem)
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
