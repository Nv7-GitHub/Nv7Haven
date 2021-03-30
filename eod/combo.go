package eod

import (
	"encoding/json"
	"fmt"
	"strings"
)

const blueCircle = "🔵"

func (b *EoD) combine(elems []string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	inps := make([]interface{}, len(elems))
	for i, val := range elems {
		_, exists = dat.elemCache[strings.ToLower(val)]
		if !exists {
			rsp.ErrorMessage(fmt.Sprintf("Element %s doesn't exist!", val))
			return
		}
		_, exists = dat.invCache[m.Author.ID][strings.ToLower(val)]
		if !exists {
			rsp.ErrorMessage(fmt.Sprintf("You don't have %s!", val))
			return
		}
		inps[i] = interface{}(strings.ToLower(val))
	}
	inps = append([]interface{}{m.GuildID}, inps...)

	where := "guild=?"
	for i := 0; i < len(elems); i++ {
		where += ` AND (JSON_EXTRACT(elems, CONCAT("$.", ?)) IS NOT NULL)`
	}
	cont := false
	var elem3 string
	var elemDat string
	res, err := b.db.Query("SELECT elems, elem3 FROM eod_combos WHERE "+where+" ORDER BY JSON_LENGTH(elems) ASC", inps...)
	if rsp.Error(err) {
		return
	}
	defer res.Close()
	for res.Next() {
		err = res.Scan(&elemDat, &elem3)
		if rsp.Error(err) {
			return
		}
		var comboDat map[string]empty
		err = json.Unmarshal([]byte(elemDat), &comboDat)
		if rsp.Error(err) {
			return
		}
		if (len(comboDat) != len(elems)) && (len(comboDat) != 1) {
			cont = false
		} else {
			cont = true
			break
		}
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
		elems: elems,
		elem3: "",
	}
	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()
	rsp.Resp("That combination doesn't exist! " + redCircle + "\n 	Suggest it by typing **/suggest**")
}
