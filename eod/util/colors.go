package util

import (
	"strconv"

	"github.com/lucasb-eyer/go-colorful"
)

func MixColors(colors []int) (int, error) {
	cls := make([]colorful.Color, len(colors))
	var err error
	for i, color := range colors {
		hex := strconv.FormatInt(int64(color), 16)
		if len(hex) < 6 {
			diff := 6 - len(hex)
			for i := 0; i < diff; i++ {
				hex = "0" + hex
			}
		}
		cls[i], err = colorful.Hex("#" + hex)
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
