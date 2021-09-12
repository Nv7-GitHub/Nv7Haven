package util

import (
	"fmt"
	"sort"
	"strconv"
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

var lock = &sync.RWMutex{}

var SortChoices = []*discordgo.ApplicationCommandOptionChoice{
	{
		Name:  "Name",
		Value: "name",
	},
	{
		Name:  "Length",
		Value: "length",
	},
	{
		Name:  "Date Created",
		Value: "createdon",
	},
	{
		Name:  "Complexity",
		Value: "complexity",
	},
	{
		Name:  "Difficulty",
		Value: "difficulty",
	},
	{
		Name:  "Used In",
		Value: "usedin",
	},
	{
		Name:  "Creator",
		Value: "creator",
	},
	{
		Name:  "ID",
		Value: "id",
	},
}

var sorts = map[string]func(a, b string, dat types.ServerData) bool{
	"length": func(a, b string, dat types.ServerData) bool {
		return len(a) < len(b)
	},
	"name": func(a, b string, dat types.ServerData) bool {
		return CompareStrings(a, b)
	},
	"createdon": func(a, b string, dat types.ServerData) bool {
		el1, res := dat.GetElement(a, true)
		el2, res2 := dat.GetElement(b, true)
		if !res.Exists || !res2.Exists {
			return false
		}
		return el1.CreatedOn.Before(el2.CreatedOn)
	},
	"id": func(a, b string, dat types.ServerData) bool {
		el1, res := dat.GetElement(a, true)
		el2, res2 := dat.GetElement(b, true)
		if !res.Exists || !res2.Exists {
			return false
		}
		return el1.ID < el2.ID
	},
	"complexity": func(a, b string, dat types.ServerData) bool {
		el1, res := dat.GetElement(a, true)
		el2, res2 := dat.GetElement(b, true)
		if !res.Exists || !res2.Exists {
			return false
		}
		return el1.Complexity < el2.Complexity
	},
	"difficulty": func(a, b string, dat types.ServerData) bool {
		el1, res := dat.GetElement(a, true)
		el2, res2 := dat.GetElement(b, true)
		if !res.Exists || !res2.Exists {
			return false
		}
		return el1.Difficulty < el2.Difficulty
	},
	"usedin": func(a, b string, dat types.ServerData) bool {
		el1, res := dat.GetElement(a, true)
		el2, res2 := dat.GetElement(b, true)
		if !res.Exists || !res2.Exists {
			return false
		}
		return el1.UsedIn < el2.UsedIn
	},
	"creator": func(a, b string, dat types.ServerData) bool {
		el1, res := dat.GetElement(a, true)
		el2, res2 := dat.GetElement(b, true)
		if !res.Exists || !res2.Exists {
			return false
		}
		return el1.Creator < el2.Creator
	},
}

var getters = map[string]func(el types.Element) string{
	"createdon": func(el types.Element) string {
		return fmt.Sprintf(" - <t:%d>", el.CreatedOn.Unix())
	},
	"id": func(el types.Element) string {
		return fmt.Sprintf(" - #%d", el.ID)
	},
	"complexity": func(el types.Element) string {
		return fmt.Sprintf(" - %d", el.Complexity)
	},
	"difficulty": func(el types.Element) string {
		return fmt.Sprintf(" - %d", el.Difficulty)
	},
	"usedin": func(el types.Element) string {
		return fmt.Sprintf(" - %d", el.UsedIn)
	},
	"creator": func(el types.Element) string {
		return fmt.Sprintf(" - <@%s>", el.Creator)
	},
}

// Less
func CompareStrings(a, b string) bool {
	fl1, err := strconv.ParseFloat(a, 32)
	fl2, err2 := strconv.ParseFloat(b, 32)
	if err == nil && err2 == nil {
		return fl1 < fl2
	}
	return a < b
}

func SortElemList(elems []string, sortName string, dat types.ServerData) {
	lock.RLock()
	sorter := sorts[sortName]
	lock.RUnlock()

	dat.Lock.RLock()
	sort.Slice(elems, func(i, j int) bool {
		return sorter(elems[i], elems[j], dat)
	})

	lock.RLock()
	getter, exists := getters[sortName]
	lock.RUnlock()
	if exists {
		for i, val := range elems {
			el, res := dat.GetElement(val, true)
			if res.Exists {
				elems[i] = val + getter(el)
			}
		}
	}
	dat.Lock.RUnlock()
}
