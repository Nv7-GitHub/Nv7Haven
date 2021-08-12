package eod

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

const blueCircle = "🔵"

func elems2txt(elems []string) string {
	for i, elem := range elems {
		elems[i] = strings.ToLower(elem)
	}
	sort.Strings(elems)
	return strings.Join(elems, "+")
}

func (b *EoD) combine(elems []string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	dat.Lock.RLock()
	inv, exists := dat.InvCache[m.Author.ID]
	dat.Lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}

	validElems := make([]string, len(elems))
	validCnt := 0
	for _, elem := range elems {
		if len(elem) > 0 {
			validElems[validCnt] = elem
			validCnt++
		}
	}
	elems = validElems[:validCnt]
	if len(elems) == 1 {
		rsp.ErrorMessage("You must combine at least 2 elements!")
		return
	}

	for _, elem := range elems {
		dat.Lock.RLock()
		_, exists := dat.ElemCache[strings.ToLower(elem)]
		dat.Lock.RUnlock()
		if !exists {
			notExists := make(map[string]types.Empty)
			for _, el := range elems {
				dat.Lock.RLock()
				_, exists := dat.ElemCache[strings.ToLower(el)]
				if !exists {
					notExists["**"+el+"**"] = types.Empty{}
				}
				dat.Lock.RUnlock()
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
			dat.Lock.RLock()
			_, exists := dat.CombCache[m.Author.ID]
			if exists {
				dat.Lock.RLock()
				delete(dat.CombCache, m.Author.ID)
				dat.Lock.RUnlock()

				lock.Lock()
				b.dat[m.GuildID] = dat
				lock.Unlock()
			}
			dat.Lock.RUnlock()

			notFound := make(map[string]types.Empty)
			for _, el := range elems {
				_, exists := inv[strings.ToLower(el)]
				if !exists {
					dat.Lock.RLock()
					notFound["**"+dat.ElemCache[strings.ToLower(el)].Name+"**"] = types.Empty{}
					dat.Lock.RUnlock()
				}
			}

			if len(notFound) == 1 {
				dat.Lock.RLock()
				rsp.ErrorMessage(fmt.Sprintf("You don't have **%s**!", dat.ElemCache[strings.ToLower(elem)].Name))
				dat.Lock.RUnlock()
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
		if dat.CombCache == nil {
			dat.CombCache = make(map[string]types.Comb)
		}

		dat.Lock.Lock()
		dat.CombCache[m.Author.ID] = types.Comb{
			Elems: elems,
			Elem3: elem3,
		}
		dat.Lock.Unlock()

		dat.Lock.RLock()
		_, exists := dat.InvCache[m.Author.ID][strings.ToLower(elem3)]
		dat.Lock.RUnlock()
		if !exists {
			dat.Lock.Lock()
			dat.InvCache[m.Author.ID][strings.ToLower(elem3)] = types.Empty{}
			dat.Lock.Unlock()
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

	if dat.CombCache == nil {
		dat.CombCache = make(map[string]types.Comb)
	}

	dat.Lock.Lock()
	dat.CombCache[m.Author.ID] = types.Comb{
		Elems: elems,
		Elem3: "",
	}
	dat.Lock.Unlock()

	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()
	rsp.Resp("That combination doesn't exist! " + redCircle + "\n 	Suggest it by typing **/suggest**")
}

func joinTxt(elemDat map[string]types.Empty, ending string) string {
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
