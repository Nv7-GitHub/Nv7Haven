package discord

import (
	"encoding/json"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) giveNum(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "roles") {
		mem, err := s.GuildMember(m.GuildID, m.Mentions[0].ID)
		if b.handle(err, m) {
			return
		}
		roleNames := make([]string, len(mem.Roles))
		guildRoles, err := s.GuildRoles(m.GuildID)
		if b.handle(err, m) {
			return
		}
		for i, role := range mem.Roles {
			for _, grole := range guildRoles {
				if grole.ID == role {
					roleNames[i] = grole.Name
				}
			}
		}
		dat, err := json.Marshal(roleNames)
		if b.handle(err, m) {
			return
		}
		s.ChannelMessageSend(m.ChannelID, string(dat))
	}
}
