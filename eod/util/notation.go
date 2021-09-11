package util

import (
	"strings"
)

const notationChars = "BCDGHIJKLMNOPQRSTUVXYZ"

func Num2Char(num int) string {
	out := &strings.Builder{}
	encodeNum(num, out)
	return out.String()
}

func encodeNum(num int, out *strings.Builder) {
	if num/len(notationChars) != 0 {
		encodeNum(num/len(notationChars), out)
	}
	out.WriteByte(notationChars[num%len(notationChars)])
}
