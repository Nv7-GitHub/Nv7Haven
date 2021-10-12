package polls

import (
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

func (b *Polls) mark(guild string, elem string, mark string, creator string, controversial string) {
	b.lock.RLock()
	dat, exists := b.dat[guild]
	b.lock.RUnlock()
	if !exists {
		return
	}
	el, res := dat.GetElement(elem)
	if !res.Exists {
		return
	}

	el.Comment = mark
	dat.SetElement(el)

	b.lock.Lock()
	b.dat[guild] = dat
	b.lock.Unlock()

	query := "UPDATE eod_elements SET comment=? WHERE guild=? AND name LIKE ?"
	if util.IsWildcard(el.Name) {
		query = strings.ReplaceAll(query, " LIKE ", "=")
	}
	b.db.Exec(query, mark, guild, el.Name)
	if creator != "" {
		b.dg.ChannelMessageSend(dat.NewsChannel, "📝 Signed - **"+el.Name+"** (By <@"+creator+">)"+controversial)
	}
}

func (b *Polls) image(guild string, elem string, image string, creator string, controversial string) {
	b.lock.RLock()
	dat, exists := b.dat[guild]
	b.lock.RUnlock()
	if !exists {
		return
	}
	el, res := dat.GetElement(elem)
	if !res.Exists {
		return
	}

	el.Image = image
	dat.SetElement(el)

	b.lock.Lock()
	b.dat[guild] = dat
	b.lock.Unlock()

	query := "UPDATE eod_elements SET image=? WHERE guild=? AND name LIKE ?"
	if util.IsWildcard(el.Name) {
		query = strings.ReplaceAll(query, " LIKE ", "=")
	}
	b.db.Exec(query, image, guild, el.Name)
	if creator != "" {
		b.dg.ChannelMessageSend(dat.NewsChannel, "📸 Added Image - **"+el.Name+"** (By <@"+creator+">)"+controversial)
	}
}

func (b *Polls) color(guild string, elem string, color int, creator string, controversial string) {
	b.lock.RLock()
	dat, exists := b.dat[guild]
	b.lock.RUnlock()
	if !exists {
		return
	}
	el, res := dat.GetElement(elem)
	if !res.Exists {
		return
	}

	el.Color = color
	dat.SetElement(el)

	b.lock.Lock()
	b.dat[guild] = dat
	b.lock.Unlock()

	query := "UPDATE eod_elements SET color=? WHERE guild=? AND name LIKE ?"
	if util.IsWildcard(el.Name) {
		query = strings.ReplaceAll(query, " LIKE ", "=")
	}
	b.db.Exec(query, color, guild, el.Name)
	if creator != "" {
		emoji, err := util.GetEmoji(color)
		if err != nil {
			emoji = types.RedCircle
		}
		b.dg.ChannelMessageSend(dat.NewsChannel, emoji+" Set Color - **"+el.Name+"** (By <@"+creator+">)"+controversial)
	}
}
