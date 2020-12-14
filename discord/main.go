package discord

import "github.com/bwmarrin/discordgo"

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

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
