package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/schollz/progressbar/v3"
)

type Combo struct {
	Elems string
	Res   string
}

type Guild struct {
	ID       string
	Combos   []Combo
	Elements []types.OldElement
	Cats     map[string]types.OldCategory
	Invs     map[string]map[string]types.Empty
	Config   *types.ServerConfig
}

func NewGuild(id string) *Guild {
	return &Guild{
		ID:       id,
		Combos:   make([]Combo, 0),
		Elements: make([]types.OldElement, 0),
		Cats:     make(map[string]types.OldCategory),
		Invs:     make(map[string]map[string]types.Empty),
		Config:   types.NewServerConfig(),
	}
}

var glds = make(map[string]*Guild)

func loadDB(refresh bool) {
	// No refresh?
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

	// Config
	res, err := db.Query("SELECT * FROM eod_serverdata WHERE 1")
	handle(err)
	defer res.Close()

	var guild string
	var kind types.ServerDataType
	var value1 string
	var intval int
	for res.Next() {
		err = res.Scan(&guild, &kind, &value1, &intval)
		handle(err)

		_, exists := glds[guild]
		if !exists {
			glds[guild] = NewGuild(guild)
		}
		switch kind {
		case types.NewsChannel:
			glds[guild].Config.NewsChannel = value1

		case types.PlayChannel:
			glds[guild].Config.PlayChannels[value1] = types.Empty{}

		case types.VotingChannel:
			glds[guild].Config.VotingChannel = value1

		case types.VoteCount:
			glds[guild].Config.VoteCount = intval

		case types.PollCount:
			glds[guild].Config.PollCount = intval

		case types.ModRole:
			glds[guild].Config.ModRole = value1

		case types.UserColor:
			glds[guild].Config.UserColors[value1] = intval

		}
	}

	// Elements
	fmt.Println("Elements...")
	bar := progressbar.New(cnt)
	elems, err := db.Query("SELECT name, image, color, guild, comment, creator, createdon, parents, complexity, difficulty, usedin, treesize FROM `eod_elements` ORDER BY createdon ASC")
	handle(err)
	defer elems.Close()
	elem := types.OldElement{}
	var createdon int64
	var parentDat string
	for elems.Next() {
		err = elems.Scan(&elem.Name, &elem.Image, &elem.Color, &elem.Guild, &elem.Comment, &elem.Creator, &createdon, &parentDat, &elem.Complexity, &elem.Difficulty, &elem.UsedIn, &elem.TreeSize)
		if err != nil {
			return
		}
		elem.CreatedOn = time.Unix(createdon, 0)

		if len(parentDat) == 0 {
			elem.Parents = make([]string, 0)
		} else {
			elem.Parents = strings.Split(parentDat, "+")
		}

		if elem.Guild == "" {
			continue
		}
		elem.ID = len(glds[elem.Guild].Elements) + 1
		glds[elem.Guild].Elements = append(glds[elem.Guild].Elements, elem)

		bar.Add(1)
	}
	bar.Finish()
	fmt.Println()

	// Combos
	fmt.Println("Combos...")
	err = db.QueryRow("SELECT COUNT(1) FROM eod_combos").Scan(&cnt)
	if err != nil {
		panic(err)
	}
	bar = progressbar.New(cnt)

	combs, err := db.Query("SELECT * FROM `eod_combos`")
	if err != nil {
		panic(err)
	}
	defer combs.Close()
	var elemsVal string
	var elem3 string
	num := 0
	for combs.Next() {
		err = combs.Scan(&guild, &elemsVal, &elem3)
		if err != nil {
			return
		}
		glds[guild].Combos = append(glds[guild].Combos, Combo{
			Elems: elemsVal,
			Res:   elem3,
		})

		bar.Add(1)
		num++
	}
	bar.Finish()
	fmt.Println()

	// Invs
	fmt.Println("Invs...")
	err = db.QueryRow("SELECT COUNT(1) FROM eod_inv").Scan(&cnt)
	handle(err)

	bar = progressbar.New(cnt)

	invs, err := db.Query("SELECT guild, user, inv, made FROM eod_inv WHERE 1")
	handle(err)
	defer invs.Close()
	var invDat string
	var user string
	var inv map[string]types.Empty
	var madecnt int
	for invs.Next() {
		inv = make(map[string]types.Empty)
		err = invs.Scan(&guild, &user, &invDat, &madecnt)
		handle(err)
		err = json.Unmarshal([]byte(invDat), &inv)
		handle(err)
		glds[guild].Invs[user] = inv

		bar.Add(1)
	}
	bar.Finish()
	fmt.Println()

	// Cats
	err = db.QueryRow("SELECT COUNT(1) FROM eod_categories").Scan(&cnt)
	handle(err)
	fmt.Println("Cats...")
	bar = progressbar.New(cnt)

	cats, err := db.Query("SELECT * FROM eod_categories")
	handle(err)
	defer cats.Close()
	var elemDat string
	cat := types.OldCategory{}
	for cats.Next() {
		err = cats.Scan(&guild, &cat.Name, &elemDat, &cat.Image, &cat.Color)
		if err != nil {
			return
		}

		cat.Guild = guild

		cat.Elements = make(map[string]types.Empty)
		err := json.Unmarshal([]byte(elemDat), &cat.Elements)
		handle(err)

		glds[guild].Cats[strings.ToLower(cat.Name)] = cat

		bar.Add(1)
	}
	bar.Finish()
	fmt.Println()

	// Save?
	/*cache, err := os.Create("data.gob")
	handle(err)
	defer cache.Close()

	enc := gob.NewEncoder(cache)
	err = enc.Encode(glds)
	handle(err)*/
}
