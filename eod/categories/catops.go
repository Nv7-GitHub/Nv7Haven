package categories

import (
	"net/url"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/sasha-s/go-deadlock"
)

func (b *Categories) DeleteCatCmd(category string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}

	rsp.Acknowledge()

	cat, res := db.GetCat(category)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	category = cat.Name

	// Remove elements
	suggestRm := make([]int, 0)
	rmed := 0
	cat.Lock.RLock()
	db.RLock()
	for elem := range cat.Elements {
		el, res := db.GetElement(elem, true)
		if !res.Exists {
			continue
		}

		if el.Creator == m.Author.ID {
			cat.Lock.RUnlock()
			db.RUnlock()
			err := b.polls.UnCategorize(el.ID, category, m.GuildID)
			cat.Lock.RLock()
			db.RLock()
			rsp.Error(err)
			rmed++
		} else {
			suggestRm = append(suggestRm, el.ID)
		}
	}
	db.RUnlock()
	cat.Lock.RUnlock()

	// Resp
	if len(suggestRm) > 0 {
		err := b.polls.CreatePoll(types.Poll{
			Channel:   db.Config.VotingChannel,
			Guild:     m.GuildID,
			Kind:      types.PollUnCategorize,
			Suggestor: m.Author.ID,

			PollCategorizeData: &types.PollCategorizeData{
				Elems:    suggestRm,
				Category: cat.Name,
				Title:    db.Config.LangProperty("DelCatPoll", nil),
			},
		})
		if rsp.Error(err) {
			return
		}
	}
	b.unCategorizeRsp(rmed, suggestRm, db, category, rsp)
}

func (b *Categories) CatOpCmd(op types.CategoryOperation, lhs string, rhs string, result string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}

	rsp.Acknowledge()

	lcat, res := db.GetCat(lhs)
	if !res.Exists {
		vcat, res := db.GetVCat(lhs)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
		els, res := b.base.CalcVCat(vcat, db)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
		lcat = &types.Category{}
		lcat.Elements = els
		lcat.Name = vcat.Name
		lcat.Lock = &deadlock.RWMutex{}
	}
	rcat, res := db.GetCat(rhs)
	if !res.Exists {
		vcat, res := db.GetVCat(rhs)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
		els, res := b.base.CalcVCat(vcat, db)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
		rcat = &types.Category{}
		rcat.Elements = els
		rcat.Name = vcat.Name
		rcat.Lock = &deadlock.RWMutex{}
	}
	out := make(map[int]types.Empty)

	// Perform operation
	lcat.Lock.RLock()
	rcat.Lock.RLock()
	switch op {
	case types.CatOpUnion:
		for elem := range lcat.Elements {
			out[elem] = types.Empty{}
		}
		for elem := range rcat.Elements {
			out[elem] = types.Empty{}
		}

	case types.CatOpIntersect:
		for elem := range lcat.Elements {
			if _, ok := rcat.Elements[elem]; ok {
				out[elem] = types.Empty{}
			}
		}
		for elem := range rcat.Elements {
			if _, ok := lcat.Elements[elem]; ok {
				out[elem] = types.Empty{}
			}
		}

	case types.CatOpDiff:
		for elem := range lcat.Elements {
			if _, ok := rcat.Elements[elem]; !ok {
				out[elem] = types.Empty{}
			}
		}
	}
	lcat.Lock.RUnlock()
	rcat.Lock.RUnlock()

	// Get result category
	cat, res := db.GetCat(result)
	var els map[int]types.Empty
	if res.Exists {
		result = cat.Name

		// Copy elements
		cat.Lock.RLock()
		els = make(map[int]types.Empty, len(cat.Elements))
		for el := range cat.Elements {
			els[el] = types.Empty{}
		}
		cat.Lock.RUnlock()
	} else if strings.ToLower(result) == result {
		result = util.ToTitle(result)
		els = make(map[int]types.Empty)
		if len(url.PathEscape(result)) > 1024 {
			rsp.ErrorMessage(db.Config.LangProperty("CatNameTooLong", nil))
			return
		}
	}

	// Apply changes
	suggestAdd := make([]int, 0)
	added := 0
	for val := range out {
		el, res := db.GetElement(val)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
		_, exists := els[el.ID]
		if !exists { // Only add if not already added
			if el.Creator == m.Author.ID {
				added++
				err := b.polls.Categorize(el.ID, result, m.GuildID)
				rsp.Error(err)
			} else {
				suggestAdd = append(suggestAdd, el.ID)
			}
		}
	}

	// Save
	if len(suggestAdd) > 0 {
		err := b.polls.CreatePoll(types.Poll{
			Channel:   db.Config.VotingChannel,
			Guild:     m.GuildID,
			Kind:      types.PollCategorize,
			Suggestor: m.Author.ID,

			PollCategorizeData: &types.PollCategorizeData{
				Elems:    suggestAdd,
				Category: result,
				Title:    b.base.CatOpPollTitle(op, db),
			},
		})
		if rsp.Error(err) {
			return
		}
	}

	b.categorizeRsp(added, suggestAdd, db, result, rsp)
}
