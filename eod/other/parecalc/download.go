package main

import (
	"encoding/gob"
	"os"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

func loadData(refresh bool) {
	if !refresh {
		f, err := os.Open("data.gob")
		handle(err)
		defer f.Close()

		dec := gob.NewDecoder(f)
		err = dec.Decode(&glds)
		handle(err)
		return
	}

	var cnt int
	err := db.QueryRow("SELECT COUNT(1) FROM eod_elements").Scan(&cnt)
	handle(err)

	res, err := db.Query("SELECT * FROM eod_elements")
	handle(err)
	defer res.Close()

	var elem Element
	var parents string
	var createdon int64

	bar := progressbar.New(cnt)
	for res.Next() {
		err = res.Scan(&elem.Name, &elem.Image, &elem.Guild, &elem.Comment, &elem.Creator, &createdon, &parents, &elem.Complexity, &elem.Difficulty, &elem.UsedIn)
		handle(err)
		elem.CreatedOn = time.Unix(createdon, 0)

		if len(parents) == 0 {
			elem.Parents = make([]string, 0)
		} else {
			elem.Parents = strings.Split(parents, "+")
		}

		gld, exists := glds[elem.Guild]
		if !exists {
			gld = NewGuild()
		}
		gld.Elements[strings.ToLower(elem.Name)] = elem
		glds[elem.Guild] = gld
		bar.Add(1)
	}
	bar.Clear()
	bar.Close()

	err = db.QueryRow("SELECT COUNT(1) FROM eod_combos").Scan(&cnt)
	handle(err)

	combos, err := db.Query("SELECT * FROM eod_combos")
	handle(err)
	defer combos.Close()

	var guild string
	var elems string
	var elem3 string
	done := 0

	bar = progressbar.New(cnt)
	for combos.Next() {
		err = combos.Scan(&guild, &elems, &elem3)
		handle(err)

		gld, exists := glds[guild]
		if !exists {
			gld = NewGuild()
		}
		gld.Combos[done] = Combo{
			Elems: strings.Split(elems, "+"),
			Elem3: strings.ToLower(elem3), // Lowercase to speed up rest of program
		}
		glds[guild] = gld

		bar.Add(1)
		done++
	}
	bar.Clear()
	bar.Close()

	cache, err := os.Create("data.gob")
	handle(err)
	defer cache.Close()

	enc := gob.NewEncoder(cache)
	err = enc.Encode(glds)
	handle(err)
}
