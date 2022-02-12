package eodb

import (
	"container/list"
	"fmt"
	"runtime/pprof"
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
	defer d.RUnlock()

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
					return fmt.Errorf("recalc: element %d does not exist", comb.Elem3)
				}

				// Update parents
				el.Parents = comb.Elems
				sort.Ints(el.Parents)

				// Update complexity & difficulty
				maxdiff := -1
				maxcomp := -1

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
				}

				maxcomp++
				if !issame {
					maxdiff++
				}

				el.Difficulty = maxdiff
				el.Complexity = maxcomp

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
		size, res := d.recalcGetTreeSize(el.ID)
		if !res.Exists {
			return fmt.Errorf("recalc: %s", res.Message)
		}
		el.TreeSize = size

		d.RUnlock()
		err := d.SaveElement(el)
		if err != nil {
			return err
		}
		d.RLock()
	}

	pprof.StopCPUProfile()

	// Done!
	return nil
}

func (d *DB) recalcGetTreeSize(elem int) (int, types.GetResponse) {
	todo := list.New()
	todo.PushBack(elem)
	size := 0

	// Calc tree size
	done := make(map[int]types.Empty)
	for todo.Len() > 0 {
		elem := todo.Remove(todo.Front()).(int)
		_, exists := done[elem]
		if exists {
			continue
		}

		el, res := d.GetElement(elem, true)
		if !res.Exists {
			return 0, res
		}

		// Update tree size
		size++

		// Add parents to TODO
		for _, parent := range el.Parents {
			_, exists := done[parent]
			if !exists {
				todo.PushBack(parent)
			}
		}

		// Done
		done[elem] = types.Empty{}
	}

	// Free up memory to make GC do less work
	for todo.Len() > 0 {
		todo.Remove(todo.Front())
	}

	return size, types.GetResponse{Exists: true}
}
