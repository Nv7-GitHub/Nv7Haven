package eod

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *EoD) isMod(userID string, m msg) (bool, error) {
	user, err := b.dg.GuildMember(m.GuildID, userID)
	if err != nil {
		return false, err
	}
	roles, err := b.dg.GuildRoles(m.GuildID)
	if err != nil {
		return false, err
	}

	for _, roleID := range user.Roles {
		for _, role := range roles {
			if role.ID == roleID && ((role.Permissions & discordgo.PermissionAdministrator) == discordgo.PermissionAdministrator) {
				return true, nil
			}
		}
	}
	return false, nil
}

func (b *EoD) saveInv(guild string, user string) {
	lock.RLock()
	dat, exists := b.dat[guild]
	lock.RUnlock()
	if !exists {
		return
	}

	data, err := json.Marshal(dat.invCache[user])
	if err != nil {
		return
	}

	b.db.Exec("UPDATE eod_inv SET inv=?, count=? WHERE guild=? AND user=?", data, len(dat.invCache[user]), guild, user)
}

func (b *EoD) mark(guild string, elem string, mark string, creator string) {
	lock.RLock()
	dat, exists := b.dat[guild]
	lock.RUnlock()
	if !exists {
		return
	}
	el, exists := dat.elemCache[strings.ToLower(elem)]
	if !exists {
		return
	}

	el.Comment = mark
	dat.elemCache[strings.ToLower(elem)] = el

	lock.Lock()
	b.dat[guild] = dat
	lock.Unlock()

	b.db.Exec("UPDATE eod_elements SET comment=? WHERE guild=? AND name=?", mark, guild, el.Name)
	if creator != "" {
		b.dg.ChannelMessageSend(dat.newsChannel, "📝 Signed - **"+el.Name+"** (By <@"+creator+">)")
	}
}

func (b *EoD) image(guild string, elem string, image string, creator string) {
	lock.RLock()
	dat, exists := b.dat[guild]
	lock.RUnlock()
	if !exists {
		return
	}
	el, exists := dat.elemCache[strings.ToLower(elem)]
	if !exists {
		return
	}

	el.Image = image
	dat.elemCache[strings.ToLower(elem)] = el

	lock.Lock()
	b.dat[guild] = dat
	lock.Unlock()

	b.db.Exec("UPDATE eod_elements SET image=? WHERE guild=? AND name=?", image, guild, el.Name)
	if creator != "" {
		b.dg.ChannelMessageSend(dat.newsChannel, "📸 Added Image - **"+el.Name+"** (By <@"+creator+">)")
	}
}

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

	rsp.Embed(&discordgo.MessageEmbed{
		Title:       el.Name + " Info",
		Description: fmt.Sprintf("Created by <@%s>\nCreated on %s\nComplexity: %d\n<@%s> **You %shave this.**\n\n%s", el.Creator, el.CreatedOn.Format("January 2, 2006, 3:04 PM"), el.Complexity, m.Author.ID, has, el.Comment),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: el.Image,
		},
	})
}
