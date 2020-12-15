package discord

import (
	"encoding/json"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const token = "Nzg4MTg1MzY1NTMzNTU2NzM2.X9f00g.krA6cjfFWYdzbqOPXq8NvRjxb3k"

// Bot is a discord bot
type Bot struct {
	dg *discordgo.Session
}

func (b *Bot) handlers() {
	b.dg.AddHandler(messageCreate)
}

// InitDiscord creates a discord bot
func InitDiscord() Bot {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	b := Bot{
		dg: dg,
	}
	b.handlers()
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	err = dg.Open()
	if err != nil {
		panic(err)
	}
	return b
}

// Close cleans up
func (b *Bot) Close() {
	b.dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "roles") {
		dat, _ := json.Marshal(m.Member.Roles)
		s.ChannelMessageSend(m.ChannelID, string(dat))
	}
}
