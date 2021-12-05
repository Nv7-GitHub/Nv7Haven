package polls

import (
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

func (b *Polls) mark(guild string, elem string, mark string, creator string, controversial string) {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return
	}
	el, res := db.GetElementByName(elem)
	if !res.Exists {
		return
	}

	el.Comment = mark
	_ = db.SaveElement(el)

	if creator != "" {
		b.dg.ChannelMessageSend(db.Config.NewsChannel, "📝 Signed - **"+el.Name+"** (By <@"+creator+">)"+controversial)
	}
}

func (b *Polls) image(guild string, elem string, image string, creator string, controversial string) {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return
	}
	el, res := db.GetElementByName(elem)
	if !res.Exists {
		return
	}

	el.Image = image
	_ = db.SaveElement(el)
	if creator != "" {
		b.dg.ChannelMessageSend(db.Config.NewsChannel, "📸 Added Image - **"+el.Name+"** (By <@"+creator+">)"+controversial)
	}
}

func (b *Polls) color(guild string, elem string, color int, creator string, controversial string) {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return
	}
	el, res := db.GetElementByName(elem)
	if !res.Exists {
		return
	}

	el.Color = color
	_ = db.SaveElement(el)
	if creator != "" {
		emoji, err := util.GetEmoji(color)
		if err != nil {
			emoji = types.RedCircle
		}
		b.dg.ChannelMessageSend(db.Config.NewsChannel, emoji+" Set Color - **"+el.Name+"** (By <@"+creator+">)"+controversial)
	}
}
