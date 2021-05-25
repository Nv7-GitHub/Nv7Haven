package eod

import (
	"fmt"
	"math/rand"
	"strings"
)

func (b *EoD) ideaCmd(count int, m msg, rsp rsp) {
	if count > maxComboLength {
		rsp.ErrorMessage(fmt.Sprintf("You can only combine up to %d elements!", maxComboLength))
		return
	}

	if count < 2 {
		rsp.ErrorMessage("You must combine at least 2 elements!")
		return
	}

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

	var elem3 string
	cont := true
	var elems []string
	tries := 0
	for cont {
		elems = make([]string, count)
		for i := range elems {
			cnt := rand.Intn(len(inv))
			j := 0
			for k := range inv {
				if j == cnt {
					elems[i] = k
					break
				}
				j++
			}
		}

		query := "SELECT elem3 FROM eod_combos WHERE elems LIKE ? AND guild=?"
		els := elems2txt(elems)
		if isASCII(els) {
			query = "SELECT elem3 FROM eod_combos WHERE CONVERT(elems USING utf8mb4) LIKE CONVERT(? USING utf8mb4) AND guild=CONVERT(? USING utf8mb4) COLLATE utf8mb4_general_ci"
		}
		row := b.db.QueryRow(query, els, m.GuildID)
		err := row.Scan(&elem3)
		if err != nil {
			cont = false
		}
		tries++

		if tries > 10 {
			rsp.ErrorMessage("Couldn't find a random unused combination, maybe try again later?")
			return
		}
	}

	text := ""
	for i, el := range elems {
		text += dat.elemCache[strings.ToLower(el)].Name
		if i != len(elems)-1 {
			text += " + "
		}
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

	rsp.Resp(fmt.Sprintf("Your random unused combination is... **%s**\n 	Suggest it by typing **/suggest**", text))
}
