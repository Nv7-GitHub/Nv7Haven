package discord

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) mod(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "warn") {
		if !(len(m.Mentions) > 0) {
			s.ChannelMessageSend(m.ChannelID, "You need to mention the person you are going to warn!")
			return
		}
		b.checkuserwithid(m, m.Mentions[0].ID)

		log.Println(m.Content)
	}
}
