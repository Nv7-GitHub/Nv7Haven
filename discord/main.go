package discord

import (
	"github.com/bwmarrin/discordgo"
)

const token = "Nzg4MTg1MzY1NTMzNTU2NzM2.X9f00g.krA6cjfFWYdzbqOPXq8NvRjxb3k"

// Bot is a discord bot
type Bot struct {
	dg *discordgo.Session
}

func (b *Bot) handlers() {
	b.dg.AddHandler(b.giveNum)
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

func (b *Bot) handle(err error, m *discordgo.MessageCreate) bool {
	if err != nil {
		b.dg.ChannelMessageSend(m.ChannelID, err.Error())
		return true
	}
	return false
}

// Close cleans up
func (b *Bot) Close() {
	b.dg.Close()
}
