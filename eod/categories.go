package eod

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

const x = "❌"
const check = "✅"

func (b *EoD) categoryCmd(elems []string, category string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	if len(category) == 0 {
		rsp.ErrorMessage("Category name can't be blank!")
		return
	}

	inv, exists := dat.invCache[m.Author.ID]
	if !exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}

	suggestAdd := make([]string, 0)
	added := make([]string, 0)
	for _, val := range elems {
		el, exists := dat.elemCache[strings.ToLower(val)]
		if !exists {
			rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", val))
			return
		}

		_, exists = inv[strings.ToLower(el.Name)]
		if !exists {
			rsp.ErrorMessage(fmt.Sprintf("Element **%s** is not in your inventory!", el.Name))
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

func (b *EoD) categorize(elem string, catName string, guild string) error {
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

	cat, exists := dat.catCache[strings.ToLower(catName)]
	if !exists {
		cat = category{
			Name:     catName,
			Guild:    guild,
			Elements: make(map[string]empty),
			Image:    "",
		}

		_, err := b.db.Exec("INSERT INTO eod_categories VALUES (?, ?, ?, ?)", guild, cat.Name, "{}", cat.Image)
		if err != nil {
			return err
		}
	}
	cat.Elements[el.Name] = empty{}
	dat.catCache[strings.ToLower(catName)] = cat

	els, err := json.Marshal(cat.Elements)
	if err != nil {
		return err
	}

	_, err = b.db.Exec("UPDATE eod_categories SET elements=? WHERE guild=? AND name=?", string(els), cat.Guild, cat.Name)
	if err != nil {
		return err
	}

	lock.Lock()
	b.dat[guild] = dat
	lock.Unlock()

	return nil
}

func (b *EoD) unCategorize(elem string, category string, guild string) error {
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

	cat, exists := dat.catCache[strings.ToLower(category)]
	if !exists {
		return nil
	}
	delete(cat.Elements, el.Name)
	if len(cat.Elements) == 0 {
		_, err := b.db.Exec("DELETE FROM eod_categories WHERE name=? AND guild=?", cat.Name, cat.Elements)
		if err != nil {
			return err
		}
	} else {
		data, err := json.Marshal(cat.Elements)
		if err != nil {
			return err
		}
		_, err = b.db.Exec("UPDATE eod_categories SET elements=? WHERE guild=? AND name=?", string(data), cat.Guild, cat.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

const (
	catSortAlphabetical = 0
	catSortByFound      = 1
	catSortByNotFound   = 2
)

func (b *EoD) catCmd(category string, sortKind int, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	inv, exists := dat.invCache[m.Author.ID]
	if !exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}

	cat, exists := dat.catCache[strings.ToLower(category)]
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Category **%s** doesn't exist!", category))
	}

	out := make([]struct {
		found int
		text  string
		name  string
	}, len(cat.Elements))

	found := 0
	i := 0
	fnd := 0
	var text string

	for name := range cat.Elements {
		_, exists := inv[strings.ToLower(name)]
		if exists {
			text = name + " " + check
			found++
			fnd = 1
		} else {
			text = name + " " + x
			fnd = 0
		}

		out[i] = struct {
			found int
			text  string
			name  string
		}{
			found: fnd,
			text:  text,
			name:  name,
		}

		i++
	}

	switch sortKind {
	case catSortAlphabetical:
		sort.Slice(out, func(i, j int) bool {
			return out[i].name < out[j].name
		})

	case catSortByFound:
		sort.Slice(out, func(i, j int) bool {
			return out[i].found < out[j].found
		})

	case catSortByNotFound:
		sort.Slice(out, func(i, j int) bool {
			return out[i].found > out[j].found
		})
	}

	o := make([]string, len(out))
	for i, val := range out {
		o[i] = val.text
	}

	b.newPageSwitcher(pageSwitcher{
		Kind:       pageSwitchInv,
		Title:      fmt.Sprintf("%s (%d, %s%%)", category, len(out), formatFloat(float32(found)/float32(len(out))*100, 2)),
		PageGetter: b.invPageGetter,
		Items:      o,
	}, m, rsp)
}

func (b *EoD) allCatCmd(m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	out := make([]string, len(dat.catCache))

	i := 0
	for _, cat := range dat.catCache {
		out[i] = cat.Name
		i++
	}

	sort.Strings(out)
	b.newPageSwitcher(pageSwitcher{
		Kind:       pageSwitchInv,
		Title:      fmt.Sprintf("All Categories (%d)", len(out)),
		PageGetter: b.invPageGetter,
		Items:      out,
	}, m, rsp)
}

func (b *EoD) rmCategoryCmd(elems []string, category string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	inv, exists := dat.invCache[m.Author.ID]
	if !exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}

	suggestRm := make([]string, 0)
	rmed := make([]string, 0)
	for _, val := range elems {
		el, exists := dat.elemCache[strings.ToLower(val)]
		if !exists {
			rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", val))
			return
		}

		_, exists = inv[strings.ToLower(el.Name)]
		if !exists {
			rsp.ErrorMessage(fmt.Sprintf("Element **%s** is not in your inventory!", el.Name))
			return
		}

		cat, exists := dat.catCache[strings.ToLower(category)]
		if !exists {
			rsp.ErrorMessage(fmt.Sprintf("Category %s doesn't exist!", category))
			return
		}

		_, exists = cat.Elements[el.Name]
		if !exists {
			rsp.ErrorMessage(fmt.Sprintf("Element %s isn't in category %s!", el.Name, cat.Name))
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
