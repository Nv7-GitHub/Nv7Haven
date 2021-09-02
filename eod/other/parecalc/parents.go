package main

import (
	"strings"
)

func recalcPars() {
	for id, gld := range glds {
		for _, elem := range starters {
			gld.Finished[strings.ToLower(elem)] = empty{}
		}

		changed := -1
		for changed != 0 {
			changed = 0
			newfinished := make(map[string]empty)
			for _, comb := range gld.Combos {
				// Lowercased when loading data
				//el3 := strings.ToLower(comb.Elem3)
				el3 := comb.Elem3
				_, exists := gld.Finished[el3]
				if exists {
					_, exists := newfinished[el3]
					if !exists {
						continue
					}
				}

				// Check if comb has all elems finished
				valid := true
				for _, elem := range comb.Elems {
					_, exists := gld.Finished[strings.ToLower(elem)]
					if !exists {
						valid = false
						break
					} else {
						_, exists = newfinished[el3]
						if exists {
							valid = false
							break
						}
					}
				}

				if valid {
					el := gld.Elements[el3]
					el.Parents = comb.Elems
					gld.Elements[el3] = el
					gld.Finished[el3] = empty{}
					newfinished[el3] = empty{}
					changed++
				}
			}
		}

		glds[id] = gld
	}
}
