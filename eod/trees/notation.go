package trees

import (
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

var elemNotations = map[int]string{
	1: "A",
	2: "E",
	3: "F",
	4: "W",
}

type notationTree struct {
	*strings.Builder

	added map[int]string
	num   int

	db *eodb.DB
}

func NewNotationTree(db *eodb.DB) *notationTree {
	return &notationTree{
		added:   make(map[int]string),
		Builder: &strings.Builder{},
		num:     0,
		db:      db,
	}
}

func (n *notationTree) AddElem(elem int) (string, bool) {
	val, exists := n.added[elem]
	if exists {
		return val, true
	}

	el, res := n.db.GetElement(elem)
	if !res.Exists {
		return res.Message, res.Exists
	}
	if len(el.Parents) == 1 {
		el.Parents = append(el.Parents, el.Parents[0])
	}

	for _, par := range el.Parents {
		_, exists := elemNotations[par]
		if !exists {
			msg, suc := n.AddElem(par)
			if !suc {
				return msg, suc
			}
		}
	}

	n.WriteString("(")
	for _, par := range el.Parents {
		val, exists := elemNotations[par]
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

	n.added[elem] = val
	return val, true
}
