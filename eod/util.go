package eod

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// element, guild, element, element, guild, guild, element - returns: made by x combos, used in x combos, found by x people
const elemInfoDataCount = `SELECT a.cnt, b.cnt, c.cnt FROM (SELECT COUNT(1) AS cnt FROM eod_combos WHERE elem3=? AND guild=?) a, (SELECT COUNT(1) AS cnt FROM eod_combos WHERE (elem1=?) OR (elem2=?) AND guild=?) b, (SELECT COUNT(1) as cnt FROM eod_inv WHERE guild=? AND (JSON_EXTRACT(inv, CONCAT("$.", LOWER(?))) IS NOT NULL)) c`

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

func formatFloat(num float32, prc int) string {
	var (
		zero, dot = "0", "."

		str = fmt.Sprintf("%."+strconv.Itoa(prc)+"f", num)
	)

	return strings.TrimRight(strings.TrimRight(str, zero), dot)
}
