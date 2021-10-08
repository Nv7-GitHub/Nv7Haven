package eod

import (
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/util"
)

func (b *EoD) mark(guild string, elem string, mark string, creator string, controversial string) {
	lock.RLock()
	dat, exists := b.dat[guild]
	lock.RUnlock()
	if !exists {
		return
	}
	el, res := dat.GetElement(elem)
	if !res.Exists {
		return
	}

	el.Comment = mark
	dat.SetElement(el)

	lock.Lock()
	b.dat[guild] = dat
	lock.Unlock()

	query := "UPDATE eod_elements SET comment=? WHERE guild=? AND name LIKE ?"
	if util.IsASCII(el.Name) {
		query = "UPDATE eod_elements SET comment=? WHERE CONVERT(guild USING utf8mb4)=CONVERT(? using utf8mb4) AND CONVERT(name USING utf8mb4) LIKE CONVERT(? USING utf8mb4) COLLATE utf8mb4_general_ci"
	}
	if util.IsWildcard(el.Name) {
		query = strings.ReplaceAll(query, " LIKE ", "=")
	}
	b.db.Exec(query, mark, guild, el.Name)
	if creator != "" {
		b.dg.ChannelMessageSend(dat.NewsChannel, "📝 Signed - **"+el.Name+"** (By <@"+creator+">)"+controversial)
	}
}

func (b *EoD) image(guild string, elem string, image string, creator string, controversial string) {
	lock.RLock()
	dat, exists := b.dat[guild]
	lock.RUnlock()
	if !exists {
		return
	}
	el, res := dat.GetElement(elem)
	if !res.Exists {
		return
	}

	el.Image = image
	dat.SetElement(el)

	lock.Lock()
	b.dat[guild] = dat
	lock.Unlock()

	query := "UPDATE eod_elements SET image=? WHERE guild=? AND name=?"
	if util.IsASCII(el.Name) {
		query = "UPDATE eod_elements SET image=? WHERE CONVERT(guild USING utf8mb4)=? AND CONVERT(name USING utf8mb4)=? COLLATE utf8mb4_general_ci"
	}
	if util.IsWildcard(el.Name) {
		query = strings.ReplaceAll(query, " LIKE ", "=")
	}
	b.db.Exec(query, image, guild, el.Name)
	if creator != "" {
		b.dg.ChannelMessageSend(dat.NewsChannel, "📸 Added Image - **"+el.Name+"** (By <@"+creator+">)"+controversial)
	}
}

func (b *EoD) color(guild string, elem string, color int, creator string, controversial string) {
	lock.RLock()
	dat, exists := b.dat[guild]
	lock.RUnlock()
	if !exists {
		return
	}
	el, res := dat.GetElement(elem)
	if !res.Exists {
		return
	}

	el.Color = color
	dat.SetElement(el)

	lock.Lock()
	b.dat[guild] = dat
	lock.Unlock()

	query := "UPDATE eod_elements SET color=? WHERE guild=? AND name LIKE ?"
	if util.IsASCII(el.Name) {
		query = "UPDATE eod_elements SET color=? WHERE CONVERT(guild USING utf8mb4)=CONVERT(? using utf8mb4) AND CONVERT(name USING utf8mb4) LIKE CONVERT(? USING utf8mb4) COLLATE utf8mb4_general_ci"
	}
	if util.IsWildcard(el.Name) {
		query = strings.ReplaceAll(query, " LIKE ", "=")
	}
	b.db.Exec(query, color, guild, el.Name)
	if creator != "" {
		emoji, err := util.GetEmoji(color)
		if err != nil {
			emoji = redCircle
		}
		b.dg.ChannelMessageSend(dat.NewsChannel, emoji+" Set Color - **"+el.Name+"** (By <@"+creator+">)"+controversial)
	}
}
