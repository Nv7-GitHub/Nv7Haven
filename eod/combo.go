package eod

import (
	"fmt"
	"sort"
	"strings"
)

const blueCircle = "🔵"

func elems2txt(elems []string) string {
	for i, elem := range elems {
		elems[i] = strings.ToLower(elem)
	}
	sort.Strings(elems)
	return strings.Join(elems, "+")
}

func (b *EoD) combine(elems []string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	inv, exists := dat.invCache[m.Author.ID]
	if !exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}

	for _, elem := range elems {
		_, exists := dat.elemCache[strings.ToLower(elem)]
		if !exists {
			rsp.ErrorMessage(fmt.Sprintf("Element %s doesn't exist!", elem))
			return
		}

		_, hasElement := inv[strings.ToLower(elem)]
		if !hasElement {
			rsp.ErrorMessage(fmt.Sprintf("You don't have **%s**!", dat.elemCache[strings.ToLower(elem)].Name))
			_, exists := dat.combCache[m.Author.ID]
			if exists {
				delete(dat.combCache, m.Author.ID)
				lock.Lock()
				b.dat[m.GuildID] = dat
				lock.Unlock()
			}
			return
		}
	}
	var elem3 string
	cont := true
	row := b.db.QueryRow("SELECT elem3 FROM eod_combos WHERE elems=? AND guild=?", elems2txt(elems), m.GuildID)
	err := row.Scan(&elem3)
	if err != nil {
		cont = false
	}
	if cont {
		if dat.combCache == nil {
			dat.combCache = make(map[string]comb)
		}

		dat.combCache[m.Author.ID] = comb{
			elems: elems,
			elem3: elem3,
		}
		_, exists := dat.invCache[m.Author.ID][strings.ToLower(elem3)]
		if !exists {
			dat.invCache[m.Author.ID][strings.ToLower(elem3)] = empty{}
			b.saveInv(m.GuildID, m.Author.ID, false)

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
		elems: elems,
		elem3: "",
	}
	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()
	rsp.Resp("That combination doesn't exist! " + redCircle + "\n 	Suggest it by typing **/suggest**")
}
