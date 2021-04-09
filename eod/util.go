package eod

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *EoD) isMod(userID string, guildID string, m msg) (bool, error) {
	lock.RLock()
	dat, inited := b.dat[guildID]
	lock.RUnlock()

	user, err := b.dg.GuildMember(m.GuildID, userID)
	if err != nil {
		return false, err
	}
	hasLoadedRoles := false
	var roles []*discordgo.Role

	for _, roleID := range user.Roles {
		if inited && (roleID == dat.modRole) {
			return true, nil
		}
		role, err := b.dg.State.Role(guildID, roleID)
		if err != nil {
			if !hasLoadedRoles {
				roles, err = b.dg.GuildRoles(m.GuildID)
				if err != nil {
					return false, err
				}
			}

			for _, role := range roles {
				if role.ID == roleID && ((role.Permissions & discordgo.PermissionAdministrator) == discordgo.PermissionAdministrator) {
					return true, nil
				}
			}
		} else {
			if (role.Permissions & discordgo.PermissionAdministrator) == discordgo.PermissionAdministrator {
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

func formatFloat(num float32, prc int) string {
	var (
		zero, dot = "0", "."

		str = fmt.Sprintf("%."+strconv.Itoa(prc)+"f", num)
	)

	return strings.TrimRight(strings.TrimRight(str, zero), dot)
}

func formatInt(n int) string {
	in := strconv.FormatInt(int64(n), 10)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits-- // First character is the - sign (not a digit)
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}

func (b *EoD) getRole(id string, guild string) (*discordgo.Role, error) {
	role, err := b.dg.State.Role(guild, id)
	if err == nil {
		return role, nil
	}

	roles, err := b.dg.GuildRoles(guild)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		if role.ID == id {
			return role, nil
		}
	}

	return nil, errors.New("eod: role not found")
}
