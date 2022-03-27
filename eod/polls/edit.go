package polls

import (
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

func (b *Polls) mark(guild string, elem int, mark string, creator string, controversial string, lasted string, news bool) {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return
	}
	el, res := db.GetElement(elem)
	if !res.Exists {
		return
	}

	if el.Commenter != "" {
		inv := db.GetInv(el.Commenter)
		inv.SignedCnt--
		_ = db.SaveInv(inv)
	}

	el.Comment = mark
	el.Commenter = creator
	_ = db.SaveElement(el)

	inv := db.GetInv(creator)
	inv.SignedCnt++
	_ = db.SaveInv(inv)

	if news {
		b.dg.ChannelMessageSend(db.Config.NewsChannel, db.Config.LangProperty("SignedElemNews", map[string]any{
			"Element":    el.Name,
			"LastedText": lasted,
			"Creator":    creator,
		})+controversial)
	}
}

func (b *Polls) image(guild string, elem int, image string, creator string, changed bool, controversial string, lasted string, news bool) {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return
	}
	el, res := db.GetElement(elem)
	if !res.Exists {
		return
	}

	if el.Imager != "" {
		inv := db.GetInv(el.Imager)
		inv.ImagedCnt--
		_ = db.SaveInv(inv)
	}

	el.Image = image
	el.Imager = creator
	_ = db.SaveElement(el)

	inv := db.GetInv(creator)
	inv.ImagedCnt++
	_ = db.SaveInv(inv)

	if news {
		newsMsgProp := "AddedImageNews"
		if changed {
			newsMsgProp = "ChangedImageNews"
		}
		b.dg.ChannelMessageSend(db.Config.NewsChannel, db.Config.LangProperty(newsMsgProp, map[string]any{
			"Element":    el.Name,
			"LastedText": lasted,
			"Creator":    creator,
		})+controversial)
	}
}

func (b *Polls) color(guild string, elem int, color int, creator string, controversial string, lasted string, news bool) {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return
	}
	el, res := db.GetElement(elem)
	if !res.Exists {
		return
	}

	if el.Colorer != "" {
		inv := db.GetInv(el.Colorer)
		inv.ColoredCnt--
		_ = db.SaveInv(inv)
	}

	el.Color = color
	el.Colorer = creator
	_ = db.SaveElement(el)

	inv := db.GetInv(creator)
	inv.ColoredCnt++
	_ = db.SaveInv(inv)

	if news {
		emoji, err := util.GetEmoji(color)
		if err != nil {
			emoji = types.RedCircle
		}
		b.dg.ChannelMessageSend(db.Config.NewsChannel, emoji+" "+db.Config.LangProperty("ColoredElemNews", map[string]any{
			"Element":    el.Name,
			"LastedText": lasted,
			"Creator":    creator,
		})+controversial)
	}
}
