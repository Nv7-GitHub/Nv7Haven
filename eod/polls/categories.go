package polls

import (
	"errors"
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
	db, res := b.GetDB(guild)
	if !res.Exists {
		return errors.New(res.Message)
	}

	for catName, catFn := range autocats {
		if catFn(elem) {
			id, res := db.GetIDByName(elem)
			if !res.Exists {
				return errors.New(res.Message)
			}
			err := b.Categorize(id, catName, guild)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Polls) Categorize(elem int, catName string, guild string) error {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return nil
	}

	el, res := db.GetElement(elem)
	if !res.Exists {
		return nil
	}

	cat, res := db.GetCat(catName)
	if !res.Exists {
		cat = db.NewCat(catName)
	} else {
		// Already exists, don't need to do anything
		_, exists := cat.Elements[el.ID]
		if exists {
			return nil
		}
	}

	cat.Lock.Lock()
	cat.Elements[el.ID] = types.Empty{}
	cat.Lock.Unlock()
	err := db.SaveCat(cat)
	if err != nil {
		return err
	}

	if types.RecalcAutocats {
		fmt.Println(el.Name)
	}

	return nil
}

func (b *Polls) UnCategorize(elem int, catName string, guild string) error {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return nil
	}

	el, res := db.GetElement(elem)
	if !res.Exists {
		return nil
	}

	cat, res := db.GetCat(catName)
	if !res.Exists {
		cat = db.NewCat(catName)
	}
	cat.Lock.Lock()
	delete(cat.Elements, el.ID)
	cat.Lock.Unlock()
	err := db.SaveCat(cat) // Will delete if empty
	if err != nil {
		return err
	}

	return nil
}

func (b *Polls) catImage(guild string, catName string, image string, creator string, changed bool, controversial string, news bool) {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return
	}
	cat, res := db.GetCat(catName)
	if !res.Exists {
		return
	}

	cat.Image = image
	err := db.SaveCat(cat)
	if err != nil {
		return
	}

	if cat.Imager != "" {
		inv := db.GetInv(cat.Imager)
		inv.CatImagedCnt--
		_ = db.SaveInv(inv)
	}
	inv := db.GetInv(creator)
	inv.CatImagedCnt++
	_ = db.SaveInv(inv)

	if news {
		word := "Added"
		if changed {
			word = "Changed"
		}
		b.dg.ChannelMessageSend(db.Config.NewsChannel, "📸 "+word+" Category Image - **"+cat.Name+"** (By <@"+creator+">)"+controversial)
	}
}

func (b *Polls) catColor(guild string, catName string, color int, creator string, controversial string, news bool) {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return
	}
	cat, res := db.GetCat(catName)
	if !res.Exists {
		return
	}

	cat.Color = color
	err := db.SaveCat(cat)
	if err != nil {
		return
	}

	if cat.Colorer != "" {
		inv := db.GetInv(cat.Colorer)
		inv.CatColoredCnt--
		_ = db.SaveInv(inv)
	}
	inv := db.GetInv(creator)
	inv.CatColoredCnt++
	_ = db.SaveInv(inv)

	if news {
		if color == 0 {
			b.dg.ChannelMessageSend(db.Config.NewsChannel, "Reset Category Color - **"+cat.Name+"** (By <@"+creator+">)"+controversial)
		}
		emoji, err := util.GetEmoji(color)
		if err != nil {
			emoji = types.RedCircle
		}
		b.dg.ChannelMessageSend(db.Config.NewsChannel, emoji+" Set Category Color - **"+cat.Name+"** (By <@"+creator+">)"+controversial)
	}
}
