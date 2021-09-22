package util

import "strings"

func TrimArray(vals []string) []string {
	var trimmedVals []string
	for _, val := range vals {
		if len(val) > 0{
			trimmedVals = append(trimmedVals,strings.TrimSpace(val))
		}	
	}
	return trimmedVals
}

func EscapeElement(elem string) string {
	return strings.ReplaceAll(elem, "\\", "\\\\")
}
