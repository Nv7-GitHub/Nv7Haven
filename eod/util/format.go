package util

import (
	"strconv"
	"strings"
	"unicode"
)

// https://stackoverflow.com/a/31046325/11388343
func FormatInt(n int) string {
	in := strconv.FormatInt(int64(n), 10)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits-- // First character is the - sign (not a digit)
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}

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
