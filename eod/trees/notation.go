package trees

import (
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

var elemNotations = map[string]string{
	"air":   "A",
	"earth": "E",
	"fire":  "F",
	"water": "W",
}

type notationTree struct {
	notations map[string]string

	dat types.ServerData
}

func NewNotationTree(dat types.ServerData) *notationTree {
	return &notationTree{
		notations: make(map[string]string),
		dat:       dat,
	}
}

//IMPORTANT: RLock before getting notation
func (n *notationTree) GetNotation(elem string) (string, bool) {
	elem = strings.ToLower(elem)

	notation, exists := elemNotations[elem]
	if exists {
		return notation, true
	}
	notation, exists = n.notations[elem]
	if exists {
		return notation, true
	}

	el, res := n.dat.GetElement(elem, true)
	if !res.Exists {
		return res.Message, res.Exists
	}

	out := &strings.Builder{}
	for _, par := range el.Parents {
		notation, suc := n.GetNotation(par)
		if !suc {
			return notation, suc
		}

		out.WriteString("(")
		out.WriteString(notation)
		out.WriteString(")")
	}

	n.notations[elem] = out.String()
	return out.String(), true
}
