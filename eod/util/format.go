package util

import (
	"strconv"
	"strings"
	"unicode"
)

var smallWords = map[string]struct{}{
	"of":  {},
	"an":  {},
	"on":  {},
	"the": {},
	"to":  {},
}

func Capitalize(s string) string {
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

func FormatHex(color int) string {
	hex := strconv.FormatInt(int64(color), 16)
	if len(hex) < 6 {
		diff := 6 - len(hex)
		for i := 0; i < diff; i++ {
			hex = "0" + hex
		}
	}
	return "#" + hex
}
