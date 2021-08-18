package eod

import (
	"context"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

const newText = "🆕"

var datafile *os.File
var createLock = &sync.Mutex{}

func (b *EoD) elemCreate(name string, parents []string, creator string, guild string) {
	lock.RLock()
	dat, exists := b.dat[guild]
	lock.RUnlock()
	if !exists {
		return
	}

	data := elems2txt(parents)
	_, res := dat.GetCombo(data)
	if res.Exists {
		return
	}

	_, res = dat.GetElement(name)
	text := "Combination"

	createLock.Lock()
	tx, err := b.db.BeginTx(context.Background(), nil)
	if err != nil {
		tx.Rollback()
		return
	}

	var postTxt string
	if !res.Exists {
		diff := -1
		compl := -1
		areUnique := false
		for _, val := range parents {
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
		}
		compl++
		if areUnique {
			diff++
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
		}
		postTxt = " - Element **#" + strconv.Itoa(elem.ID) + "**"

		dat.SetElement(elem)

		_, err = tx.Exec("INSERT INTO eod_elements VALUES ( ?,  ?, ?, ?, ?, ?, ?, ?, ?, ? )", elem.Name, elem.Image, elem.Guild, elem.Comment, elem.Creator, int(elem.CreatedOn.Unix()), elems2txt(parents), elem.Complexity, elem.Difficulty, 0)
		if err != nil {
			dat.DeleteElement(elem.Name)

			datafile.WriteString(err.Error() + "\n")
			tx.Rollback()
			return
		}
		text = "Element"
	} else {
		el, res := dat.GetElement(name)
		if !res.Exists {
			datafile.WriteString("Doesn't exist\n")
			tx.Rollback()
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

		datafile.WriteString(err.Error() + "\n")
		tx.Rollback()
		return
	}
	dat.AddComb(data, name)

	params := make(map[string]types.Empty)
	for _, val := range parents {
		params[val] = types.Empty{}
	}
	for k := range params {
		query := "UPDATE eod_elements SET usedin=usedin+1 WHERE name LIKE ? AND guild LIKE ?"
		if isASCII(k) {
			query = "UPDATE eod_elements SET usedin=usedin+1 WHERE CONVERT(name USING utf8mb4) LIKE CONVERT(? USING utf8mb4) AND CONVERT(guild USING utf8mb4) LIKE CONVERT(? USING utf8mb4)"
		}
		if isWildcard(k) {
			query = strings.ReplaceAll(query, " LIKE ", "=")
		}
		_, err = tx.Exec(query, k, guild)
		if err != nil {
			dat.DeleteElement(name)

			datafile.WriteString(err.Error() + "\n")
			tx.Rollback()
			return
		}

		el, res := dat.GetElement(k)
		if res.Exists {
			el.UsedIn++
			dat.SetElement(el)
		}
	}

	txt := newText + " " + text + " - **" + name + "** (By <@" + creator + ">)" + postTxt

	b.dg.ChannelMessageSend(dat.NewsChannel, txt)

	err = tx.Commit()
	if err != nil {
		dat.DeleteElement(name)

		datafile.WriteString(err.Error() + "\n")
		return
	}

	createLock.Unlock()

	// Add Element to Inv
	inv, _ := dat.GetInv(creator, true)
	inv.Add(name)
	dat.SetInv(creator, inv)

	lock.Lock()
	b.dat[guild] = dat
	lock.Unlock()
	b.saveInv(guild, creator, true)

	err = b.autocategorize(name, guild)
	if err != nil {
		datafile.WriteString(err.Error() + "\n")
		return
	}
}
