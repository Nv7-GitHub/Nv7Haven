package base

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func (b *Base) GetColor(guild, id string) (int, error) {
	b.lock.RLock()
	dat, exists := b.dat[guild]
	b.lock.RUnlock()
	if exists {
		col, exists := dat.UserColors[id]
		if exists {
			return col, nil
		}
	}

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
		role, err := b.GetRole(roleID, guild)
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

func (b *Base) GetRole(id string, guild string) (*discordgo.Role, error) {
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

func (b *Base) SaveInv(guild string, user string, newmade bool, recalculate ...bool) {
	b.lock.RLock()
	dat, exists := b.dat[guild]
	b.lock.RUnlock()
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
