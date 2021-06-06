package eod

import (
	"fmt"
	"strings"
)

func (b *EoD) categoryCmd(elems []string, category string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	category = strings.TrimSpace(category)

	if len(category) == 0 {
		rsp.ErrorMessage("Category name can't be blank!")
		return
	}

	cat, exists := dat.catCache[strings.ToLower(category)]
	if exists {
		category = cat.Name
	}

	suggestAdd := make([]string, 0)
	added := make([]string, 0)
	for _, val := range elems {
		el, exists := dat.elemCache[strings.ToLower(val)]
		if !exists {
			rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", val))
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
		err := b.createPoll(poll{
			Channel: dat.votingChannel,
			Guild:   m.GuildID,
			Kind:    pollCategorize,
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

func (b *EoD) rmCategoryCmd(elems []string, category string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	cat, exists := dat.catCache[strings.ToLower(category)]
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Category **%s** doesn't exist!", category))
		return
	}

	category = cat.Name

	suggestRm := make([]string, 0)
	rmed := make([]string, 0)
	for _, val := range elems {
		el, exists := dat.elemCache[strings.ToLower(val)]
		if !exists {
			rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", val))
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
		err := b.createPoll(poll{
			Channel: dat.votingChannel,
			Guild:   m.GuildID,
			Kind:    pollUnCategorize,
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
		rsp.Message(fmt.Sprintf("Un-categorized and suggested to remove **%d elements** tfrom**%s** 🗃️", len(suggestRm), category))
	} else {
		rsp.Message("Successfully un-categorized! 🗃️")
	}
}

func (b *EoD) catImgCmd(catName string, url string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	cat, exists := dat.catCache[strings.ToLower(catName)]
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Category **%s** doesn't exist!", catName))
		return
	}

	err := b.createPoll(poll{
		Channel: dat.votingChannel,
		Guild:   m.GuildID,
		Kind:    pollCatImage,
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
