package eod

import (
	"encoding/json"
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

	b.db.Exec("UPDATE eod_inv SET inv=? WHERE guild=? AND user=?", data, guild, user)
}

func (b *EoD) mark(guild string, elem string, mark string) {
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
}
