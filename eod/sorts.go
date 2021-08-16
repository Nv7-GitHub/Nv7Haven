package eod

import (
	"sort"
	"strconv"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

var sortChoices = []*discordgo.ApplicationCommandOptionChoice{
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
}

var sorts = map[string]func(a, b string, dat types.ServerData) bool{
	"length": func(a, b string, dat types.ServerData) bool {
		return len(a) < len(b)
	},
	"name": func(a, b string, dat types.ServerData) bool {
		return compareStrings(a, b)
	},
	"createdon": func(a, b string, dat types.ServerData) bool {
		el1, res := dat.GetElement(a, true)
		el2, res2 := dat.GetElement(b, true)
		if !res.Exists || !res2.Exists {
			return false
		}
		return el1.CreatedOn.Before(el2.CreatedOn)
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

// Less
func compareStrings(a, b string) bool {
	fl1, err := strconv.ParseFloat(a, 32)
	fl2, err2 := strconv.ParseFloat(b, 32)
	if err == nil && err2 == nil {
		return fl1 < fl2
	}
	return a < b
}

func sortStrings(arr []string) {
	sort.Slice(arr, func(i, j int) bool {
		return compareStrings(arr[i], arr[j])
	})
}

func sortElemList(elems []string, sortName string, dat types.ServerData) {
	dat.Lock.RLock()
	sort.Slice(elems, func(i, j int) bool {
		return sorts[sortName](elems[i], elems[j], dat)
	})
	dat.Lock.RUnlock()
}
