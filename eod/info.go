package eod

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *EoD) infoCmd(elem string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	el, exists := dat.elemCache[strings.ToLower(elem)]
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Element %s doesn't exist!", elem))
		return
	}

	has := ""
	exists = false
	if dat.invCache != nil {
		_, exists = dat.invCache[m.Author.ID]
		if exists {
			_, exists = dat.invCache[m.Author.ID][strings.ToLower(el.Name)]
		}
	}
	if !exists {
		has = "don't "
	}

	row := b.db.QueryRow(elemInfoDataCount, el.Name, el.Guild, el.Name, el.Name, el.Guild, el.Guild, el.Name)
	var madeby int
	var usedby int
	var foundby int
	err := row.Scan(&madeby, &usedby, &foundby)
	if rsp.Error(err) {
		return
	}

	usedbysuff := "s"
	if usedby == 1 {
		usedbysuff = ""
	}
	madebysuff := "s"
	if madeby == 1 {
		madebysuff = ""
	}
	foundbysuff := "s"
	if foundby == 1 {
		foundbysuff = ""
	}

	rsp.Embed(&discordgo.MessageEmbed{
		Title:       el.Name + " Info",
		Description: fmt.Sprintf("Created by <@%s>\nCreated on %s\nUsed in %d combo%s\nMade with %d combo%s\nFound by %d player%s\nComplexity: %d\nDifficulty: %d\n<@%s> **You %shave this.**\n\n%s", el.Creator, el.CreatedOn.Format("January 2, 2006, 3:04 PM"), usedby, usedbysuff, madeby, madebysuff, foundby, foundbysuff, el.Complexity, el.Difficulty, m.Author.ID, has, el.Comment),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: el.Image,
		},
	})
}
