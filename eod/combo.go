package eod

import (
	"fmt"
	"strings"
)

const blueCircle = "🔵"

func (b *EoD) combine(elem1 string, elem2 string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	elem1dat, exists := dat.elemCache[strings.ToLower(elem1)]
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Element %s doesn't exist!", elem1))
		return
	}
	elem2dat, exists := dat.elemCache[strings.ToLower(elem2)]
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Element %s doesn't exist!", elem2))
		return
	}
	_, exists = dat.invCache[m.Author.ID][strings.ToLower(elem1)]
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("You don't have %s!", elem1))
		return
	}
	_, exists = dat.invCache[m.Author.ID][strings.ToLower(elem2)]
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("You don't have %s!", elem2))
		return
	}
	row := b.db.QueryRow("SELECT COUNT(1) FROM eod_combos WHERE guild=? AND ((elem1=? AND elem2=?) OR (elem1=? AND elem2=?))", m.GuildID, elem1, elem2, elem2, elem1)
	var count int
	err := row.Scan(&count)
	if rsp.Error(err) {
		return
	}

	if count > 0 {
		var elem3 string
		row = b.db.QueryRow("SELECT elem3 FROM eod_combos WHERE guild=? AND ((elem1=? AND elem2=?) OR (elem1=? AND elem2=?))", m.GuildID, elem1, elem2, elem2, elem1)
		err = row.Scan(&elem3)
		if rsp.Error(err) {
			return
		}

		if dat.combCache == nil {
			dat.combCache = make(map[string]comb)
		}

		dat.combCache[m.Author.ID] = comb{
			elem1: elem1dat.Name,
			elem2: elem2dat.Name,
			elem3: elem3,
		}
		_, exists := dat.invCache[m.Author.ID][strings.ToLower(elem3)]
		if !exists {
			dat.invCache[m.Author.ID][strings.ToLower(elem3)] = empty{}
			b.saveInv(m.GuildID, m.Author.ID)

			rsp.Resp(fmt.Sprintf("You made **%s** "+newText, elem3))
			return
		}

		rsp.Resp(fmt.Sprintf("You made **%s**, but already have it "+blueCircle, elem3))

		lock.Lock()
		b.dat[m.GuildID] = dat
		lock.Unlock()
		return
	}

	if dat.combCache == nil {
		dat.combCache = make(map[string]comb)
	}

	dat.combCache[m.Author.ID] = comb{
		elem1: elem1dat.Name,
		elem2: elem2dat.Name,
		elem3: "",
	}
	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()
	rsp.Resp("That combination doesn't exist! " + redCircle + "\n 	Suggest it by typing **/suggest**")
}
