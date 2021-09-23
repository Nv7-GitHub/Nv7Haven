package eod

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var autocats = map[string]func(string) bool{
	"Characters": func(s string) bool { return len([]rune(s)) == 1 },
	"Letters":    func(s string) bool { return len([]rune(s)) == 1 && strings.Contains(alphabet, strings.ToLower(s)) },
	"Briks":      func(s string) bool { return strings.HasSuffix(strings.ToLower(s), "brik") },
	"Cheesy":     func(s string) bool { return strings.HasPrefix(strings.ToLower(s), "cheesy") },
	"Bloops":     func(s string) bool { return strings.HasSuffix(strings.ToLower(s), "bloop") },
	"Melons":     func(s string) bool { return strings.HasSuffix(strings.ToLower(s), "melon") },
	"Numbers": func(s string) bool {
		_, err := strconv.ParseFloat(s, 32)
		return err == nil
	},
	"Vukkies":              func(s string) bool { return strings.Contains(strings.ToLower(s), "vukky") },
	"All \"All\" Elements": func(s string) bool { return strings.HasPrefix(strings.ToLower(s), "all ") },
}

func (b *EoD) autocategorize(elem string, guild string) error {
	for catName, catFn := range autocats {
		if catFn(elem) {
			err := b.categorize(elem, catName, guild)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *EoD) categorize(elem string, catName string, guild string) error {
	lock.RLock()
	dat, exists := b.dat[guild]
	lock.RUnlock()
	if !exists {
		return nil
	}

	el, res := dat.GetElement(elem)
	if !res.Exists {
		return nil
	}

	cat, res := dat.GetCategory(catName)
	if !res.Exists {
		cat = types.Category{
			Name:     catName,
			Guild:    guild,
			Elements: make(map[string]types.Empty),
			Image:    "",
		}

		_, err := b.db.Exec("INSERT INTO eod_categories VALUES (?, ?, ?, ?)", guild, cat.Name, "{}", cat.Image)
		if err != nil {
			return err
		}
	} else {
		// Already exists, don't need to do anything
		_, exists = cat.Elements[el.Name]
		if exists {
			return nil
		}
	}

	cat.Elements[el.Name] = types.Empty{}
	dat.SetCategory(cat)

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

	el, res := dat.GetElement(elem)
	if !res.Exists {
		return nil
	}

	cat, res := dat.GetCategory(catName)
	if !res.Exists {
		return nil
	}
	delete(cat.Elements, el.Name)
	dat.SetCategory(cat)

	if len(cat.Elements) == 0 {
		_, err := b.db.Exec("DELETE FROM eod_categories WHERE name=? AND guild=?", cat.Name, cat.Guild)
		if err != nil {
			return err
		}
		dat.DeleteCategory(catName)
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

func (b *EoD) catImage(guild string, catName string, image string, creator string, controversial string) {
	lock.RLock()
	dat, exists := b.dat[guild]
	lock.RUnlock()
	if !exists {
		return
	}
	cat, res := dat.GetCategory(catName)
	if !res.Exists {
		return
	}

	cat.Image = image
	dat.SetCategory(cat)

	lock.Lock()
	b.dat[guild] = dat
	lock.Unlock()

	b.db.Exec("UPDATE eod_categories SET image=? WHERE guild=? AND name=?", image, cat.Guild, cat.Name)
	if creator != "" {
		b.dg.ChannelMessageSend(dat.NewsChannel, "📸 Added Category Image - **"+cat.Name+"** (By <@"+creator+">)"+controversial)
	}
}

func removeDuplicates(elems []string) []string {
	mp := make(map[string]types.Empty, len(elems))
	for _, elem := range elems {
		mp[elem] = types.Empty{}
	}
	out := make([]string, len(mp))
	i := 0
	for k := range mp {
		out[i] = k
		i++
	}
	return out
}
