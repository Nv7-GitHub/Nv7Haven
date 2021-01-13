package discord

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mitchellh/mapstructure"
)

var warnMatch = regexp.MustCompile(`warn <@!?\d+> (.+)`)

type warning struct {
	Mod   string //ID
	Text  string
	Date  int64
	Guild string //ID
}

func (b *Bot) mod(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "warn ") {
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

		b.checkuserwithid(m, m.Mentions[0].ID)

		if b.isMod(m, m.Author.ID) {
			warn := warning{
				Mod:   m.Author.ID,
				Text:  message,
				Date:  time.Now().Unix(),
				Guild: m.GuildID,
			}
			user, suc := b.getuser(m, m.Mentions[0].ID)
			if !suc {
				return
			}
			var existing []interface{}
			_, exists := user.Metadata["warns"]
			if !exists {
				existing = make([]interface{}, 0)
			} else {
				existing = user.Metadata["warns"].([]interface{})
			}
			existing = append(existing, warn)
			user.Metadata["warns"] = existing
			suc = b.updateuser(m, user)
			if !suc {
				return
			}
			s.ChannelMessageSend(m.ChannelID, `Successfully warned user.`)
			return
		}
		s.ChannelMessageSend(m.ChannelID, `You need to have permission "Administrator" to use this command.`)
		return
	}

	if strings.HasPrefix(m.Content, "warns") {
		if !b.isMod(m, m.Author.ID) {
			s.ChannelMessageSend(m.ChannelID, `You need to have permission "Administrator" to use this command.`)
			return
		}

		users := make([]string, 0)
		if !(len(m.Mentions) > 0) {
			res, err := b.db.Query("SELECT user FROM currency WHERE guilds LIKE ?", "%"+m.GuildID+"%")
			if b.handle(err, m) {
				return
			}
			defer res.Close()
			for res.Next() {
				var user string
				err = res.Scan(&user)
				if b.handle(err, m) {
					return
				}
				users = append(users, user)
			}
		} else {
			users = []string{m.Mentions[0].ID}
			b.checkuserwithid(m, m.Mentions[0].ID)
		}

		text := ""
		for _, userID := range users {
			text = ""
			user, suc := b.getuser(m, userID)
			if !suc {
				return
			}

			var existing = make([]interface{}, 0)
			_, exists := user.Metadata["warns"]
			if !exists {
				existing = make([]interface{}, 0)
			} else {
				existing = user.Metadata["warns"].([]interface{})
			}

			var warn warning
			warnCount := 0
			for _, warnVal := range existing {
				mapstructure.Decode(warnVal, &warn)
				if warn.Guild == m.GuildID {
					warnCount++
					user, err := s.User(warn.Mod)
					if b.handle(err, m) {
						return
					}
					text += fmt.Sprintf("Warned by **%s** on **%s**\nWarning: **%s**\n\n", user.Username+"#"+user.Discriminator, time.Unix(warn.Date, 0).Format("Jan 2 2006"), warn.Text)
				}
			}
			usr, err := s.User(userID)
			if b.handle(err, m) {
				return
			}
			if !(len(users) > 1 && warnCount == 0) {
				s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
					Title:       fmt.Sprintf("Warnings for **%s**", usr.Username+"#"+usr.Discriminator),
					Description: text,
				})
			}
		}
		return
	}

	if strings.HasPrefix(m.Content, "addrole") {
		var name string
		_, err := fmt.Sscanf(m.Content, "addrole %s", &name)
		if b.handle(err, m) {
			return
		}

		if !b.isMod(m, m.Author.ID) {
			s.ChannelMessageSend(m.ChannelID, "You need to have permission `Administrator` to use this command!")
			return
		}

		role, err := s.GuildRoleCreate(m.GuildID)
		if b.handle(err, m) {
			return
		}
		s.GuildRoleEdit(m.GuildID, role.ID, name, role.Color, role.Hoist, role.Permissions, true)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully created role %s", name))
	}
}
