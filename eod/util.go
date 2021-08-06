package eod

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"

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

func (b *EoD) saveInv(guild string, user string, newmade bool, recalculate ...bool) {
	lock.RLock()
	dat, exists := b.dat[guild]
	lock.RUnlock()
	if !exists {
		return
	}

	dat.lock.RLock()
	inv := dat.invCache[user]
	dat.lock.RUnlock()

	dat.lock.RLock()
	data, err := json.Marshal(inv)
	dat.lock.RUnlock()
	if err != nil {
		return
	}

	if newmade {
		m := "made+1"
		if len(recalculate) > 0 {
			count := 0
			for val := range inv {
				creator := ""
				dat.lock.RLock()
				elem, exists := dat.elemCache[strings.ToLower(val)]
				dat.lock.RUnlock()
				if exists {
					creator = elem.Creator
				}
				if creator == user {
					count++
				}
			}
			m = strconv.Itoa(count)
		}
		b.db.Exec(fmt.Sprintf("UPDATE eod_inv SET inv=?, count=?, made=%s WHERE guild=? AND user=?", m), data, len(inv), guild, user)
		return
	}

	b.db.Exec("UPDATE eod_inv SET inv=?, count=? WHERE guild=? AND user=?", data, len(inv), guild, user)
}

func (b *EoD) mark(guild string, elem string, mark string, creator string) {
	lock.RLock()
	dat, exists := b.dat[guild]
	lock.RUnlock()
	if !exists {
		return
	}
	dat.lock.RLock()
	el, exists := dat.elemCache[strings.ToLower(elem)]
	dat.lock.RUnlock()
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
	dat.lock.RLock()
	el, exists := dat.elemCache[strings.ToLower(elem)]
	dat.lock.RUnlock()
	if !exists {
		return
	}

	el.Image = image

	dat.lock.Lock()
	dat.elemCache[strings.ToLower(elem)] = el
	dat.lock.Unlock()

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

func (b *EoD) getColor(guild, id string) (int, error) {
	mem, err := b.dg.State.Member(guild, id)
	if err != nil {
		mem, err = b.dg.GuildMember(guild, id)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
	}
	roles := make([]*discordgo.Role, len(mem.Roles))
	for i, roleID := range mem.Roles {
		role, err := b.getRole(roleID, guild)
		if err != nil {
			return 0, err
		}
		roles[i] = role
	}

	sorted := discordgo.Roles(roles)
	sort.Sort(sorted)
	for _, role := range sorted {
		if role.Color != 0 {
			return role.Color, nil
		}
	}

	return 0, errors.New("eod: color not found")
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func isWildcard(s string) bool {
	for _, char := range s {
		if char == '%' || char == '*' || char == '?' || char == '[' || char == ']' || char == '!' || char == '-' || char == '#' || char == '^' || char == '_' {
			return true
		}
	}
	return false
}

var smallWords = map[string]empty{
	"of":  {},
	"an":  {},
	"on":  {},
	"the": {},
	"to":  {},
}

func toTitle(s string) string {
	words := strings.Split(strings.ToLower(s), " ")
	for i, word := range words {
		if len(word) < 1 {
			continue
		}
		w := []rune(word)
		ind := -1

		if w[0] > unicode.MaxASCII {
			continue
		}

		if i == 0 {
			ind = 0
		} else {
			_, exists := smallWords[word]
			if !exists {
				ind = 0
			}
		}

		if w[0] == '(' {
			ind = 1
		}

		if ind != -1 {
			w[ind] = rune(strings.ToUpper(string(word[ind]))[0])
			words[i] = string(w)
		}
	}
	return strings.Join(words, " ")
}
