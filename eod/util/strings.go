package util

import (
	"strings"
	"unicode"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func ToTitle(s string) string {
	words := strings.Split(strings.ToLower(s), " ")
	for i, word := range words {
		if len(word) < 1 || (w[0] == '(' && len(word) < 2) {
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

		if w[0] == '(' {
			ind = 1
		}

		if ind != -1 {
			w[ind] = rune(strings.ToUpper(string(word[ind]))[0])
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
