package eod

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

const newText = "🆕"

var datafile *os.File

func (b *EoD) elemCreate(name string, parents []string, creator string, guild string) {
	lock.RLock()
	dat, exists := b.dat[guild]
	lock.RUnlock()
	if !exists {
		return
	}

	tx, err := b.db.BeginTx(context.Background(), nil)
	if err != nil {
		tx.Rollback()
		return
	}

	data := elems2txt(parents)
	query := "SELECT COUNT(1) FROM eod_combos WHERE guild=? AND elems LIKE ?"
	if isWildcard(data) {
		query = strings.ReplaceAll(query, " LIKE ", "=")
	}

	row := b.db.QueryRow(query, guild, data)
	var count int
	err = row.Scan(&count)
	if err != nil {
		tx.Rollback()
		return
	}
	if count != 0 {
		tx.Rollback()
		return
	}

	dat.Lock.RLock()
	_, exists = dat.ElemCache[strings.ToLower(name)]
	text := "Combination"
	dat.Lock.RUnlock()

	var postTxt string
	if !exists {
		diff := -1
		compl := -1
		areUnique := false
		for _, val := range parents {
			dat.Lock.RLock()
			elem := dat.ElemCache[strings.ToLower(val)]
			dat.Lock.RUnlock()
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
		dat.Lock.RLock()
		elem := types.Element{
			ID:         len(dat.ElemCache) + 1,
			Name:       name,
			Guild:      guild,
			Comment:    "None",
			Creator:    creator,
			CreatedOn:  time.Now(),
			Parents:    parents,
			Complexity: compl,
			Difficulty: diff,
		}
		dat.Lock.RUnlock()
		postTxt = " - Element **#" + strconv.Itoa(elem.ID) + "**"

		dat.Lock.Lock()
		dat.ElemCache[strings.ToLower(elem.Name)] = elem
		dat.Lock.Unlock()

		_, err = tx.Exec("INSERT INTO eod_elements VALUES ( ?,  ?, ?, ?, ?, ?, ?, ?, ?, ? )", elem.Name, elem.Image, elem.Guild, elem.Comment, elem.Creator, int(elem.CreatedOn.Unix()), elems2txt(parents), elem.Complexity, elem.Difficulty, 0)
		if err != nil {
			dat.Lock.RLock()
			delete(dat.ElemCache, strings.ToLower(elem.Name))
			dat.Lock.RUnlock()

			fmt.Println(err)
			tx.Rollback()
			return
		}
		text = "Element"
	} else {
		dat.Lock.RLock()
		el, exists := dat.ElemCache[strings.ToLower(name)]
		dat.Lock.RUnlock()
		if !exists {
			fmt.Println("Doesn't exist")
			tx.Rollback()
			return
		}
		name = el.Name

		var id int
		res := tx.QueryRow("SELECT COUNT(1) FROM eod_combos WHERE guild=?", guild)
		err = res.Scan(&id)
		if err == nil {
			postTxt = " - Combination **#" + strconv.Itoa(id) + "**"
		}
	}

	_, err = tx.Exec("INSERT INTO eod_combos VALUES ( ?, ?, ? )", guild, data, name)
	if err != nil {
		dat.Lock.RLock()
		delete(dat.ElemCache, strings.ToLower(name))
		dat.Lock.RUnlock()

		fmt.Println(err)
		tx.Rollback()
		return
	}

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
			dat.Lock.RLock()
			delete(dat.ElemCache, strings.ToLower(name))
			dat.Lock.RUnlock()

			fmt.Println(err)
			tx.Rollback()
			return
		}

		dat.Lock.RLock()
		el := dat.ElemCache[strings.ToLower(k)]
		dat.Lock.RUnlock()
		el.UsedIn++
		dat.Lock.Lock()
		dat.ElemCache[strings.ToLower(k)] = el
		dat.Lock.Unlock()
	}

	txt := newText + " " + text + " - **" + name + "** (By <@" + creator + ">)" + postTxt

	b.dg.ChannelMessageSend(dat.NewsChannel, txt)
	datafile.Write([]byte(fmt.Sprintf("%s %s\n", name, parents)))

	// Add Element to Inv
	dat.Lock.Lock()
	dat.InvCache[creator][strings.ToLower(name)] = types.Empty{}
	dat.Lock.Unlock()

	lock.Lock()
	b.dat[guild] = dat
	lock.Unlock()
	b.saveInv(guild, creator, true)

	err = tx.Commit()
	if err != nil {
		dat.Lock.RLock()
		delete(dat.ElemCache, strings.ToLower(name))
		dat.Lock.RUnlock()

		fmt.Println(err)
		return
	}

	err = b.autocategorize(name, guild)
	if err != nil {
		fmt.Println(err)
		return
	}
}
