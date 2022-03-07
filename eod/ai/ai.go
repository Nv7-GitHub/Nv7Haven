package ai

import (
	"sort"
	"strconv"
	"strings"

	"github.com/sasha-s/go-deadlock"
)

type AI struct {
	lock     *deadlock.RWMutex
	Links    map[int]*Probability
	Starters *Probability
}

func NewAI() *AI {
	return &AI{
		lock:     &deadlock.RWMutex{},
		Links:    make(map[int]*Probability),
		Starters: NewProbability(),
	}
}

func (a *AI) AddCombo(combo string, nolock bool) error {
	// Parse combo
	elems := strings.Split(combo, "+")
	vals := make([]int, len(elems))
	for i, elem := range elems {
		val, err := strconv.Atoi(elem)
		if err != nil {
			return err
		}
		vals[i] = val
	}
	sort.Ints(vals)

	// Add links
	for i, val := range vals {
		if i == 0 {
			a.Starters.Add(val, nolock) // Start link
		} else {
			a.AddLink(vals[i-1], val, nolock)
		}

		if i == len(vals)-1 {
			a.AddLink(val, -1, nolock) // End link
		}
	}

	return nil
}

func (a *AI) AddLink(start int, end int, nolock bool) {
	// Check if exists, make if it doesnt
	if !nolock {
		a.lock.RLock()
	}
	link, exists := a.Links[start]
	if !exists {
		if !nolock {
			a.lock.RUnlock()
			a.lock.Lock()
		}
		link = NewProbability()
		a.Links[start] = link
		if !nolock {
			a.lock.Unlock()
		}
	} else if !nolock {
		a.lock.RUnlock()
	}

	// Add to link
	link.Add(end, nolock)
}
