package util

import (
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func TrimArray(vals []string) []string {
	for i, val := range vals {
		vals[i] = strings.TrimSpace(val)
	}
	return vals
}

func EscapeElement(elem string) string {
	return strings.ReplaceAll(elem, "\\", "\\\\")
}

func RemoveDuplicates(elems []string) []string {
	mp := make(map[string]types.Empty, len(elems))
	for _, elem := range elems {
		mp[elem] = types.Empty{}
	}
	out := make([]string, len(mp))
	i := 0
	for k := range mp {
		out[i] = k
		i++
	}
	return out
}
