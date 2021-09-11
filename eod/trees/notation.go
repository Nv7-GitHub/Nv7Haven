package trees

import (
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

var elemNotations = map[string]string{
	"air":   "A",
	"earth": "E",
	"fire":  "F",
	"water": "W",
}

type notationTree struct {
	*strings.Builder

	added map[string]string
	num   int

	dat types.ServerData
}

func NewNotationTree(dat types.ServerData) *notationTree {
	return &notationTree{
		added:   make(map[string]string),
		Builder: &strings.Builder{},
		num:     0,
		dat:     dat,
	}
}

func (n *notationTree) AddElem(elem string) (string, bool) {
	val, exists := n.added[elem]
	if exists {
		return val, true
	}

	el, res := n.dat.GetElement(elem)
	if !res.Exists {
		return res.Message, res.Exists
	}

	for _, par := range el.Parents {
		_, exists := elemNotations[strings.ToLower(par)]
		if !exists {
			msg, suc := n.AddElem(par)
			if !suc {
				return msg, suc
			}
		}
	}

	n.WriteString("(")
	for _, par := range el.Parents {
		val, exists := elemNotations[strings.ToLower(par)]
		if !exists {
			msg, suc := n.AddElem(par)
			if !suc {
				return msg, suc
			}
			if len(msg) > 1 {
				n.WriteString("{")
				n.WriteString(msg)
				n.WriteString("}")
			} else {
				n.WriteString(msg)
			}
		} else {
			n.WriteString(val)
		}
	}
	n.WriteString(")")

	val = util.Num2Char(n.num)
	n.num++

	n.WriteString("[")
	n.WriteString(val)
	n.WriteString("]")

	n.added[strings.ToLower(elem)] = val
	return val, true
}
