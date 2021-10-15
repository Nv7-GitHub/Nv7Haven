package polls

import (
	"context"
	"errors"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/logs"
	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

var createLock = &sync.Mutex{}

func (b *Polls) elemCreate(name string, parents []string, creator string, controversial string, guild string) {
	b.lock.RLock()
	dat, exists := b.dat[guild]
	b.lock.RUnlock()
	if !exists {
		return
	}

	data := util.Elems2Txt(parents)
	_, res := dat.GetCombo(data)
	if res.Exists {
		return
	}

	_, res = dat.GetElement(name)
	text := "Combination"

	createLock.Lock()
	tx, err := b.db.GetSqlDB().BeginTx(context.Background(), nil)
	if err != nil {
		_ = tx.Rollback()
		createLock.Unlock()
		return
	}

	handle := func(err error) {
		log.SetOutput(logs.DataFile)
		log.Println(err)
		tx.Rollback()
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
			elem, _ := dat.GetElement(val)
			if elem.Difficulty > diff {
				diff = elem.Difficulty
			}
			if elem.Complexity > compl {
				compl = elem.Complexity
			}
			if !strings.EqualFold(parents[0], val) {
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
		size, suc, msg := trees.ElemCreateSize(parents, dat)
		if !suc {
			handle(errors.New(msg))
			return
		}
		elem := types.Element{
			ID:         len(dat.Elements) + 1,
			Name:       name,
			Guild:      guild,
			Comment:    "None",
			Creator:    creator,
			CreatedOn:  time.Now(),
			Parents:    parents,
			Complexity: compl,
			Difficulty: diff,
			Color:      col,
			TreeSize:   size,
		}
		postTxt = " - Element **#" + strconv.Itoa(elem.ID) + "**"
		dat.SetElement(elem)

		_, err = tx.Exec("INSERT INTO eod_elements VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )", elem.Name, elem.Image, elem.Color, elem.Guild, elem.Comment, elem.Creator, int(elem.CreatedOn.Unix()), util.Elems2Txt(parents), elem.Complexity, elem.Difficulty, 0, elem.TreeSize)
		if err != nil {
			dat.DeleteElement(elem.Name)

			handle(err)
			return
		}
		text = "Element"
	} else {
		el, res := dat.GetElement(name)
		if !res.Exists {
			log.SetOutput(logs.DataFile)
			log.Println("Doesn't exist")

			_ = tx.Rollback()
			createLock.Unlock()
			return
		}
		name = el.Name

		id := len(dat.Combos)
		if err == nil {
			postTxt = " - Combination **#" + strconv.Itoa(id) + "**"
		}
	}

	_, err = tx.Exec("INSERT INTO eod_combos VALUES ( ?, ?, ? )", guild, data, name)
	if err != nil {
		dat.DeleteElement(name)
		createLock.Unlock()

		log.SetOutput(logs.DataFile)
		log.Println(err)
		_ = tx.Rollback()
		return
	}
	dat.AddComb(data, name)

	params := make(map[string]types.Empty)
	for _, val := range parents {
		params[val] = types.Empty{}
	}
	for k := range params {
		query := "UPDATE eod_elements SET usedin=usedin+1 WHERE name LIKE ? AND guild LIKE ?"
		if util.IsASCII(k) {
			query = "UPDATE eod_elements SET usedin=usedin+1 WHERE CONVERT(name USING utf8mb4) LIKE CONVERT(? USING utf8mb4) AND CONVERT(guild USING utf8mb4) LIKE CONVERT(? USING utf8mb4)"
		}
		if util.IsWildcard(k) {
			query = strings.ReplaceAll(query, " LIKE ", "=")
		}
		_, err = tx.Exec(query, k, guild)
		if err != nil {
			dat.DeleteElement(name)
			createLock.Unlock()

			log.SetOutput(logs.DataFile)
			log.Println(err)
			_ = tx.Rollback()
			return
		}

		el, res := dat.GetElement(k)
		if res.Exists {
			el.UsedIn++
			dat.SetElement(el)
		}
	}

	txt := types.NewText + " " + text + " - **" + name + "** (By <@" + creator + ">)" + postTxt + controversial

	_, _ = b.dg.ChannelMessageSend(dat.NewsChannel, txt)

	err = tx.Commit()
	if err != nil {
		dat.DeleteElement(name)
		createLock.Unlock()

		log.SetOutput(logs.DataFile)
		log.Println(err)
		return
	}

	createLock.Unlock()

	// Add Element to Inv
	inv, _ := dat.GetInv(creator, true)
	inv.Add(name)
	dat.SetInv(creator, inv)

	b.lock.Lock()
	b.dat[guild] = dat
	b.lock.Unlock()
	b.base.SaveInv(guild, creator, true)

	err = b.Autocategorize(name, guild)
	if err != nil {
		log.SetOutput(logs.DataFile)
		log.Println(err)
		return
	}
}
