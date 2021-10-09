package eod

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

func (b *EoD) categoryCmd(elems []string, category string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
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
		rsp.ErrorMessage("Category name can't be blank!")
		return
	}

	cat, res := dat.GetCategory(category)
	if res.Exists {
		category = cat.Name
	} else if strings.ToLower(category) == category {
		category = util.ToTitle(category)
	}

	suggestAdd := make([]string, 0)
	added := make([]string, 0)
	for _, val := range elems {
		el, res := dat.GetElement(val)
		if !res.Exists {
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

			rsp.ErrorMessage("Elements " + joinTxt(notExists, "and") + " don't exist!")
			return
		}

		if el.Creator == m.Author.ID {
			added = append(added, el.Name)
			err := b.polls.Categorize(el.Name, category, m.GuildID)
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
		err := b.polls.CreatePoll(types.Poll{
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

	elems = util.RemoveDuplicates(elems)

	cat, res := dat.GetCategory(category)
	if !res.Exists {
		rsp.ErrorMessage(fmt.Sprintf("Category **%s** doesn't exist!", category))
		return
	}

	category = cat.Name

	// Error messages
	notincat := false
	elExists := true
	for _, val := range elems {
		el, res := dat.GetElement(val)
		if !res.Exists {
			elExists = false
			break
		}

		_, cont := cat.Elements[el.Name]
		if !cont {
			notincat = true
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

		rsp.ErrorMessage("Elements " + joinTxt(notExists, "and") + " don't exist!")
		return
	}

	if notincat {
		_, res := dat.GetComb(m.Author.ID)
		if res.Exists {
			dat.DeleteComb(m.Author.ID)

			lock.Lock()
			b.dat[m.GuildID] = dat
			lock.Unlock()
		}

		notFound := make(map[string]types.Empty)
		for _, el := range elems {
			elem, _ := dat.GetElement(el)
			_, exists := cat.Elements[elem.Name]
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
			rsp.ErrorMessage(fmt.Sprintf("Element **%s** isn't in category **%s**!", el, cat.Name))
			return
		}

		rsp.ErrorMessage(fmt.Sprintf("Elements %s aren't in category **%s**!", joinTxt(notFound, "and"), cat.Name))
		return
	}

	// Actually remove
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
			err := b.polls.UnCategorize(el.Name, category, m.GuildID)
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
		err := b.polls.CreatePoll(types.Poll{
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
