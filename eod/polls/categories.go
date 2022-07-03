package polls

import (
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

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

	return nil
}

func (b *Polls) UnCategorize(elems []int, catName string, guild string) error {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return nil
	}

	cat, res := db.GetCat(catName)
	if !res.Exists {
		cat = db.NewCat(catName)
	}

	for _, elem := range elems {
		el, res := db.GetElement(elem)
		if !res.Exists {
			return nil
		}
		cat.Lock.Lock()
		delete(cat.Elements, el.ID)
		cat.Lock.Unlock()
	}

	err := db.SaveCat(cat) // Will delete if empty
	if err != nil {
		return err
	}

	return nil
}

func (b *Polls) catImage(guild string, catName string, image string, creator string, changed bool, controversial string, lasted string, news bool) {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return
	}
	cat, res := db.GetCat(catName)
	if !res.Exists {
		vcat, res := db.GetVCat(catName)
		if !res.Exists {
			return
		}

		if vcat.Imager != "" {
			inv := db.GetInv(vcat.Imager)
			inv.CatImagedCnt--
			_ = db.SaveInv(inv)
		}

		vcat.Image = image
		vcat.Imager = creator
		err := db.SaveVCat(vcat)
		if err != nil {
			return
		}

		inv := db.GetInv(creator)
		inv.CatImagedCnt++
		_ = db.SaveInv(inv)

		if news {
			word := "Added"
			if changed {
				word = "Changed"
			}
			b.dg.ChannelMessageSend(db.Config.NewsChannel, "📸 "+word+" Category Image - **"+vcat.Name+"** ("+lasted+"By <@"+creator+">)"+controversial)
		}
		return
	}

	if cat.Imager != "" {
		inv := db.GetInv(cat.Imager)
		inv.CatImagedCnt--
		_ = db.SaveInv(inv)
	}

	cat.Image = image
	cat.Imager = creator
	err := db.SaveCat(cat)
	if err != nil {
		return
	}

	inv := db.GetInv(creator)
	inv.CatImagedCnt++
	_ = db.SaveInv(inv)

	if news {
		word := "Added"
		if changed {
			word = "Changed"
		}
		b.dg.ChannelMessageSend(db.Config.NewsChannel, "📸 "+word+" Category Image - **"+cat.Name+"** ("+lasted+"By <@"+creator+">)"+controversial)
	}
}

func (b *Polls) catColor(guild string, catName string, color int, creator string, controversial string, lasted string, news bool) {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return
	}
	cat, res := db.GetCat(catName)
	if !res.Exists {
		vcat, res := db.GetVCat(catName)
		if !res.Exists {
			return
		}

		if vcat.Colorer != "" {
			inv := db.GetInv(vcat.Colorer)
			inv.CatColoredCnt--
			_ = db.SaveInv(inv)
		}

		vcat.Color = color
		vcat.Colorer = creator
		err := db.SaveVCat(vcat)
		if err != nil {
			return
		}

		inv := db.GetInv(creator)
		inv.CatColoredCnt++
		_ = db.SaveInv(inv)

		if news {
			if color == 0 {
				b.dg.ChannelMessageSend(db.Config.NewsChannel, db.Config.LangProperty("ResetCatColorNews", map[string]any{
					"Category":   vcat.Name,
					"LastedText": lasted,
					"Creator":    creator,
				})+controversial)
			}
			emoji, err := util.GetEmoji(color)
			if err != nil {
				emoji = types.RedCircle
			}
			b.dg.ChannelMessageSend(db.Config.NewsChannel, emoji+" "+db.Config.LangProperty("SetCatColorNews", map[string]any{
				"Category":   vcat.Name,
				"LastedText": lasted,
				"Creator":    creator,
			})+controversial)
		}

		return
	}

	if cat.Colorer != "" {
		inv := db.GetInv(cat.Colorer)
		inv.CatColoredCnt--
		_ = db.SaveInv(inv)
	}

	cat.Color = color
	cat.Colorer = creator
	err := db.SaveCat(cat)
	if err != nil {
		return
	}

	inv := db.GetInv(creator)
	inv.CatColoredCnt++
	_ = db.SaveInv(inv)

	if news {
		if color == 0 {
			b.dg.ChannelMessageSend(db.Config.NewsChannel, db.Config.LangProperty("ResetCatColorNews", map[string]any{
				"Category":   cat.Name,
				"LastedText": lasted,
				"Creator":    creator,
			})+controversial)
		}
		emoji, err := util.GetEmoji(color)
		if err != nil {
			emoji = types.RedCircle
		}
		b.dg.ChannelMessageSend(db.Config.NewsChannel, emoji+" "+db.Config.LangProperty("SetCatColorNews", map[string]any{
			"Category":   cat.Name,
			"LastedText": lasted,
			"Creator":    creator,
		})+controversial)
	}
}

func (b *Polls) deleteVCat(guild string, catName string, creator string, controversial string, lasted string, news bool) {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return
	}
	vcat, res := db.GetVCat(catName)
	if !res.Exists {
		return
	}

	// Update stats
	if vcat.Imager != "" {
		inv := db.GetInv(vcat.Imager)
		inv.ImagedCnt--
		_ = db.SaveInv(inv) // ignore error
	}
	if vcat.Colorer != "" {
		inv := db.GetInv(vcat.Colorer)
		inv.ColoredCnt--
		_ = db.SaveInv(inv) // ignore error
	}

	// Delete
	err := db.DeleteVCat(vcat.Name)
	if err != nil {
		return
	}

	// News
	if news {
		b.dg.ChannelMessageSend(db.Config.NewsChannel, db.Config.LangProperty("DeleteVCatNews", map[string]any{
			"Category":   vcat.Name,
			"LastedText": lasted,
			"Creator":    creator,
		})+controversial)
	}
}

func (b *Polls) catSign(guild string, catName string, mark string, creator string, controversial string, lasted string, news bool) {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return
	}
	cat, res := db.GetCat(catName)
	if !res.Exists {
		vcat, res := db.GetVCat(catName)
		if !res.Exists {
			return
		}

		if vcat.Commenter != "" {
			inv := db.GetInv(vcat.Commenter)
			inv.CatSignedCnt--
			_ = db.SaveInv(inv)
		}

		vcat.Comment = mark
		vcat.Commenter = creator
		err := db.SaveVCat(vcat)
		if err != nil {
			return
		}

		inv := db.GetInv(creator)
		inv.CatSignedCnt++
		_ = db.SaveInv(inv)

		if news {
			b.dg.ChannelMessageSend(db.Config.NewsChannel, "📝 Category Signed - **"+vcat.Name+"** ("+lasted+"By <@"+creator+">)"+controversial)
		}
		return
	}

	if cat.Commenter != "" {
		inv := db.GetInv(cat.Commenter)
		inv.CatSignedCnt--
		_ = db.SaveInv(inv)
	}

	cat.Comment = mark
	cat.Commenter = creator
	err := db.SaveCat(cat)
	if err != nil {
		return
	}

	inv := db.GetInv(creator)
	inv.CatSignedCnt++
	_ = db.SaveInv(inv)

	if news {
		b.dg.ChannelMessageSend(db.Config.NewsChannel, "📝 Category Signed - **"+cat.Name+"** ("+lasted+"By <@"+creator+">)"+controversial)
	}
}
