package eod

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *EoD) categoryCmd(elems []string, category string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	elems = removeDuplicates(elems)

	category = strings.TrimSpace(category)

	for _, elem := range elems {
		if isFoolsMode && !isFool(elem) {
			rsp.ErrorMessage(makeFoolResp(elem))
			return
		}
	}
	if isFoolsMode && !isFool(category) {
		rsp.ErrorMessage(makeFoolResp(category))
		return
	}

	if len(category) == 0 {
		rsp.ErrorMessage("Category name can't be blank!")
		return
	}

	cat, res := dat.GetCategory(category)
	if res.Exists {
		category = cat.Name
	}

	suggestAdd := make([]string, 0)
	added := make([]string, 0)
	for _, val := range elems {
		el, res := dat.GetElement(val)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}

		if el.Creator == m.Author.ID {
			added = append(added, el.Name)
			err := b.categorize(el.Name, category, m.GuildID)
			rsp.Error(err)
		} else {
			suggestAdd = append(suggestAdd, el.Name)
		}
	}
	if len(added) > 0 {
		lock.Lock()
		b.dat[m.GuildID] = dat
		lock.Unlock()
	}
	if len(suggestAdd) > 0 {
		err := b.createPoll(types.Poll{
			Channel: dat.VotingChannel,
			Guild:   m.GuildID,
			Kind:    types.PollCategorize,
			Value1:  category,
			Value4:  m.Author.ID,
			Data:    map[string]interface{}{"elems": suggestAdd},
		})
		if rsp.Error(err) {
			return
		}
	}
	if len(added) > 0 && len(suggestAdd) == 0 {
		rsp.Message("Successfully categorized! 🗃️")
	} else if len(added) == 0 && len(suggestAdd) == 1 {
		rsp.Message(fmt.Sprintf("Suggested to add **%s** to **%s** 🗃️", suggestAdd[0], category))
	} else if len(added) == 0 && len(suggestAdd) > 1 {
		rsp.Message(fmt.Sprintf("Suggested to add **%d elements** to **%s** 🗃️", len(suggestAdd), category))
	} else if len(added) > 0 && len(suggestAdd) == 1 {
		rsp.Message(fmt.Sprintf("Categorized and suggested to add **%s** to **%s** 🗃️", suggestAdd[0], category))
	} else if len(added) > 0 && len(suggestAdd) > 1 {
		rsp.Message(fmt.Sprintf("Categorized and suggested to add **%d elements** to **%s** 🗃️", len(suggestAdd), category))
	} else {
		rsp.Message("Successfully categorized! 🗃️")
	}
}

func (b *EoD) rmCategoryCmd(elems []string, category string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	elems = removeDuplicates(elems)

	cat, res := dat.GetCategory(category)
	if !res.Exists {
		rsp.ErrorMessage(fmt.Sprintf("Category **%s** doesn't exist!", category))
		return
	}

	category = cat.Name

	suggestRm := make([]string, 0)
	rmed := make([]string, 0)
	for _, val := range elems {
		el, res := dat.GetElement(val)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}

		_, exists = cat.Elements[el.Name]
		if !exists {
			rsp.ErrorMessage(fmt.Sprintf("Element **%s** isn't in category **%s**!", el.Name, cat.Name))
			return
		}

		if el.Creator == m.Author.ID {
			rmed = append(rmed, el.Name)
			err := b.unCategorize(el.Name, category, m.GuildID)
			rsp.Error(err)
		} else {
			suggestRm = append(suggestRm, el.Name)
		}
	}
	if len(rmed) > 0 {
		lock.Lock()
		b.dat[m.GuildID] = dat
		lock.Unlock()
	}
	if len(suggestRm) > 0 {
		err := b.createPoll(types.Poll{
			Channel: dat.VotingChannel,
			Guild:   m.GuildID,
			Kind:    types.PollUnCategorize,
			Value1:  category,
			Value4:  m.Author.ID,
			Data:    map[string]interface{}{"elems": suggestRm},
		})
		if rsp.Error(err) {
			return
		}
	}
	if len(rmed) > 0 && len(suggestRm) == 0 {
		rsp.Message("Successfully un-categorized! 🗃️")
	} else if len(rmed) == 0 && len(suggestRm) == 1 {
		rsp.Message(fmt.Sprintf("Suggested to remove **%s** from **%s** 🗃️", suggestRm[0], category))
	} else if len(rmed) == 0 && len(suggestRm) > 1 {
		rsp.Message(fmt.Sprintf("Suggested to remove **%d elements** from **%s** 🗃️", len(suggestRm), category))
	} else if len(rmed) > 0 && len(suggestRm) == 1 {
		rsp.Message(fmt.Sprintf("Un-categorized and suggested to remove **%s** from **%s** 🗃️", suggestRm[0], category))
	} else if len(rmed) > 0 && len(suggestRm) > 1 {
		rsp.Message(fmt.Sprintf("Un-categorized and suggested to remove **%d elements** from **%s** 🗃️", len(suggestRm), category))
	} else {
		rsp.Message("Successfully un-categorized! 🗃️")
	}
}

func (b *EoD) catImgCmd(catName string, url string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	cat, res := dat.GetCategory(catName)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	err := b.createPoll(types.Poll{
		Channel: dat.VotingChannel,
		Guild:   m.GuildID,
		Kind:    types.PollCatImage,
		Value1:  cat.Name,
		Value2:  url,
		Value3:  cat.Image,
		Value4:  m.Author.ID,
	})
	if rsp.Error(err) {
		return
	}
	rsp.Message(fmt.Sprintf("Suggested an image for category **%s** 📷", cat.Name))
}
