package discord

import (
	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "givenum",
			Description: "Give yourself a number!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "number",
					Description: "Number to give yourself",
					Required:    true,
				},
			},
		},
		{
			Name:        "getnum",
			Description: "Get your number or someone else's number!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "Optionally, provide a user to get the number of",
					Required:    false,
				},
			},
		},
		{
			Name:        "randselect",
			Description: "Select a random number out of all the numbers people have and congratulate them!",
		},
		{
			Name:        "meme",
			Description: "Get a meme fresh off of reddit! (r/memes)",
		},
		{
			Name:        "cmeme",
			Description: "Get a clean meme fresh off of reddit! (r/cleanmemes)",
		},
		{
			Name:        "pmeme",
			Description: "Get a programming meme fresh off of reddit! (r/ProgrammerHumor)",
		},
		{
			Name:        "warn",
			Description: "Warn a user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to warn",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "The warning's description",
					Required:    true,
				},
			},
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"givenum": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.giveNumCmd(int(i.Data.Options[0].IntValue()), bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"getnum": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			mention := ""
			if len(i.Data.Options) > 0 {
				mention = i.Data.Options[0].UserValue(bot.dg).ID
			}
			bot.getNumCmd(len(i.Data.Options) > 0, mention, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"randselect": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.randselectCmd(bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"meme": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.memeCommand(bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"cmeme": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.cmemeCommand(bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"pmeme": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.pmemeCommand(bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"warn": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.warnCmd(i.Data.Options[0].UserValue(s).ID, i.Member.User.ID, i.Data.Options[1].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
		},
	}
)
