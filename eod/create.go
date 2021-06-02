package eod

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
		log.Println(err)
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
		log.Println(err)
		return
	}
	if count != 0 {
		return
	}

	row = b.db.QueryRow("SELECT COUNT(1) FROM eod_elements WHERE name=? AND guild=?", name, guild)
	err = row.Scan(&count)
	if err != nil {
		log.Println(err)
		return
	}
	text := "Combination"

	var id int
	if count == 0 {
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
			Categories: make(map[string]empty),
			Guild:      guild,
			Comment:    "None",
			Creator:    creator,
			CreatedOn:  time.Now(),
			Parents:    parents,
			Complexity: compl,
			Difficulty: diff,
		}
		id = elem.ID

		dat.elemCache[strings.ToLower(elem.Name)] = elem
		dat.invCache[creator][strings.ToLower(elem.Name)] = empty{}
		lock.Lock()
		b.dat[guild] = dat
		lock.Unlock()
		cats, err := json.Marshal(elem.Categories)
		if err != nil {
			log.Println(err)
			return
		}

		_, err = tx.Exec("INSERT INTO eod_elements VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )", elem.Name, string(cats), elem.Image, elem.Guild, elem.Comment, elem.Creator, int(elem.CreatedOn.Unix()), elems2txt(parents), elem.Complexity, elem.Difficulty, 0)
		if err != nil {
			log.Println(err)
			return
		}
		text = "Element"

		b.saveInv(guild, creator, true)
	} else {
		el, exists := dat.elemCache[strings.ToLower(name)]
		if !exists {
			return
		}
		name = el.Name
		id = el.ID

		dat.invCache[creator][strings.ToLower(name)] = empty{}
		lock.Lock()
		b.dat[guild] = dat
		lock.Unlock()
		b.saveInv(guild, creator, false)
	}
	_, err = tx.Exec("INSERT INTO eod_combos VALUES ( ?, ?, ? )", guild, data, name)
	if err != nil {
		log.Println(err)
		return
	}

	params := make(map[string]empty)
	for _, val := range parents {
		params[val] = empty{}
	}
	for k := range params {
		_, err = tx.Exec("UPDATE eod_elements SET usedin=usedin+1 WHERE name=? AND guild=?", k, guild)
		if err != nil {
			log.Println(err)
			return
		}

		el := dat.elemCache[strings.ToLower(k)]
		el.UsedIn++
		dat.elemCache[strings.ToLower(k)] = el
	}
	lock.Lock()
	b.dat[guild] = dat
	lock.Unlock()

	txt := newText + " " + text + " - **" + name + "** (By <@" + creator + ">)" + " - Element **#" + strconv.Itoa(id) + "**"

	b.dg.ChannelMessageSend(dat.newsChannel, txt)
	if guild == "819077688371314718" {
		datafile.Write([]byte(fmt.Sprintf("%s %s\n", name, parents)))
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return
	}
}
