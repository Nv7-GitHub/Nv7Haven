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
	Mod  string //ID
	Text string
	Date int64
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
				Mod:  m.Author.ID,
				Text: message,
				Date: time.Now().Unix(),
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
		if !(len(m.Mentions) > 0) {
			s.ChannelMessageSend(m.ChannelID, "You need to mention the person you are going to warn!")
			return
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

		text := ""
		var warn warning
		for _, warnVal := range existing {
			mapstructure.Decode(warnVal, &warn)
			mem, err := s.GuildMember(m.GuildID, warn.Mod)
			if b.handle(err, m) {
				return
			}
			text += fmt.Sprintf("Warned by **%s** on **%s**\nWarning: **%s**\n\n", mem.Nick+"#"+mem.User.Discriminator, time.Unix(warn.Date, 0).Format("Jan 2 2006"), warn.Text)
		}
		mem, err := s.GuildMember(m.GuildID, m.Mentions[0].ID)
		if b.handle(err, m) {
			return
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("Warnings for **%s**", mem.Nick+"#"+mem.User.Discriminator),
			Description: text,
		})
		return
	}
}
