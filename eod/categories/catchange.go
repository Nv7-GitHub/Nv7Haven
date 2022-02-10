package categories

import (
	"net/url"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
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

	elems = util.RemoveDuplicates(elems)

	category = strings.TrimSpace(category)

	for _, elem := range elems {
		if base.IsFoolsMode && !base.IsFool(elem) {
			rsp.ErrorMessage(base.MakeFoolResp(elem))
			return
		}
	}
	if base.IsFoolsMode && !base.IsFool(category) {
		rsp.ErrorMessage(base.MakeFoolResp(category))
		return
	}

	if len(category) == 0 {
		rsp.ErrorMessage(db.Config.LangProperty("CatNameBlank", nil))
		return
	}

	cat, res := db.GetCat(category)
	if res.Exists {
		category = cat.Name
	} else if strings.ToLower(category) == category {
		category = util.ToTitle(category)
		if len(url.PathEscape(category)) > 1024 {
			rsp.ErrorMessage(db.Config.LangProperty("CatNameTooLong", nil))
			return
		}
	}

	suggestAdd := make([]int, 0)
	added := make([]string, 0)
	for _, val := range elems {
		el, res := db.GetElementByName(val)
		if !res.Exists {
			b.GetNotExists(db, elems, m, rsp)
			return
		}

		if el.Creator == m.Author.ID {
			added = append(added, el.Name)
			err := b.polls.Categorize(el.ID, category, m.GuildID)
			rsp.Error(err)
		} else {
			suggestAdd = append(suggestAdd, el.ID)
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
			},
		})
		if rsp.Error(err) {
			return
		}
	}
	if len(added) > 0 && len(suggestAdd) == 0 {
		rsp.Message(db.Config.LangProperty("Categorized", nil))
	} else if len(added) == 0 && len(suggestAdd) == 1 {
		el, _ := db.GetElement(suggestAdd[0])
		rsp.Message(db.Config.LangProperty("SuggestCategorized", map[string]interface{}{
			"Element":  el.Name,
			"Category": category,
		}))
	} else if len(added) == 0 && len(suggestAdd) > 1 {
		rsp.Message(db.Config.LangProperty("SuggestCategorizedMult", map[string]interface{}{
			"Elements": len(suggestAdd),
			"Category": category,
		}))
	} else if len(added) > 0 && len(suggestAdd) == 1 {
		el, _ := db.GetElement(suggestAdd[0])
		rsp.Message(db.Config.LangProperty("CategorizeMultSuggestCategorized", map[string]interface{}{
			"Element":  el.Name,
			"Category": category,
		}))
	} else if len(added) > 0 && len(suggestAdd) > 1 {
		rsp.Message(db.Config.LangProperty("CategorizeMultSuggestCategorizedMult", map[string]interface{}{
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

	elems = util.RemoveDuplicates(elems)

	cat, res := db.GetCat(category)
	if !res.Exists {
		rsp.ErrorMessage(db.Config.LangProperty("CatNoExist", category))
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
			rsp.ErrorMessage(db.Config.LangProperty("NotInCat", map[string]interface{}{
				"Element":  el,
				"Category": cat.Name,
			}))
			return
		}

		rsp.ErrorMessage(db.Config.LangProperty("NotInCatMult", map[string]interface{}{
			"Elements": util.JoinTxt(notFound, db.Config.LangProperty("DoesntExistJoiner", nil)),
			"Category": cat.Name,
		}))
		return
	}

	// Actually remove
	suggestRm := make([]int, 0)
	rmed := make([]string, 0)
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
			rsp.ErrorMessage(db.Config.LangProperty("NotInCat", map[string]interface{}{
				"Element":  el.Name,
				"Category": cat.Name,
			}))
			return
		}

		if el.Creator == m.Author.ID {
			rmed = append(rmed, el.Name)
			err := b.polls.UnCategorize(el.ID, category, m.GuildID)
			rsp.Error(err)
		} else {
			suggestRm = append(suggestRm, el.ID)
		}
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
			},
		})
		if rsp.Error(err) {
			return
		}
	}
	if len(rmed) > 0 && len(suggestRm) == 0 {
		rsp.Message(db.Config.LangProperty("UnCategorized", nil))
	} else if len(rmed) == 0 && len(suggestRm) == 1 {
		el, _ := db.GetElement(suggestRm[0])
		rsp.Message(db.Config.LangProperty("SuggestUnCategorized", map[string]interface{}{
			"Element":  el.Name,
			"Category": category,
		}))
	} else if len(rmed) == 0 && len(suggestRm) > 1 {
		rsp.Message(db.Config.LangProperty("SuggestUnCategorizedMult", map[string]interface{}{
			"Elements": len(suggestRm),
			"Category": category,
		}))
	} else if len(rmed) > 0 && len(suggestRm) == 1 {
		el, _ := db.GetElement(suggestRm[0])
		rsp.Message(db.Config.LangProperty("UnCategorizeMultSuggestUnCategorized", map[string]interface{}{
			"Element":  el.Name,
			"Category": category,
		}))
	} else if len(rmed) > 0 && len(suggestRm) > 1 {
		rsp.Message(db.Config.LangProperty("UnCategorizeMultSuggestUnCategorizedMult", map[string]interface{}{
			"Elements": len(suggestRm),
			"Category": category,
		}))
	} else {
		rsp.Message(db.Config.LangProperty("UnCategorized", nil))
	}
}
