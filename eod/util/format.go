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
func HexToRGB(hex string) (int, int, int) {
	colorhex := strings.Trim(hex, "#")
	colorr := ""
	colorr += string(colorhex[0]) + string(colorhex[1])
	colorg := ""
	colorg += string(colorhex[2]) + string(colorhex[3])
	colorb := ""
	colorb += string(colorhex[4]) + string(colorhex[5])
	colorbint, _ := strconv.ParseInt(colorb, 16, 64)
	colorgint, _ := strconv.ParseInt(colorg, 16, 64)
	colorrint, _ := strconv.ParseInt(colorr, 16, 64)
	return int(colorrint), int(colorgint), int(colorbint)

}
func RGBToHex(R int, G int, B int) string {

	hexr := strconv.FormatInt(int64(R), 16)
	if len(hexr) < 2 {
		hexr = "0" + hexr
	}
	hexg := strconv.FormatInt(int64(G), 16)
	if len(hexg) < 2 {
		hexg = "0" + hexg
	}
	hexb := strconv.FormatInt(int64(B), 16)
	if len(hexb) < 2 {
		hexb = "0" + hexb
	}

	return "#" + hexr + hexg + hexb
}
