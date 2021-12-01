package util

import (
	"sort"
	"strconv"
	"strings"
)

func FormatCombo(elems []int) string {
	sort.Sort(sort.IntSlice(elems))

	out := &strings.Builder{}
	for i, v := range elems {
		out.WriteString(strconv.Itoa(v))
		if i != len(elems)-1 {
			out.WriteString("+")
		}
	}
	return out.String()
}
