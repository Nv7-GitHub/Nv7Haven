package base

import (
	"errors"
	"fmt"
	"sort"

	"github.com/bwmarrin/discordgo"
)

func (b *Base) GetColor(guild, id string) (int, error) {
	db, res := b.GetDB(guild)
	if res.Exists {
		db.Config.RLock()
		col, exists := db.Config.UserColors[id]
		db.Config.RUnlock()
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
