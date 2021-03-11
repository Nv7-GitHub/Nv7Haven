package eod

import "github.com/bwmarrin/discordgo"

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
