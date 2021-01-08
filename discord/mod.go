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
		b.checkuserwithid(m, m.Mentions[0].ID)

		log.Println(m.Content)
	}
}
