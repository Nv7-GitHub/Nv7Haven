package eodb

import (
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

type recalcCombo struct {
	Elems []int
	Elem3 int
	Done  bool
}

func (d *DB) Recalc() error {
	d.RLock()

	d.BeginTransaction()

	// Recalc trees
	done := make(map[int]types.Empty)
	for i := 1; i <= 4; i++ { // Starters done
		done[i] = types.Empty{}
	}

	// Build combination cache
	combos := make([]recalcCombo, len(d.combos))
	i := 0
	for k, v := range d.combos {
		items := make([]int, 0)
		for _, val := range strings.Split(k, "+") {
			num, err := strconv.Atoi(val)
			if err != nil {
				d.RUnlock()
				return err
			}
			items = append(items, num)
		}
		combos[i] = recalcCombo{
			Elems: items,
			Elem3: v,
		}
		i++
	}

	// Keep on recalcing until done
	changed := -1
	for changed != 0 {
		changed = 0

		// Loop through combos
		for i, comb := range combos {
			if comb.Done {
				continue
			}

			// Check if done
			_, exists := done[comb.Elem3]
			if exists {
				combos[i].Done = true
				continue
			}

			// If not done, check if can be done
			valid := true
			for _, elem := range comb.Elems {
				_, exists := done[elem] // Check if element has been done, if it hasnt then not ready
				if !exists {
					valid = false
					break
				}
			}

			// If its valid, do it
			if valid {
				// Save element
				el, res := d.GetElement(comb.Elem3, true)
				if !res.Exists {
					d.RUnlock()
					return fmt.Errorf("recalc: element %d does not exist", comb.Elem3)
				}

				// Update parents
				el.Parents = comb.Elems
				sort.Ints(el.Parents)

				// Update complexity & difficulty, element stats
				maxdiff := -1
				maxcomp := -1
				air := big.NewInt(0)
				earth := big.NewInt(0)
				fire := big.NewInt(0)
				water := big.NewInt(0)

				issame := false
				first := el.Parents[0]
				for _, elem := range el.Parents {
					el, res := d.GetElement(elem, true)
					if res.Exists { // Check complexity, difficulty
						if el.Difficulty > maxdiff {
							maxdiff = el.Difficulty
						}

						if el.Complexity > maxcomp {
							maxcomp = el.Complexity
						}
					}
					if elem != first { // Check if all are same (for difficulty)
						issame = false
					}
					air.Add(air, el.Air)
					earth.Add(earth, el.Earth)
					fire.Add(fire, el.Fire)
					water.Add(water, el.Water)
				}

				maxcomp++
				if !issame {
					maxdiff++
				}

				el.Difficulty = maxdiff
				el.Complexity = maxcomp
				el.Air = air
				el.Earth = earth
				el.Fire = fire
				el.Water = water

				// Save
				d.RUnlock()
				err := d.SaveElement(el)
				if err != nil {
					return err
				}
				d.RLock()

				// Update done
				done[comb.Elem3] = types.Empty{}
				changed++
			}
		}
	}

	// Recalc tree size
	for _, el := range d.Elements {
		done := make(map[int]types.Empty, el.TreeSize)
		res := d.recalcGetTreeSize(el.ID, done)
		if !res.Exists {
			d.RUnlock()
			return fmt.Errorf("recalc: %s", res.Message)
		}
		el.TreeSize = len(done)

		d.RUnlock()
		err := d.SaveElement(el)
		if err != nil {
			return err
		}
		d.RLock()
	}

	// Persist
	d.RUnlock()
	return d.CommitTransaction()
}

func (d *DB) recalcGetTreeSize(elem int, done map[int]types.Empty) types.GetResponse {
	_, exists := done[elem]
	if exists {
		return types.GetResponse{Exists: true}
	}
	el, res := d.GetElement(elem, true)
	if !res.Exists {
		return res
	}

	for _, parent := range el.Parents {
		res = d.recalcGetTreeSize(parent, done)
		if !res.Exists {
			return res
		}
	}

	done[elem] = types.Empty{}

	return types.GetResponse{Exists: true}
}
