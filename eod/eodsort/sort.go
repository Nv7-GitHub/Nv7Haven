package eodsort

import (
	"fmt"
	"sort"
	"strconv"
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
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
	{
		Name:  "Tree Size",
		Value: "treesize",
	},
	{
		Name:  "Color",
		Value: "color",
	},
	{
		Name:  "Found",
		Value: "found",
	},
}

var sorts = map[string]func(a, b int, db *eodb.DB, data any) bool{
	"length": func(a, b int, db *eodb.DB, data any) bool {
		name1, res := db.GetElement(a, true)
		name2, res2 := db.GetElement(b, true)
		if !res.Exists || !res2.Exists {
			return false
		}
		return len(name1.Name) < len(name2.Name)
	},
	"name": func(a, b int, db *eodb.DB, data any) bool {
		name1, res := db.GetElement(a, true)
		name2, res2 := db.GetElement(b, true)
		if !res.Exists || !res2.Exists {
			return false
		}
		return CompareStrings(name1.Name, name2.Name)
	},
	"createdon": func(a, b int, db *eodb.DB, data any) bool {
		el1, res := db.GetElement(a, true)
		el2, res2 := db.GetElement(b, true)
		if !res.Exists || !res2.Exists {
			return false
		}
		return el1.CreatedOn.Time.Before(el2.CreatedOn.Time)
	},
	"id": func(a, b int, db *eodb.DB, data any) bool {
		el1, res := db.GetElement(a, true)
		el2, res2 := db.GetElement(b, true)
		if !res.Exists || !res2.Exists {
			return false
		}
		return el1.ID < el2.ID
	},
	"complexity": func(a, b int, db *eodb.DB, data any) bool {
		el1, res := db.GetElement(a, true)
		el2, res2 := db.GetElement(b, true)
		if !res.Exists || !res2.Exists {
			return false
		}
		return el1.Complexity < el2.Complexity
	},
	"difficulty": func(a, b int, db *eodb.DB, data any) bool {
		el1, res := db.GetElement(a, true)
		el2, res2 := db.GetElement(b, true)
		if !res.Exists || !res2.Exists {
			return false
		}
		return el1.Difficulty < el2.Difficulty
	},
	"usedin": func(a, b int, db *eodb.DB, data any) bool {
		el1, res := db.GetElement(a, true)
		el2, res2 := db.GetElement(b, true)
		if !res.Exists || !res2.Exists {
			return false
		}
		return el1.UsedIn < el2.UsedIn
	},
	"creator": func(a, b int, db *eodb.DB, data any) bool {
		el1, res := db.GetElement(a, true)
		el2, res2 := db.GetElement(b, true)
		if !res.Exists || !res2.Exists {
			return false
		}
		return el1.Creator < el2.Creator
	},
	"treesize": func(a, b int, db *eodb.DB, data any) bool {
		el1, res := db.GetElement(a, true)
		el2, res2 := db.GetElement(b, true)
		if !res.Exists || !res2.Exists {
			return false
		}
		return el1.TreeSize < el2.TreeSize
	},
	"color": func(a, b int, db *eodb.DB, data any) bool {
		el1, res := db.GetElement(a, true)
		el2, res2 := db.GetElement(b, true)
		if !res.Exists || !res2.Exists {
			return false
		}
		return el1.Color < el2.Color
	},
	"found": func(a, b int, db *eodb.DB, data any) bool {
		inv := data.(*types.Inventory)
		v1 := 1
		v2 := 1
		if inv.Contains(a, true) {
			v1--
		}
		if inv.Contains(b, true) {
			v2--
		}
		return v1 < v2
	},
}

var getters = map[string]func(el types.Element, data any) string{
	"createdon": func(el types.Element, data any) string {
		return fmt.Sprintf(" - <t:%d>", el.CreatedOn.Unix())
	},
	"id": func(el types.Element, data any) string {
		return fmt.Sprintf(" - #%d", el.ID)
	},
	"complexity": func(el types.Element, data any) string {
		return fmt.Sprintf(" - %d", el.Complexity)
	},
	"difficulty": func(el types.Element, data any) string {
		return fmt.Sprintf(" - %d", el.Difficulty)
	},
	"usedin": func(el types.Element, data any) string {
		return fmt.Sprintf(" - %d", el.UsedIn)
	},
	"creator": func(el types.Element, data any) string {
		return fmt.Sprintf(" - <@%s>", el.Creator)
	},
	"treesize": func(el types.Element, data any) string {
		return fmt.Sprintf(" - %d", el.TreeSize)
	},
	"color": func(el types.Element, data any) string {
		col, err := util.GetEmoji(el.Color)
		if err == nil {
			return fmt.Sprintf(" - %s", col)
		}
		return ""
	},
	"found": func(el types.Element, data any) string {
		inv := data.(*types.Inventory)
		if inv.Contains(el.ID, true) {
			return " " + types.Check
		}
		return " " + types.X
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

func Sort(vals any, length int, elemGet func(index int) int, elemTxt func(index int) string, elemSet func(index int, val string), sortName string, user string, db *eodb.DB, postfix bool) {
	lock.RLock()
	sorter := sorts[sortName]
	lock.RUnlock()

	var data any
	switch sortName {
	case "found":
		data = db.GetInv(user)
	}

	db.RLock()
	// Lock data
	switch sortName {
	case "found":
		data.(*types.Inventory).Lock.RLock()
	}

	// Sort
	sort.Slice(vals, func(i, j int) bool {
		return sorter(elemGet(i), elemGet(j), db, data)
	})

	// Unlock data
	switch sortName {
	case "found":
		data.(*types.Inventory).Lock.RUnlock()
	}

	if postfix {
		lock.RLock()
		getter, exists := getters[sortName]
		lock.RUnlock()
		if exists {
			// Lock data
			switch sortName {
			case "found":
				data.(*types.Inventory).Lock.RLock()
			}

			for i := 0; i < length; i++ {
				el, res := db.GetElement(elemGet(i), true)
				if res.Exists {
					elemSet(i, elemTxt(i)+getter(el, data))
				}
			}

			// Unlock data
			switch sortName {
			case "found":
				data.(*types.Inventory).Lock.RLock()
			}
		}
	}
	db.RUnlock()
}
