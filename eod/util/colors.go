package util

import (
	"math"
	"strconv"

	"github.com/lucasb-eyer/go-colorful"
)

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

func MixColors(colors []int) (int, error) {
	cls := make([]colorful.Color, len(colors))
	var err error
	for i, color := range colors {
		cls[i], err = colorful.Hex(FormatHex(color))
		if err != nil {
			return 0, err
		}
	}

	var h, s, v float64
	for _, val := range cls {
		hv, sv, vv := val.Hsv()
		h += hv
		s += sv
		v += vv
	}
	length := float64(len(colors))
	h /= length
	s /= length
	v /= length

	out := colorful.Hsv(h, s, v)
	outv, err := strconv.ParseInt(out.Hex()[1:], 16, 64)
	return int(outv), err
}

// map[emoji]HSV
var emojiColors = map[[3]float64]string{
	{0, 0, 0}:         "⚫",
	{240, 1, 1}:       "🔵",
	{180, 1, 1}:       "🔵",
	{32, 0.61, 0.67}:  "🟤",
	{120, 1, 1}:       "🟢",
	{35, 1, 1}:        "🟠",
	{301, 0.78, 0.58}: "🟣",
	{0, 1, 1}:         "🔴",
	{317, 0.3, 1}:     "🔴",
	{0, 0, 1}:         "⚪",
	{59, 1, 1}:        "🟡",
	{112, 55, 0}:      "🟤",
}

func GetEmoji(color int) (string, error) {
	col, err := colorful.Hex(FormatHex(color))
	if err != nil {
		return "", err
	}
	h, s, v := col.Hsv()

	curr := ""
	var dist float64 = -1

	i := 0
	for k, val := range emojiColors {
		currDist := math.Abs((h-k[0])/120) + math.Abs((s-k[1])/1) + math.Abs((v-k[2])/1) // H is only divided by 120 instead of 360 so that Hue is biased more
		if currDist < dist || dist == -1 {
			curr = val
			dist = currDist
		}
		i++
	}

	return curr, nil
}
