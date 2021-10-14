package basecmds

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

const blueCircle = "🔵"

func (b *BaseCmds) Combine(elems []string, m types.Msg, rsp types.Rsp) {
	ok := b.base.CheckServer(m, rsp)
	if !ok {
		return
	}

	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		return
	}

	inv, res := dat.GetInv(m.Author.ID, true)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	// Get rid of nulls
	validElems := make([]string, len(elems))
	validCnt := 0
	for _, elem := range elems {
		elem = strings.ReplaceAll(elem, "\n", "") // Remove newlines
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
	if len(elems) > types.MaxComboLength {
		rsp.ErrorMessage(fmt.Sprintf("You can only combine up to %d elements!", types.MaxComboLength))
		return
	}

	// Check if don't have or not exists
	donthave := false
	elExists := true
	for _, elem := range elems {
		_, res := dat.GetElement(elem)
		if !res.Exists {
			elExists = false
			break
		}

		_, hasElement := inv[strings.ToLower(elem)]
		if !hasElement {
			donthave = true
		}
	}
	if !elExists {
		notExists := make(map[string]types.Empty)
		for _, el := range elems {
			_, res = dat.GetElement(el)
			if !res.Exists {
				notExists["**"+el+"**"] = types.Empty{}
			}
		}
		if len(notExists) == 1 {
			el := ""
			for k := range notExists {
				el = k
				break
			}
			rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", el))
			return
		}

		rsp.ErrorMessage("Elements " + util.JoinTxt(notExists, "and") + " don't exist!")
		return
	}
	if donthave {
		_, res := dat.GetComb(m.Author.ID)
		if res.Exists {
			dat.DeleteComb(m.Author.ID)

			b.lock.Lock()
			b.dat[m.GuildID] = dat
			b.lock.Unlock()
		}

		notFound := make(map[string]types.Empty)
		for _, el := range elems {
			_, exists := inv[strings.ToLower(el)]
			if !exists {
				elem, _ := dat.GetElement(el)
				notFound["**"+elem.Name+"**"] = types.Empty{}
			}
		}

		if len(notFound) == 1 {
			el := ""
			for k := range notFound {
				el = k
				break
			}
			id := rsp.ErrorMessage(fmt.Sprintf("You don't have **%s**!", el))
			dat.SetMsgElem(id, el[2:len(el)-2])
			b.lock.Lock()
			b.dat[m.GuildID] = dat
			b.lock.Unlock()
			return
		}

		rsp.ErrorMessage("You don't have " + util.JoinTxt(notFound, "or") + "!")
		return
	}

	// Combine elements
	elem3, res := dat.GetCombo(util.Elems2Txt(elems))
	if res.Exists {
		dat.SetComb(m.Author.ID, types.Comb{
			Elems: elems,
			Elem3: elem3,
		})

		inv, res := dat.GetInv(m.Author.ID, true)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}

		exists = inv.Contains(elem3)
		if !exists {
			inv.Add(elem3)
			dat.SetInv(m.Author.ID, inv)
			b.base.SaveInv(m.GuildID, m.Author.ID, false)

			id := rsp.Message(fmt.Sprintf("You made **%s** "+types.NewText, elem3))
			dat.SetMsgElem(id, elem3)

			b.lock.Lock()
			b.dat[m.GuildID] = dat
			b.lock.Unlock()
			return
		}

		id := rsp.Message(fmt.Sprintf("You made **%s**, but already have it "+blueCircle, elem3))
		dat.SetMsgElem(id, elem3)

		b.lock.Lock()
		b.dat[m.GuildID] = dat
		b.lock.Unlock()
		return
	}

	dat.SetComb(m.Author.ID, types.Comb{
		Elems: elems,
		Elem3: "",
	})

	b.lock.Lock()
	b.dat[m.GuildID] = dat
	b.lock.Unlock()
	rsp.Resp("That combination doesn't exist! " + types.RedCircle + "\n 	Suggest it by typing **/suggest**")
}
