package polls

import (
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

func (b *Polls) mark(guild string, elem int, mark string, creator string, controversial string, news bool) {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return
	}
	el, res := db.GetElement(elem)
	if !res.Exists {
		return
	}

	el.Comment = mark
	el.Commenter = creator
	_ = db.SaveElement(el)

	if el.Commenter != "" {
		inv := db.GetInv(el.Commenter)
		inv.SignedCnt--
		_ = db.SaveInv(inv)
	}
	inv := db.GetInv(creator)
	inv.SignedCnt++
	_ = db.SaveInv(inv)

	if news {
		b.dg.ChannelMessageSend(db.Config.NewsChannel, "📝 Signed - **"+el.Name+"** (By <@"+creator+">)"+controversial)
	}
}

func (b *Polls) image(guild string, elem int, image string, creator string, changed bool, controversial string, news bool) {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return
	}
	el, res := db.GetElement(elem)
	if !res.Exists {
		return
	}

	el.Image = image
	el.Imager = creator
	_ = db.SaveElement(el)

	if el.Imager != "" {
		inv := db.GetInv(el.Imager)
		inv.ImagedCnt--
		_ = db.SaveInv(inv)
	}
	inv := db.GetInv(creator)
	inv.ImagedCnt++
	_ = db.SaveInv(inv)

	if news {
		word := "Added"
		if changed {
			word = "Changed"
		}
		b.dg.ChannelMessageSend(db.Config.NewsChannel, "📸 "+word+" Image - **"+el.Name+"** (By <@"+creator+">)"+controversial)
	}
}

func (b *Polls) color(guild string, elem int, color int, creator string, controversial string, news bool) {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return
	}
	el, res := db.GetElement(elem)
	if !res.Exists {
		return
	}

	el.Color = color
	el.Colorer = creator
	_ = db.SaveElement(el)

	if el.Colorer != "" {
		inv := db.GetInv(el.Colorer)
		inv.ColoredCnt--
		_ = db.SaveInv(inv)
	}
	inv := db.GetInv(creator)
	inv.ColoredCnt++
	_ = db.SaveInv(inv)

	if news {
		emoji, err := util.GetEmoji(color)
		if err != nil {
			emoji = types.RedCircle
		}
		b.dg.ChannelMessageSend(db.Config.NewsChannel, emoji+" Set Color - **"+el.Name+"** (By <@"+creator+">)"+controversial)
	}
}
