package categories

import (
	"github.com/Nv7-Github/Nv7Haven/eod/types"
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
	toRm := make([]int, 0)
	cat.Lock.RLock()
	db.RLock()
	for elem := range cat.Elements {
		el, res := db.GetElement(elem, true)
		if !res.Exists {
			continue
		}

		if el.Creator == m.Author.ID {
			toRm = append(toRm, el.ID)
		} else {
			suggestRm = append(suggestRm, el.ID)
		}
	}
	db.RUnlock()
	cat.Lock.RUnlock()

	err := b.polls.UnCategorize(toRm, category, m.GuildID)
	if rsp.Error(err) {
		return
	}

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
	b.unCategorizeRsp(len(toRm), suggestRm, db, category, rsp)
}
