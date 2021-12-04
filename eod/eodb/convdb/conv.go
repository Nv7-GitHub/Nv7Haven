package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func convDB() {
	home, err := os.UserHomeDir()
	handle(err)
	dbPath := filepath.Join(home, "go/src/github.com/Nv7-Github/Nv7Haven/data/eod")

	err = os.RemoveAll(dbPath)
	handle(err)
	err = os.MkdirAll(dbPath, os.ModePerm)
	handle(err)

	for _, gld := range glds {
		db, err := eodb.NewDB(gld.ID, filepath.Join(dbPath, gld.ID))
		handle(err)

		// Conv elements
		names := make(map[string]int)
		for _, elem := range gld.Elements {
			names[strings.ToLower(elem.Name)] = elem.ID
		}
		for _, elem := range gld.Elements {
			// Conv parents
			cont := false
			pars := make([]int, len(elem.Parents))
			for i, par := range elem.Parents {
				id, exists := names[par]
				if !exists {
					cont = true
					break
				}
				pars[i] = id
			}
			if cont {
				continue
			}

			el := types.Element{
				ID:         elem.ID,
				Name:       elem.Name,
				Image:      elem.Image,
				Color:      elem.Color,
				Guild:      elem.Guild,
				Comment:    elem.Comment,
				Creator:    elem.Creator,
				CreatedOn:  elem.CreatedOn,
				Parents:    pars,
				Complexity: elem.Complexity,
				Difficulty: elem.Difficulty,
				UsedIn:     elem.UsedIn,
				TreeSize:   elem.TreeSize,
			}
			err = db.SaveElement(el, true)
			handle(err)
		}

		// Conv combos
		for _, comb := range gld.Combos {
			vals := strings.Split(comb.Elems, "+")
			elems := make([]int, len(vals))
			cont := false
			for i, v := range vals {
				id, res := db.GetIDByName(v)
				if !res.Exists {
					fmt.Println(v)
					id, exists := names[v]
					if !exists {
						cont = true
						break
					}
					elems[i] = id
					continue
				}
				elems[i] = id
			}
			if cont {
				continue
			}

			out, res := db.GetIDByName(comb.Res)
			if !res.Exists {
				fmt.Println(comb.Res)
				id, exists := names[comb.Res]
				if !exists {
					continue
				}
				out = id
			}

			err = db.AddCombo(elems, out)
			handle(err)
		}

		db.Close()
	}
}
