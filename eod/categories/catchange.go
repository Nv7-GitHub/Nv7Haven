package categories

import (
	"net/url"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

func (b *Categories) GetNotExists(db *eodb.DB, elems []string, m types.Msg, rsp types.Rsp) {
	// Not exists message
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
		rsp.ErrorMessage(db.Config.LangProperty("DoesntExist", el))
		return
	}

	rsp.ErrorMessage(db.Config.LangProperty("DoesntExistMultiple", util.JoinTxt(notExists, db.Config.LangProperty("DoesntExistJoiner", nil))))

}

func (b *Categories) CategoryCmd(elems []string, category string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}

	rsp.Acknowledge()

	elems = util.RemoveDuplicates(elems)

	category = strings.TrimSpace(category)

	vcat, res := db.GetVCat(category)
	if res.Exists {
		rsp.ErrorMessage(db.Config.LangProperty("CatAlreadyExist", vcat.Name))
		return
	}

	if len(category) == 0 {
		rsp.ErrorMessage(db.Config.LangProperty("CatNameBlank", nil))
		return
	}

	cat, res := db.GetCat(category)
	var els map[int]types.Empty
	if res.Exists {
		category = cat.Name

		// Copy elements
		cat.Lock.RLock()
		els = make(map[int]types.Empty, len(cat.Elements))
		for el := range cat.Elements {
			els[el] = types.Empty{}
		}
		cat.Lock.RUnlock()
	} else if strings.ToLower(category) == category {
		category = util.ToTitle(category)
		els = make(map[int]types.Empty)
		if len(url.PathEscape(category)) > 1024 {
			rsp.ErrorMessage(db.Config.LangProperty("CatNameTooLong", nil))
			return
		}
	}

	suggestAdd := make([]int, 0)
	added := 0
	for _, val := range elems {
		el, res := db.GetElementByName(val)
		if !res.Exists {
			b.GetNotExists(db, elems, m, rsp)
			return
		}
		_, exists := els[el.ID]
		if !exists { // Only add if not already added
			if el.Creator == m.Author.ID {
				added++
				err := b.polls.Categorize(el.ID, category, m.GuildID)
				rsp.Error(err)
			} else {
				suggestAdd = append(suggestAdd, el.ID)
			}
		}
	}

	if len(suggestAdd) > 0 {
		err := b.polls.CreatePoll(types.Poll{
			Channel:   db.Config.VotingChannel,
			Guild:     m.GuildID,
			Kind:      types.PollCategorize,
			Suggestor: m.Author.ID,

			PollCategorizeData: &types.PollCategorizeData{
				Elems:    suggestAdd,
				Category: category,
				Title:    db.Config.LangProperty("AddCatPoll", nil),
			},
		})
		if rsp.Error(err) {
			return
		}
	}

	b.categorizeRsp(added, suggestAdd, db, category, rsp)
}

func (c *Categories) categorizeRsp(added int, suggestAdd []int, db *eodb.DB, category string, rsp types.Rsp) {
	if added > 0 && len(suggestAdd) == 0 {
		rsp.Message(db.Config.LangProperty("Categorized", nil))
	} else if added == 0 && len(suggestAdd) == 1 {
		el, _ := db.GetElement(suggestAdd[0])
		rsp.Message(db.Config.LangProperty("SuggestCategorized", map[string]any{
			"Element":  el.Name,
			"Category": category,
		}))
	} else if added == 0 && len(suggestAdd) > 1 {
		rsp.Message(db.Config.LangProperty("SuggestCategorizedMult", map[string]any{
			"Elements": len(suggestAdd),
			"Category": category,
		}))
	} else if added > 0 && len(suggestAdd) == 1 {
		el, _ := db.GetElement(suggestAdd[0])
		rsp.Message(db.Config.LangProperty("CategorizeMultSuggestCategorized", map[string]any{
			"Element":  el.Name,
			"Category": category,
		}))
	} else if added > 0 && len(suggestAdd) > 1 {
		rsp.Message(db.Config.LangProperty("CategorizeMultSuggestCategorizedMult", map[string]any{
			"Elements": len(suggestAdd),
			"Category": category,
		}))
	} else {
		rsp.Message(db.Config.LangProperty("Categorized", nil))
	}
}

func (b *Categories) RmCategoryCmd(elems []string, category string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}

	rsp.Acknowledge()

	elems = util.RemoveDuplicates(elems)

	cat, res := db.GetCat(category)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	category = cat.Name

	// Error messages
	notincat := false
	elExists := true
	for _, val := range elems {
		el, res := db.GetElementByName(val)
		if !res.Exists {
			elExists = false
			break
		}

		cat.Lock.RLock()
		_, cont := cat.Elements[el.ID]
		cat.Lock.RUnlock()
		if !cont {
			notincat = true
		}
	}

	if !elExists {
		b.GetNotExists(db, elems, m, rsp)
		return
	}

	if notincat {
		notFound := make(map[string]types.Empty)
		for _, el := range elems {
			elem, _ := db.GetElementByName(el)
			cat.Lock.RLock()
			_, exists := cat.Elements[elem.ID]
			cat.Lock.RUnlock()
			if !exists {
				elem, _ := db.GetElementByName(el)
				notFound["**"+elem.Name+"**"] = types.Empty{}
			}
		}

		if len(notFound) == 1 {
			el := ""
			for k := range notFound {
				el = k
				break
			}
			rsp.ErrorMessage(db.Config.LangProperty("NotInCat", map[string]any{
				"Element":  el,
				"Category": cat.Name,
			}))
			return
		}

		rsp.ErrorMessage(db.Config.LangProperty("NotInCatMult", map[string]any{
			"Elements": util.JoinTxt(notFound, db.Config.LangProperty("DoesntExistJoiner", nil)),
			"Category": cat.Name,
		}))
		return
	}

	// Actually remove
	suggestRm := make([]int, 0)
	toRm := make([]int, 0)
	for _, val := range elems {
		el, res := db.GetElementByName(val)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}

		cat.Lock.RLock()
		_, exists := cat.Elements[el.ID]
		cat.Lock.RUnlock()
		if !exists {
			rsp.ErrorMessage(db.Config.LangProperty("NotInCat", map[string]any{
				"Element":  el.Name,
				"Category": cat.Name,
			}))
			return
		}

		if el.Creator == m.Author.ID {
			toRm = append(toRm, el.ID)
		} else {
			suggestRm = append(suggestRm, el.ID)
		}
	}

	err := b.polls.UnCategorize(toRm, category, m.GuildID)
	if rsp.Error(err) {
		return
	}

	if len(suggestRm) > 0 {
		err := b.polls.CreatePoll(types.Poll{
			Channel:   db.Config.VotingChannel,
			Guild:     m.GuildID,
			Kind:      types.PollUnCategorize,
			Suggestor: m.Author.ID,

			PollCategorizeData: &types.PollCategorizeData{
				Elems:    suggestRm,
				Category: category,
				Title:    db.Config.LangProperty("RmCatPoll", nil),
			},
		})
		if rsp.Error(err) {
			return
		}
	}

	b.unCategorizeRsp(len(toRm), suggestRm, db, category, rsp)
}

func (c *Categories) unCategorizeRsp(rmed int, suggestRm []int, db *eodb.DB, category string, rsp types.Rsp) {
	if rmed > 0 && len(suggestRm) == 0 {
		rsp.Message(db.Config.LangProperty("UnCategorized", nil))
	} else if rmed == 0 && len(suggestRm) == 1 {
		el, _ := db.GetElement(suggestRm[0])
		rsp.Message(db.Config.LangProperty("SuggestUnCategorized", map[string]any{
			"Element":  el.Name,
			"Category": category,
		}))
	} else if rmed == 0 && len(suggestRm) > 1 {
		rsp.Message(db.Config.LangProperty("SuggestUnCategorizedMult", map[string]any{
			"Elements": len(suggestRm),
			"Category": category,
		}))
	} else if rmed > 0 && len(suggestRm) == 1 {
		el, _ := db.GetElement(suggestRm[0])
		rsp.Message(db.Config.LangProperty("UnCategorizeMultSuggestUnCategorized", map[string]any{
			"Element":  el.Name,
			"Category": category,
		}))
	} else if rmed > 0 && len(suggestRm) > 1 {
		rsp.Message(db.Config.LangProperty("UnCategorizeMultSuggestUnCategorizedMult", map[string]any{
			"Elements": len(suggestRm),
			"Category": category,
		}))
	} else {
		rsp.Message(db.Config.LangProperty("UnCategorized", nil))
	}
}
