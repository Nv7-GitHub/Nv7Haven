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
		el, exists := dat.elemCache[strings.ToLower(name)]
		if !exists {
			return
		}
		if len(el.Parents) == 2 {
			par1, exists := dat.elemCache[strings.ToLower(parent1)]
			if !exists {
				return
			}
			par2, exists := dat.elemCache[strings.ToLower(parent2)]
			if !exists {
				return
			}

			comp := 0
			if par1.Complexity > par2.Complexity {
				comp = par1.Complexity
			} else {
				comp = par2.Complexity
			}
			comp++

			if comp < el.Complexity {
				b.db.Exec("UPDATE eod_elements SET parent1=?,parent2=?,complexity=? WHERE name=? AND guild=?", par1.Name, par2.Name, comp, el.Name, el.Guild)

				el.Complexity = comp
				el.Parents = []string{par1.Name, par2.Name}
				dat.elemCache[strings.ToLower(el.Name)] = el
				lock.Lock()
				b.dat[guild] = dat
				lock.Unlock()
			}
		}
		name = el.Name

		dat.invCache[creator][strings.ToLower(name)] = empty{}
		lock.Lock()
		b.dat[guild] = dat
		lock.Unlock()
		b.saveInv(guild, creator)
	}
	row = b.db.QueryRow("SELECT COUNT(1) FROM eod_combos WHERE guild=? AND (elem1=? AND elem2=?) OR (elem1=? AND elem2=?)", guild, parent1, parent2, parent2, parent1)
	err = row.Scan(&count)
	if err != nil {
		return
	}
	if count == 0 {
		b.db.Exec("INSERT INTO eod_combos VALUES ( ?, ?, ?, ? )", guild, parent1, parent2, name)
	}
	b.dg.ChannelMessageSend(dat.newsChannel, newText+" "+text+" - **"+name+"** (By <@"+creator+">)")
}

func max(val1, val2 int) int {
	if val1 > val2 {
		return val1
	}
	return val2
}
