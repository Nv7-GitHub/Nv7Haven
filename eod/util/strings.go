package util

import (
	"sort"
	"strings"
	"unicode"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func ToTitle(s string) string {
	words := strings.Split(strings.ToLower(s), " ")
	for i, word := range words {
		if len(word) < 1 {
			continue
		}
		w := []rune(word)
		ind := -1

		if w[0] > unicode.MaxASCII {
			continue
		}

		if i == 0 {
			ind = 0
		} else {
			_, exists := smallWords[word]
			if !exists {
				ind = 0
			}
		}

		if w[0] == '(' && len(word) > 1 {
			ind = 1
		}

		if ind != -1 {
			w[ind] = []rune(strings.ToUpper(string([]rune(word)[ind])))[0]
			words[i] = string(w)
		}
	}
	return strings.Join(words, " ")
}

var smallWords = map[string]types.Empty{
	"of":  {},
	"an":  {},
	"on":  {},
	"the": {},
	"to":  {},
}

func JoinTxt(elemDat map[string]types.Empty, ending string) string {
	elems := make([]string, len(elemDat))
	i := 0
	for k := range elemDat {
		elems[i] = k
		i++
	}
	sort.Strings(elems)

	out := ""
	for i, elem := range elems {
		out += elem
		if i != len(elems)-1 && len(elems) != 2 {
			out += ", "
		} else if i != len(elems)-1 {
			out += " "
		}

		if i == len(elems)-2 {
			out += ending + " "
		}
	}

	return out
}
