package main

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

// ~/go/src/github.com/Nv7-Github/Nv7Haven/data/eod/705084182673621033
func convDB() {
	home, err := os.UserHomeDir()
	handle(err)
	dbPath := filepath.Join(home, "go/src/github.com/Nv7-Github/Nv7haven/data/eod")

	err = os.RemoveAll(dbPath)
	handle(err)
	err = os.MkdirAll(dbPath, os.ModePerm)
	handle(err)

	for _, gld := range glds {
		db, err := eodb.NewDB(gld.ID, filepath.Join(dbPath, gld.ID))
		handle(err)

		db.Config = gld.Config
		err = db.SaveConfig()
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

		getId := func(name string) (int, bool) {
			out, res := db.GetIDByName(name)
			if !res.Exists {
				//fmt.Println(name)
				id, exists := names[name]
				if !exists {
					return 0, false
				}
				out = id
			}
			return out, true
		}

		// Conv combos
		for _, comb := range gld.Combos {
			vals := strings.Split(comb.Elems, "+")
			elems := make([]int, len(vals))
			cont := false
			for i, v := range vals {
				id, suc := getId(v)
				if !suc {
					cont = true
					break
				}
				elems[i] = id
			}
			if cont {
				continue
			}

			out, suc := getId(comb.Res)
			if !suc {
				continue
			}

			err = db.AddCombo(elems, out)
			handle(err)
		}

		// Conv invs
		for id, inv := range gld.Invs {
			i := db.GetInv(id)
			for elem := range inv {
				id, suc := getId(elem)
				if !suc {
					continue
				}
				i.Add(id)
			}
			err = db.SaveInv(i)
			handle(err)
		}

		// Conv cats
		for _, cat := range gld.Cats {
			txt := url.PathEscape(cat.Name)
			if len(txt) > 1024 {
				continue
			}
			c := db.NewCat(cat.Name)
			c.Color = cat.Color
			c.Image = cat.Image
			for elem := range cat.Elements {
				id, ok := getId(strings.ToLower(elem))
				if !ok {
					continue
				}
				c.Elements[id] = types.Empty{}
			}
			err = db.SaveCat(c)
			//handle(err)
		}

		db.Close()
	}
}
