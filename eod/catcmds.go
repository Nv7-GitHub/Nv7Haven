package eod

import (
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
		return
	}
	category = cat.Name

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
			return out[i].found > out[j].found
		})

	case catSortByNotFound:
		sort.Slice(out, func(i, j int) bool {
			return out[i].found < out[j].found
		})
	}

	o := make([]string, len(out))
	for i, val := range out {
		o[i] = val.text
	}

	b.newPageSwitcher(pageSwitcher{
		Kind:       pageSwitchInv,
		Thumbnail:  cat.Image,
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

func (b *EoD) catImage(guild string, catName string, image string, creator string) {
	lock.RLock()
	dat, exists := b.dat[guild]
	lock.RUnlock()
	if !exists {
		return
	}
	cat, exists := dat.catCache[strings.ToLower(catName)]
	if !exists {
		return
	}

	cat.Image = image
	dat.catCache[strings.ToLower(cat.Name)] = cat

	lock.Lock()
	b.dat[guild] = dat
	lock.Unlock()

	b.db.Exec("UPDATE eod_categories SET image=? WHERE guild=? AND name=?", image, cat.Guild, cat.Name)
	if creator != "" {
		b.dg.ChannelMessageSend(dat.newsChannel, "📸 Added Image - **"+cat.Name+"** (By <@"+creator+">)")
	}
}
