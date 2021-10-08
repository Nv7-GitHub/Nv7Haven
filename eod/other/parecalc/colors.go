package main

import (
	"fmt"
	"strconv"

	"github.com/lucasb-eyer/go-colorful"
)

var defaultColors = map[string]int{
	"air":   12764099, // #C2C3C3
	"earth": 11172162, // #AA7942
	"fire":  16749824, // #FF9500
	"water": 275455,   // #0433FF
}

func calcColor(elem string, gld Guild) (Guild, int) {
	col, exists := gld.Colors[elem]
	if exists {
		return gld, col
	}

	el := gld.Elements[elem]
	var color int
	cols := make([]int, len(el.Parents))
	for i, par := range el.Parents {
		gld, color = calcColor(par, gld)
		cols[i] = color
	}

	out, err := MixColors(cols)
	if err != nil {
		fmt.Println(el.Parents, cols)
	}
	handle(err)
	gld.Colors[elem] = out
	return gld, out
}

func recalcColors() {
	for id, gld := range glds {
		for name, col := range defaultColors {
			gld.Colors[name] = col
		}

		var col int
		for elem, el := range gld.Elements {
			gld, col = calcColor(elem, gld)
			el.Color = col
			gld.Elements[elem] = el
		}
		glds[id] = gld
	}
}

// Copied from util
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
