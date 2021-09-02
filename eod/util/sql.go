package util

import (
	"unicode"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func IsASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

var Wildcards = map[rune]types.Empty{
	'%': {},
	'*': {},
	'?': {},
	'[': {},
	']': {},
	'!': {},
	'-': {},
	'#': {},
	'^': {},
	'_': {},
}

func IsWildcard(s string) bool {
	for _, char := range s {
		_, exists := Wildcards[char]
		if exists {
			return true
		}
	}
	return false
}
