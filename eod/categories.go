package eod

import (
	"encoding/json"
	"strings"
)

func (b *EoD) categorize(elem string, catName string, guild string) error {
	lock.RLock()
	dat, exists := b.dat[guild]
	lock.RUnlock()
	if !exists {
		return nil
	}
	dat.lock.RLock()
	el, exists := dat.elemCache[strings.ToLower(elem)]
	if !exists {
		return nil
	}
	dat.lock.RUnlock()

	if dat.catCache == nil {
		dat.catCache = make(map[string]category)
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

func (b *EoD) unCategorize(elem string, catName string, guild string) error {
	lock.RLock()
	dat, exists := b.dat[guild]
	lock.RUnlock()
	if !exists {
		return nil
	}
	dat.lock.RLock()
	el, exists := dat.elemCache[strings.ToLower(elem)]
	if !exists {
		return nil
	}
	dat.lock.RUnlock()

	if dat.catCache == nil {
		dat.catCache = make(map[string]category)
	}

	cat, exists := dat.catCache[strings.ToLower(catName)]
	if !exists {
		return nil
	}
	delete(cat.Elements, el.Name)
	dat.catCache[strings.ToLower(catName)] = cat

	if len(cat.Elements) == 0 {
		_, err := b.db.Exec("DELETE FROM eod_categories WHERE name=? AND guild=?", cat.Name, cat.Guild)
		if err != nil {
			return err
		}
		delete(dat.catCache, strings.ToLower(catName))
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

	lock.Lock()
	b.dat[guild] = dat
	lock.Unlock()

	return nil
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
		b.dg.ChannelMessageSend(dat.newsChannel, "📸 Added Category Image - **"+cat.Name+"** (By <@"+creator+">)")
	}
}

func removeDuplicates(elems []string) []string {
	mp := make(map[string]empty, len(elems))
	for _, elem := range elems {
		mp[elem] = empty{}
	}
	out := make([]string, len(mp))
	i := 0
	for k := range mp {
		out[i] = k
		i++
	}
	return out
}
