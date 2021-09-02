package main

import (
	"strings"
)

func (g *Guild) CalcElemStats(elem string) {
	elem = strings.ToLower(elem)

	_, exists := g.Finished[elem]
	if exists {
		return
	}

	el, exists := g.Elements[elem]
	if !exists {
		g.Finished[elem] = empty{}
		return
	}

	if len(el.Parents) == 0 {
		g.Finished[elem] = empty{}
		return
	}

	for _, par := range el.Parents {
		g.CalcElemStats(par)
	}

	unique := false
	first := el.Parents[0]

	maxDiff := 0
	maxComp := 0
	for _, par := range el.Parents {
		if par != first {
			unique = true
		}

		parEl, exists := g.Elements[par]
		if exists {
			if parEl.Complexity > maxComp {
				maxComp = parEl.Complexity
			}

			if parEl.Difficulty > maxDiff {
				maxDiff = parEl.Difficulty
			}
		}
	}

	el.Complexity = maxComp + 1
	el.Difficulty = maxDiff
	if unique {
		el.Difficulty++
	}

	g.Finished[elem] = empty{}
	g.Elements[elem] = el
}

func recalcStats() {
	for id, gld := range glds {
		gld.Finished = make(map[string]empty)
		for _, elem := range starters {
			gld.Finished[strings.ToLower(elem)] = empty{}
		}

		for _, elem := range gld.Elements {
			gld.CalcElemStats(elem.Name)
		}

		glds[id] = gld
	}
}
