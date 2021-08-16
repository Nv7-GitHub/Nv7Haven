package util

import "strings"

func TrimArray(vals []string) []string {
	for i, val := range vals {
		vals[i] = strings.TrimSpace(val)
	}
	return vals
}
