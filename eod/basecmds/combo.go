package basecmds

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

func (b *BaseCmds) Combine(elems []string, m types.Msg, rsp types.Rsp) {
	ok := b.base.CheckServer(m, rsp)
	if !ok {
		return
	}

	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}

	inv := db.GetInv(m.Author.ID)

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
		rsp.ErrorMessage(fmt.Sprintf(db.Config.LangProperty("MustCombine"), 2))
		return
	}
	if len(elems) > types.MaxComboLength {
		rsp.ErrorMessage(fmt.Sprintf(db.Config.LangProperty("MaxCombine"), types.MaxComboLength))
		return
	}

	// Check if don't have or not exists
	donthave := false
	elExists := true
	for _, elem := range elems {
		id, res := db.GetIDByName(elem)
		if !res.Exists {
			elExists = false
			break
		}

		hasElement := inv.Contains(id)
		if !hasElement {
			donthave = true
		}
	}
	if !elExists {
		notExists := make(map[string]types.Empty)
		for _, el := range elems {
			_, res := db.GetElementByName(el)
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
			rsp.ErrorMessage(fmt.Sprintf(db.Config.LangProperty("DoesntExist"), el))
			return
		}

		rsp.ErrorMessage(fmt.Sprintf(db.Config.LangProperty("DoesntExistMultiple"), util.JoinTxt(notExists, db.Config.LangProperty("DoesntExistJoiner"))))
		return
	}
	if donthave {
		data, _ := b.GetData(m.GuildID)
		_, res := data.GetComb(m.Author.ID)
		if res.Exists {
			data.DeleteComb(m.Author.ID)
		}

		notFound := make(map[string]types.Empty)
		for _, el := range elems {
			id, _ := db.GetIDByName(el)
			exists := inv.Contains(id)
			if !exists {
				elem, _ := db.GetElement(id)
				notFound["**"+elem.Name+"**"] = types.Empty{}
			}
		}

		if len(notFound) == 1 {
			el := ""
			for k := range notFound {
				el = k
				break
			}
			id := rsp.ErrorMessage(fmt.Sprintf(db.Config.LangProperty("DontHave"), el))
			elID, _ := db.GetIDByName(el[2 : len(el)-2])
			data.SetMsgElem(id, elID)
			return
		}

		rsp.ErrorMessage(fmt.Sprintf(db.Config.LangProperty("DontHaveMultiple"), util.JoinTxt(notFound, db.Config.LangProperty("DontHaveJoiner"))))
		return
	}

	// Combine elements
	ids := make([]int, len(elems))
	for i, elem := range elems {
		id, res := db.GetIDByName(elem)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
		ids[i] = id
	}
	elem3, res := db.GetCombo(ids)
	data, _ := b.GetData(m.GuildID)
	if res.Exists {
		data.SetComb(m.Author.ID, types.Comb{
			Elems: ids,
			Elem3: elem3,
		})
		el3, res := db.GetElement(elem3)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}

		inv := db.GetInv(m.Author.ID)

		exists := inv.Contains(elem3)
		if !exists {
			inv.Add(elem3)
			err := db.SaveInv(inv)
			if rsp.Error(err) {
				return
			}

			id := rsp.Message(fmt.Sprintf(db.Config.LangProperty("YouMade"), el3.Name))
			data.SetMsgElem(id, elem3)
			return
		}

		id := rsp.Message(fmt.Sprintf(db.Config.LangProperty("YouHave"), el3.Name))
		data.SetMsgElem(id, elem3)
		return
	}

	data.SetComb(m.Author.ID, types.Comb{
		Elems: ids,
		Elem3: -1,
	})
	rsp.Resp(db.Config.LangProperty("ComboNoExist"))
}
