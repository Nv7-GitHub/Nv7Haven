package main

import (
	"fmt"
)

type elemTree struct {
	added map[string]empty
	gld   Guild
}

func (e *elemTree) addElem(name string) {
	_, exists := e.added[name]
	if exists {
		return
	}

	el, exists := e.gld.Elements[name]
	if !exists {
		panic(fmt.Errorf("unknown element %s", name))
	}
	for _, par := range el.Parents {
		e.addElem(par)
	}
	e.added[name] = empty{}
}

func recalcTreeSize() {
	for id, gld := range glds {
		for elem := range gld.Elements {
			tree := &elemTree{gld: gld}
			tree.addElem(elem)

			el := gld.Elements[elem]
			el.TreeSize = len(tree.added)

			glds[id].Elements[elem] = el
		}
	}
}
