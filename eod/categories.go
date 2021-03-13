package eod

import (
	"encoding/json"
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
	suggestAdd := make([]string, 0)
	added := make([]string, 0)
	for _, val := range elems {
		el, exists := dat.elemCache[strings.ToLower(val)]
		if !exists {
			rsp.ErrorMessage(fmt.Sprintf("Element %s doesn't exist!", val))
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
			Channel: m.ChannelID,
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
		rsp.Resp("Successfully categorized! 🗃️")
	} else if len(added) == 0 && len(suggestAdd) == 1 {
		rsp.Resp(fmt.Sprintf("Suggested to add **%s** to **%s** 🗃️", suggestAdd[0], category))
	} else if len(added) == 0 && len(suggestAdd) > 1 {
		rsp.Resp(fmt.Sprintf("Suggested to add **%d elements** to **%s** 🗃️", len(suggestAdd), category))
	} else if len(added) > 0 && len(suggestAdd) == 1 {
		rsp.Resp(fmt.Sprintf("Categorized and suggested to add **%s** to **%s** 🗃️", suggestAdd[0], category))
	} else if len(added) > 0 && len(suggestAdd) > 1 {
		rsp.Resp(fmt.Sprintf("Categorized and suggested to add **%d elements** to **%s** 🗃️", len(suggestAdd), category))
	} else {
		rsp.Resp("Successfully categorized! 🗃️")
	}
}

func (b *EoD) categorize(elem string, category string, guild string) error {
	lock.RLock()
	dat, exists := b.dat[guild]
	lock.RUnlock()
	if !exists {
		return nil
	}
	el, exists := dat.elemCache[strings.ToLower(elem)]
	if !exists {
		return nil
	}
	el.Categories[category] = empty{}
	dat.elemCache[strings.ToLower(elem)] = el

	data, err := json.Marshal(el.Categories)
	if err != nil {
		return err
	}
	_, err = b.db.Exec("UPDATE eod_elements SET categories=? WHERE guild=? AND name=?", string(data), el.Guild, el.Name)
	return err
}
