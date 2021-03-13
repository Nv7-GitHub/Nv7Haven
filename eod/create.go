package eod

import (
	"encoding/json"
	"strings"
	"time"
)

const newText = "🆕"

func (b *EoD) elemCreate(name string, parent1 string, parent2 string, creator string, guild string) {
	lock.RLock()
	dat, exists := b.dat[guild]
	lock.RUnlock()
	if !exists {
		return
	}
	row := b.db.QueryRow("SELECT COUNT(1) FROM eod_elements WHERE name=?", name)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return
	}
	text := "Combination"
	if count == 0 {
		elem := element{
			Name:       name,
			Categories: make(map[string]empty),
			Guild:      guild,
			Comment:    "None",
			Creator:    creator,
			CreatedOn:  time.Now(),
			Parents:    []string{parent1, parent2},
			Complexity: max(dat.elemCache[strings.ToLower(parent1)].Complexity, dat.elemCache[strings.ToLower(parent2)].Complexity) + 1,
		}
		dat.elemCache[strings.ToLower(elem.Name)] = elem
		dat.invCache[creator][strings.ToLower(elem.Name)] = empty{}
		lock.Lock()
		b.dat[guild] = dat
		lock.Unlock()
		cats, err := json.Marshal(elem.Categories)
		if err != nil {
			return
		}
		b.db.Exec("INSERT INTO eod_elements VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )", elem.Name, string(cats), elem.Image, elem.Guild, elem.Comment, elem.Creator, int(elem.CreatedOn.Unix()), elem.Parents[0], elem.Parents[1], elem.Complexity)
		if err != nil {
			return
		}
		text = "Element"

		b.saveInv(guild, creator)
	} else {
		row := b.db.QueryRow("SELECT name FROM eod_elements WHERE name=?", name)
		err = row.Scan(&name)
		if err != nil {
			return
		}
	}
	b.db.Exec("INSERT INTO eod_combos VALUES ( ?, ?, ?, ? )", guild, parent1, parent2, name)
	b.dg.ChannelMessageSend(dat.newsChannel, newText+" "+text+" - **"+name+"** (By <@"+creator+">)")
}

func max(val1, val2 int) int {
	if val1 > val2 {
		return val1
	}
	return val2
}
