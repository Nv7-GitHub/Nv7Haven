package eod

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

// Unneeded for now
/*func (b *EoD) getRoles(userID string, guild string) ([]*discordgo.Role, error) {
	user, err := b.dg.GuildMember(guild, userID)
	if err != nil {
		return nil, err
	}
	hasLoadedRoles := false
	var roles []*discordgo.Role
	out := make([]*discordgo.Role, len(user.Roles))

	for i, roleID := range user.Roles {
		role, err := b.dg.State.Role(guild, roleID)
		if err != nil {
			if !hasLoadedRoles {
				roles, err = b.dg.GuildRoles(guild)
				if err != nil {
					return nil, err
				}
			}

			for _, role := range roles {
				if role.ID == roleID {
					roles[i] = role
				}
			}
		} else {
			roles[i] = role
		}
	}
	return out, nil
}*/

func (b *EoD) isMod(userID string, guildID string, m types.Msg) (bool, error) {
	lock.RLock()
	dat, inited := b.dat[guildID]
	lock.RUnlock()

	user, err := b.dg.GuildMember(m.GuildID, userID)
	if err != nil {
		return false, err
	}
	if (user.Permissions * discordgo.PermissionAdministrator) == discordgo.PermissionAdministrator {
		return true, nil
	}

	hasLoadedRoles := false
	var roles []*discordgo.Role

	for _, roleID := range user.Roles {
		if inited && (roleID == dat.ModRole) {
			return true, nil
		}
		role, err := b.dg.State.Role(guildID, roleID)
		if err != nil {
			if !hasLoadedRoles {
				roles, err = b.dg.GuildRoles(m.GuildID)
				if err != nil {
					return false, err
				}
				hasLoadedRoles = true
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

	inv, _ := dat.GetInv(user, true)

	dat.Lock.RLock()
	data, err := json.Marshal(inv)
	dat.Lock.RUnlock()
	if err != nil {
		return
	}

	if newmade {
		m := "made+1"
		if len(recalculate) > 0 {
			count := 0
			for val := range inv {
				creator := ""

				elem, res := dat.GetElement(val)
				if res.Exists {
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

// FOOLS
//go:embed fools.txt
var foolsRaw string
var fools []string

var isFoolsMode = time.Now().Month() == time.April && time.Now().Day() == 1

func isFool(inp string) bool {
	for _, val := range fools {
		if strings.Contains(inp, val) {
			return true
		}
	}
	return false
}

func makeFoolResp(val string) string {
	return fmt.Sprintf("**%s** doesn't satisfy me!", val)
}

func splitByCombs(inp string) []string {
	for _, val := range combs {
		if strings.Contains(inp, val) {
			return strings.Split(inp, val)
		}
	}
	return []string{inp}
}

func (b *EoD) getMessageElem(id string, guild string) (string, bool) {
	lock.RLock()
	dat, exists := b.dat[guild]
	lock.RUnlock()
	if !exists {
		return "Guild not setup yet!", false
	}
	el, res := dat.GetMsgElem(id)
	if !res.Exists {
		return res.Message, false
	}
	return el, true
}
