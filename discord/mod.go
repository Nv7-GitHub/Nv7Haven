package discord

import (
	"fmt"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mitchellh/mapstructure"
)

var warnMatch = regexp.MustCompile(`warn <@!?\d+> (.+)`)

type empty struct{}

type warning struct {
	Mod   string //ID
	Text  string
	Date  int64
	Guild string //ID
}

func (b *Bot) warnCmd(user string, mod string, text string, m msg, rsp rsp) {
	if b.isUserMod(m, rsp, mod) {
		warn := warning{
			Mod:   m.Author.ID,
			Text:  text,
			Date:  time.Now().Unix(),
			Guild: m.GuildID,
		}

		serverData := b.readServerData(rsp, m.GuildID)
		_, exists := serverData["warns"]
		if !exists {
			serverData["warns"] = make(map[string]interface{})
		}
		var existing []interface{}
		_, exists = serverData["warns"].(map[string]interface{})[user]
		if !exists {
			existing = make([]interface{}, 0)
		} else {
			existing = serverData["warns"].(map[string]interface{})[user].([]interface{})
		}
		existing = append(existing, warn)
		serverData["warns"].(map[string]interface{})[user] = existing
		b.changeServerData(rsp, m.GuildID, serverData)
		rsp.Resp(`Successfully warned user.`)
		return
	}
	rsp.Resp(`You need to have permission "Administrator" to use this command.`)
	return
}

func (b *Bot) warnsCmd(hasMention bool, mentionID string, m msg, rsp rsp) {
	if !b.isUserMod(m, rsp, m.Author.ID) {
		rsp.ErrorMessage(`You need to have permission "Administrator" to use this command.`)
		return
	}

	users := make(map[string]interface{}, 0)
	serverData := b.readServerData(rsp, m.GuildID)
	_, exists := serverData["warns"]
	if !exists {
		serverData["warns"] = make(map[string]interface{})
	}
	warns := serverData["warns"].(map[string]interface{})
	if hasMention {
		_, exists := warns[mentionID]
		if exists {
			users[mentionID] = warns[mentionID]
		} else {
			rsp.ErrorMessage("That user does not have any warnings.")
			return
		}
	} else {
		users = warns
	}

	for userID, warnVals := range warns {
		var warn warning
		var text string
		for _, warning := range warnVals.([]interface{}) {
			err := mapstructure.Decode(warning, &warn)
			if rsp.Error(err) {
				return
			}
			user, err := b.dg.User(warn.Mod)
			if rsp.Error(err) {
				return
			}
			text += fmt.Sprintf("Warned by **%s** on **%s**\nWarning: **%s**\n\n", user.Username+"#"+user.Discriminator, time.Unix(warn.Date, 0).Format("Jan 2 2006"), warn.Text)
		}
		usr, err := b.dg.User(userID)
		if rsp.Error(err) {
			return
		}
		rsp.Embed(&discordgo.MessageEmbed{
			Title:       fmt.Sprintf("Warnings for **%s**", usr.Username+"#"+usr.Discriminator),
			Description: text,
		})
	}
}

func (b *Bot) mod(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if b.startsWith(m, "warn ") {
		if !(len(m.Mentions) > 0) {
			s.ChannelMessageSend(m.ChannelID, "You need to mention the person you are going to warn!")
			return
		}

		groups := warnMatch.FindAllStringSubmatch(m.Content, -1)
		if len(groups) < 1 {
			s.ChannelMessageSend(m.ChannelID, "Does not match format `warn @user <warning text>`")
			return
		}
		messageCont := groups[0]
		if len(messageCont) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Does not match format `warn @user <warning text>`")
			return
		}
		message := messageCont[1]
		b.warnCmd(m.Mentions[0].ID, m.Author.ID, message, b.newMsgNormal(m), b.newRespNormal(m))
		return
	}

	if b.startsWith(m, "warns") {
		hasMention := len(m.Mentions) > 0
		mention := ""
		if hasMention {
			mention = m.Mentions[0].ID
		}
		b.warnsCmd(hasMention, mention, b.newMsgNormal(m), b.newRespNormal(m))
	}

	if b.startsWith(m, "addrole") {
		var name string
		_, err := fmt.Sscanf(m.Content, "addrole %s", &name)
		if b.handle(err, m) {
			return
		}

		if !b.isMod(m, m.Author.ID) {
			s.ChannelMessageSend(m.ChannelID, "You need to have permission `Administrator` to use this command!")
			return
		}

		dat := b.getServerData(m, m.GuildID)
		_, exists := dat["roles"]
		if !exists {
			dat["roles"] = make(map[string]interface{})
		}
		dat["roles"].(map[string]interface{})[name] = empty{}
		b.updateServerData(m, m.GuildID, dat)

		role, err := s.GuildRoleCreate(m.GuildID)
		if b.handle(err, m) {
			return
		}
		s.GuildRoleEdit(m.GuildID, role.ID, name, role.Color, role.Hoist, role.Permissions, true)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully created role `%s`", name))
	}

	if b.startsWith(m, "rmrole") {
		var name string
		_, err := fmt.Sscanf(m.Content, "rmrole %s", &name)
		if b.handle(err, m) {
			return
		}

		if !b.isMod(m, m.Author.ID) {
			s.ChannelMessageSend(m.ChannelID, "You need to have permission `Administrator` to use this command!")
			return
		}

		roles, err := s.GuildRoles(m.GuildID)
		if b.handle(err, m) {
			return
		}

		id := ""
		for _, role := range roles {
			if role.Name == name {
				id = role.ID
			}
		}
		if id == "" {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role `%s` doesn't exist!", name))
			return
		}

		dat := b.getServerData(m, m.GuildID)
		_, exists := dat["roles"]
		if exists {
			delete(dat["roles"].(map[string]interface{}), name)
		}
		b.updateServerData(m, m.GuildID, dat)

		err = s.GuildRoleDelete(m.GuildID, id)
		if b.handle(err, m) {
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully removed role `%s`", name))
		return
	}

	if b.startsWith(m, "giverole") {
		var name string
		_, err := fmt.Sscanf(m.Content, "giverole %s", &name)
		if b.handle(err, m) {
			return
		}

		var role *discordgo.Role
		roles, err := s.GuildRoles(m.GuildID)
		if b.handle(err, m) {
			return
		}
		for _, rol := range roles {
			if rol.Name == name {
				role = rol
			}
		}
		if role == new(discordgo.Role) {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role `%s` doesn't exist!", name))
			return
		}

		dat := b.getServerData(m, m.GuildID)
		_, exists := dat["roles"]
		if !exists {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role `%s` hasn't been created by this bot!", name))
			return
		}
		_, exists = dat["roles"].(map[string]interface{})[name]
		if !exists {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role `%s` hasn't been created by this bot!", name))
			return
		}

		if role == nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role `%s` doesn't exist!", name))
			return
		}
		err = s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, role.ID)
		if b.handle(err, m) {
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully gave role `%s`", name))
		return
	}

	if b.startsWith(m, "norole") {
		var name string
		_, err := fmt.Sscanf(m.Content, "norole %s", &name)
		if b.handle(err, m) {
			return
		}

		var role *discordgo.Role
		roles, err := s.GuildRoles(m.GuildID)
		if b.handle(err, m) {
			return
		}
		for _, rol := range roles {
			if rol.Name == name {
				role = rol
			}
		}
		if role == new(discordgo.Role) {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role `%s` doesn't exist!", name))
			return
		}

		dat := b.getServerData(m, m.GuildID)
		_, exists := dat["roles"]
		if !exists {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role `%s` hasn't been created by this bot!", name))
			return
		}
		_, exists = dat["roles"].(map[string]interface{})[name]
		if !exists {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role `%s` hasn't been created by this bot!", name))
			return
		}

		if role == nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Role `%s` doesn't exist!", name))
			return
		}

		hasRole := false
		for _, rl := range m.Member.Roles {
			if rl == role.ID {
				hasRole = true
				break
			}
		}
		if !hasRole {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You don't have role `%s`!", name))
			return
		}
		err = s.GuildMemberRoleRemove(m.GuildID, m.Author.ID, role.ID)
		if b.handle(err, m) {
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully removed role `%s`", name))
		return
	}
}
