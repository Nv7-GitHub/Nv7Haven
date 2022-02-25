package util

import (
	"fmt"
	"math"
	"math/big"
	"sort"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func FormatFloat(num float32, prc int) string {
	var (
		zero, dot = "0", "."

		str = fmt.Sprintf("%."+strconv.Itoa(prc)+"f", num)
	)

	return strings.TrimRight(strings.TrimRight(str, zero), dot)
}

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

func Elems2Txt(elems []string) string {
	for i, elem := range elems {
		elems[i] = strings.ToLower(elem)
	}
	sort.Strings(elems)
	return strings.Join(elems, "+")
}

var noObscure = map[rune]types.Empty{
	' ': {},
	'.': {},
	'-': {},
	'_': {},
}

func Obscure(val string) string {
	out := make([]rune, len([]rune(val)))
	i := 0
	for _, char := range val {
		_, exists := noObscure[char]
		if exists {
			out[i] = char
		} else {
			out[i] = '?'
		}
		i++
	}
	return string(out)
}

var ten = big.NewInt(10)
var superscript = map[string]string{
	"0": "⁰",
	"1": "¹",
	"2": "²",
	"3": "³",
	"4": "⁴",
	"5": "⁵",
	"6": "⁶",
	"7": "⁷",
	"8": "⁸",
	"9": "⁹",
}

func FormatBigInt(b *big.Int) string {
	if b.IsInt64() {
		return FormatInt(int(b.Int64()))
	}

	// Get number of digits
	fac := math.Log(2) / math.Log(10)
	count := int64(fac * float64(b.BitLen()))
	if big.NewInt(0).Exp(ten, big.NewInt(count-1), nil).Cmp(b) < 0 {
		count++
	}

	// Get val
	tenCnt := count - 1
	div := big.NewInt(0).Exp(ten, big.NewInt(tenCnt), nil)
	v := big.NewFloat(0).Quo(big.NewFloat(0).SetInt(b), big.NewFloat(0).SetInt(div)) // b/tenCnt
	val, _ := v.Float64()

	// Format
	exp := strconv.Itoa(int(tenCnt))
	for k, v := range superscript {
		exp = strings.ReplaceAll(exp, k, v)
	}

	return fmt.Sprintf("%0.6f×10%s", val, exp)
}
