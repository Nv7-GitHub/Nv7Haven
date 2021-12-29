package polls

import (
	"errors"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/logs"
	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

var createLock = &sync.Mutex{}

func (b *Polls) elemCreate(name string, parents []int, creator string, controversial string, guild string) {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return
	}

	_, res = db.GetCombo(parents)
	if res.Exists {
		return
	}

	_, res = db.GetElementByName(name)
	text := "Combination"

	createLock.Lock()

	handle := func(err error) {
		log.SetOutput(logs.DataFile)
		log.Println(err)
		createLock.Unlock()
	}

	var postTxt string
	if !res.Exists {
		// Element doesnt exist
		diff := -1
		compl := -1
		areUnique := false
		parColors := make([]int, len(parents))
		for j, val := range parents {
			elem, _ := db.GetElement(val)
			if elem.Difficulty > diff {
				diff = elem.Difficulty
			}
			if elem.Complexity > compl {
				compl = elem.Complexity
			}
			if parents[0] != val {
				areUnique = true
			}
			parColors[j] = elem.Color
		}
		compl++
		if areUnique {
			diff++
		}
		col, err := util.MixColors(parColors)
		if err != nil {
			handle(err)
			return
		}
		size, suc, msg := trees.ElemCreateSize(parents, db)
		if !suc {
			handle(errors.New(msg))
			return
		}
		elem := types.Element{
			ID:         len(db.Elements) + 1,
			Name:       name,
			Guild:      guild,
			Comment:    "None",
			Creator:    creator,
			CreatedOn:  types.NewTimeStamp(time.Now()),
			Parents:    parents,
			Complexity: compl,
			Difficulty: diff,
			Color:      col,
			TreeSize:   size,
		}
		postTxt = " - Element **#" + strconv.Itoa(elem.ID) + "**"
		err = db.SaveElement(elem, true)
		if err != nil {
			handle(err)
			return
		}

		text = "Element"
	} else {
		el, res := db.GetElementByName(name)
		if !res.Exists {
			log.SetOutput(logs.DataFile)
			log.Println("Doesn't exist")

			createLock.Unlock()
			return
		}
		name = el.Name

		id := db.ComboCnt()
		postTxt = " - Combination **#" + strconv.Itoa(id) + "**"
	}

	el, _ := db.GetElementByName(name)
	err := db.AddCombo(parents, el.ID)
	if err != nil {
		handle(err)
		return
	}

	params := make(map[int]types.Empty)
	for _, val := range parents {
		params[val] = types.Empty{}
	}
	for k := range params {
		el, res := db.GetElement(k)
		if res.Exists {
			el.UsedIn++
			err := db.SaveElement(el)
			if err != nil {
				log.SetOutput(logs.DataFile)
				log.Println(err)
			}
		}
	}

	txt := types.NewText + " " + text + " - **" + name + "** (By <@" + creator + ">)" + postTxt + controversial

	_, _ = b.dg.ChannelMessageSend(db.Config.NewsChannel, txt)

	createLock.Unlock()

	// Add Element to Inv
	inv := db.GetInv(creator)
	inv.Add(el.ID)
	err = db.SaveInv(inv, true, true)
	if err != nil {
		log.SetOutput(logs.DataFile)
		log.Println(err)
	}

	err = b.Autocategorize(name, guild)
	if err != nil {
		log.SetOutput(logs.DataFile)
		log.Println(err)
		return
	}
}
