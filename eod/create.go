package eod

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
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
		return
	}
	if count != 0 {
		return
	}

	_, exists = dat.elemCache[strings.ToLower(name)]
	text := "Combination"

	var postTxt string
	if !exists {
		diff := -1
		compl := -1
		areUnique := false
		for _, val := range parents {
			elem := dat.elemCache[strings.ToLower(val)]
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
		elem := element{
			ID:         len(dat.elemCache) + 1,
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

		dat.elemCache[strings.ToLower(elem.Name)] = elem
		dat.invCache[creator][strings.ToLower(elem.Name)] = empty{}
		lock.Lock()
		b.dat[guild] = dat
		lock.Unlock()

		_, err = tx.Exec("INSERT INTO eod_elements VALUES ( ?,  ?, ?, ?, ?, ?, ?, ?, ?, ? )", elem.Name, elem.Image, elem.Guild, elem.Comment, elem.Creator, int(elem.CreatedOn.Unix()), elems2txt(parents), elem.Complexity, elem.Difficulty, 0)
		if err != nil {
			fmt.Println(err)
			return
		}
		text = "Element"

		b.saveInv(guild, creator, true)
	} else {
		el, exists := dat.elemCache[strings.ToLower(name)]
		if !exists {
			fmt.Println("Doesn't exist")
			return
		}
		name = el.Name

		var id int
		res := tx.QueryRow("SELECT COUNT(1) FROM eod_combos WHERE guild=?", guild)
		err = res.Scan(&id)
		if err == nil {
			postTxt = " - Combination **#" + strconv.Itoa(id) + "**"
		}

		dat.invCache[creator][strings.ToLower(name)] = empty{}
		lock.Lock()
		b.dat[guild] = dat
		lock.Unlock()
		b.saveInv(guild, creator, false)
	}
	_, err = tx.Exec("INSERT INTO eod_combos VALUES ( ?, ?, ? )", guild, data, name)
	if err != nil {
		fmt.Println(err)
		return
	}

	params := make(map[string]empty)
	for _, val := range parents {
		params[val] = empty{}
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
			fmt.Println(err)
			return
		}

		el := dat.elemCache[strings.ToLower(k)]
		el.UsedIn++
		dat.elemCache[strings.ToLower(k)] = el
	}
	lock.Lock()
	b.dat[guild] = dat
	lock.Unlock()

	txt := newText + " " + text + " - **" + name + "** (By <@" + creator + ">)" + postTxt

	b.dg.ChannelMessageSend(dat.newsChannel, txt)
	datafile.Write([]byte(fmt.Sprintf("%s %s\n", name, parents)))

	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		return
	}
}
