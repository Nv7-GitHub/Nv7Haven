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

	for i, val := range elems {
		if len(val) == 0 {
			buff := elems[len(elems)-1]
			elems[len(elems)-1] = val
			elems[i] = buff
			elems = elems[:len(elems)-1]
		}
	}

	for _, elem := range elems {
		_, exists := dat.elemCache[strings.ToLower(elem)]
		if !exists {
			notExists := make(map[string]empty)
			for _, el := range elems {
				_, exists := dat.elemCache[strings.ToLower(el)]
				if !exists {
					notExists["**"+el+"**"] = empty{}
				}
			}
			if len(notExists) == 1 {
				rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", elem))
				return
			}

			rsp.ErrorMessage("Elements " + joinTxt(notExists, "and") + " don't exist!")
			return
		}

		_, hasElement := inv[strings.ToLower(elem)]
		if !hasElement {
			_, exists := dat.combCache[m.Author.ID]
			if exists {
				delete(dat.combCache, m.Author.ID)
				lock.Lock()
				b.dat[m.GuildID] = dat
				lock.Unlock()
			}

			notFound := make(map[string]empty)
			for _, el := range elems {
				_, exists := inv[strings.ToLower(el)]
				if !exists {
					notFound["**"+dat.elemCache[strings.ToLower(el)].Name+"**"] = empty{}
				}
			}

			if len(notFound) == 1 {
				rsp.ErrorMessage(fmt.Sprintf("You don't have **%s**!", dat.elemCache[strings.ToLower(elem)].Name))
				return
			}

			rsp.ErrorMessage("You don't have " + joinTxt(notFound, "or") + "!")
			return
		}
	}

	var elem3 string
	cont := true
	query := "SELECT elem3 FROM eod_combos WHERE elems LIKE ? AND guild=?"
	els := elems2txt(elems)
	if isASCII(els) {
		query = "SELECT elem3 FROM eod_combos WHERE CONVERT(elems USING utf8mb4) LIKE CONVERT(? USING utf8mb4) AND guild=CONVERT(? USING utf8mb4) COLLATE utf8mb4_general_ci"
	}

	if isWildcard(els) {
		query = strings.ReplaceAll(query, " LIKE ", "=")
	}

	row := b.db.QueryRow(query, els, m.GuildID)
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

func joinTxt(elemDat map[string]empty, ending string) string {
	elems := make([]string, len(elemDat))
	i := 0
	for k := range elemDat {
		elems[i] = k
		i++
	}
	sort.Strings(elems)

	out := ""
	for i, elem := range elems {
		out += elem
		if i != len(elems)-1 && len(elems) != 2 {
			out += ", "
		} else if i != len(elems)-1 {
			out += " "
		}

		if i == len(elems)-2 {
			out += ending + " "
		}
	}

	return out
}
