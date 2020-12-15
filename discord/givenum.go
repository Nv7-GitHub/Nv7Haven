package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) giveNum(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// givenum command
	if strings.HasPrefix(m.Content, "givenum") {
		var num int
		_, err := fmt.Sscanf(m.Content, "givenum %d", &num)
		if b.handle(err, m) {
			return
		}
		res, err := b.db.Query("SELECT COUNT(1) FROM givenum WHERE guild=? AND member=? LIMIT 1", m.GuildID, m.Author.ID)
		defer res.Close()
		if b.handle(err, m) {
			return
		}
		var count int
		if (num > 100) || (num < 0) {
			s.ChannelMessageSend(m.ChannelID, "You need to choose a number from 0-100.")
			return
		}
		res.Next()
		res.Scan(&count)
		if count == 0 {
			_, err = b.db.Exec("INSERT INTO givenum VALUES ( ?, ?, ? )", m.GuildID, m.Author.ID, num)
			if b.handle(err, m) {
				return
			}
			s.ChannelMessageSend(m.ChannelID, "Successfully initialized value.")
			return
		}
		_, err = b.db.Exec("UPDATE givenum SET number=? WHERE guild=? AND member=?", num, m.GuildID, m.Author.ID)
		if b.handle(err, m) {
			return
		}
		s.ChannelMessageSend(m.ChannelID, "Successfully updated value.")
		return
	}

	// getnum command
	if strings.HasPrefix(m.Content, "getnum") {
		res, err := b.db.Query("SELECT number FROM givenum WHERE guild=? AND member=? LIMIT 1", m.GuildID, m.Author.ID)
		defer res.Close()
		if b.handle(err, m) {
			return
		}
		res.Next()
		var num int
		err = res.Scan(&num)
		if b.handle(err, m) {
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Your number is %d.", num))
		return
	}

	// roles command
	if strings.HasPrefix(m.Content, "randselect") {
		mem, err := s.GuildMember(m.GuildID, m.Author.ID)
		if b.handle(err, m) {
			return
		}
		guildRoles, err := s.GuildRoles(m.GuildID)
		if b.handle(err, m) {
			return
		}
		for _, role := range mem.Roles {
			for _, grole := range guildRoles {
				if grole.ID == role {
					if strings.ToLower(grole.Name) == "admin" {
						s.ChannelMessageSend(m.ChannelID, "This has not been implemented yet.")
						return
					}
				}
			}
		}
		s.ChannelMessageSend(m.ChannelID, `You need to have a role called "admin".`)
		return
	}
}
