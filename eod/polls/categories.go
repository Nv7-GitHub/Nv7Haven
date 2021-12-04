package polls

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
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
	"Amogus in a...":       func(s string) bool { return strings.HasPrefix(s, "Amogus in a") },
}

func (b *Polls) Autocategorize(elem string, guild string) error {
	for catName, catFn := range autocats {
		if catFn(elem) {
			err := b.Categorize(elem, catName, guild)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Polls) Categorize(elem string, catName string, guild string) error {
	b.lock.RLock()
	dat, exists := b.dat[guild]
	b.lock.RUnlock()
	if !exists {
		return nil
	}

	el, res := dat.GetElement(elem)
	if !res.Exists {
		return nil
	}

	cat, res := dat.GetCategory(catName)
	if !res.Exists {
		cat = types.OldCategory{
			Name:     catName,
			Guild:    guild,
			Elements: make(map[string]types.Empty),
			Image:    "",
		}

		_, err := b.db.Exec("INSERT INTO eod_categories VALUES (?, ?, ?, ?, ?)", guild, cat.Name, "{}", cat.Image, cat.Color)
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

	if types.RecalcAutocats {
		fmt.Println(el.Name)
	}

	_, err = b.db.Exec("UPDATE eod_categories SET elements=? WHERE guild=? AND name=?", string(els), cat.Guild, cat.Name)
	if err != nil {
		return err
	}

	b.lock.Lock()
	b.dat[guild] = dat
	b.lock.Unlock()

	return nil
}

func (b *Polls) UnCategorize(elem string, catName string, guild string) error {
	b.lock.RLock()
	dat, exists := b.dat[guild]
	b.lock.RUnlock()
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

	b.lock.Lock()
	b.dat[guild] = dat
	b.lock.Unlock()

	return nil
}

func (b *Polls) catImage(guild string, catName string, image string, creator string, controversial string) {
	b.lock.RLock()
	dat, exists := b.dat[guild]
	b.lock.RUnlock()
	if !exists {
		return
	}
	cat, res := dat.GetCategory(catName)
	if !res.Exists {
		return
	}

	cat.Image = image
	dat.SetCategory(cat)

	b.lock.Lock()
	b.dat[guild] = dat
	b.lock.Unlock()

	query := "UPDATE eod_categories SET image=? WHERE guild=? AND name LIKE ?"
	if util.IsWildcard(cat.Name) {
		query = strings.ReplaceAll(query, " LIKE ", "=")
	}
	b.db.Exec(query, image, cat.Guild, cat.Name)
	if creator != "" {
		b.dg.ChannelMessageSend(dat.NewsChannel, "📸 Added Category Image - **"+cat.Name+"** (By <@"+creator+">)"+controversial)
	}
}

func (b *Polls) catColor(guild string, catName string, color int, creator string, controversial string) {
	b.lock.RLock()
	dat, exists := b.dat[guild]
	b.lock.RUnlock()
	if !exists {
		return
	}
	cat, res := dat.GetCategory(catName)
	if !res.Exists {
		return
	}

	cat.Color = color
	dat.SetCategory(cat)

	b.lock.Lock()
	b.dat[guild] = dat
	b.lock.Unlock()

	query := "UPDATE eod_categories SET color=? WHERE guild=? AND name LIKE ?"
	if util.IsWildcard(cat.Name) {
		query = strings.ReplaceAll(query, " LIKE ", "=")
	}
	b.db.Exec(query, color, cat.Guild, cat.Name)
	if creator != "" {
		if color == 0 {
			b.dg.ChannelMessageSend(dat.NewsChannel, "Reset Category Color - **"+cat.Name+"** (By <@"+creator+">)"+controversial)
		}
		emoji, err := util.GetEmoji(color)
		if err != nil {
			emoji = types.RedCircle
		}
		b.dg.ChannelMessageSend(dat.NewsChannel, emoji+" Set Category Color - **"+cat.Name+"** (By <@"+creator+">)"+controversial)
	}
}
