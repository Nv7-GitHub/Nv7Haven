package eod

import (
	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "setvotes",
			Description: "Sets the vote count required in the server",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "votecount",
					Description: "The number of votes required for a poll to be completed.",
					Required:    true,
				},
			},
		},
		{
			Name:        "setplaychannel",
			Description: "Mark a channel as a play channel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "Channel to mark as a play channel",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "isplaychannel",
					Description: "Is it a play channel? If not given, defaults to true.",
					Required:    false,
				},
			},
		},
		{
			Name:        "setvotingchannel",
			Description: "Set a channel to be a channel for polls",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "Channel to set as a voting channel",
					Required:    true,
				},
			},
		},
		{
			Name:        "setnewschannel",
			Description: "Set a channel to be a channel for news",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "Channel to set as a news channel",
					Required:    true,
				},
			},
		},
		{
			Name:        "suggest",
			Description: "Create a suggestion!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "result",
					Description: "What the result for a combo should be",
					Required:    true,
				},
			},
		},
		{
			Name:        "mark",
			Description: "Suggest a creator mark, or add a creator mark to an element you created!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "element",
					Description: "The name of the element to add a creator mark to!",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "mark",
					Description: "What the new creator mark should be!",
					Required:    true,
				},
			},
		},
		{
			Name:        "image",
			Description: "Suggest an image for an element, or add an image to an element you created!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "element",
					Description: "The name of the element to add the image to!",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "imageurl",
					Description: "URL of an image to add to the element! You can also upload an image and then put the link here.",
					Required:    true,
				},
			},
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"setnewschannel": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.setNewsChannel(i.Data.Options[0].RoleValue(bot.dg, i.GuildID).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"setvotingchannel": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.setVotingChannel(i.Data.Options[0].RoleValue(bot.dg, i.GuildID).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"setvotes": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.setVoteCount(int(i.Data.Options[0].IntValue()), bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"setplaychannel": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			isPlayChannel := true
			if len(i.Data.Options) > 1 {
				isPlayChannel = i.Data.Options[1].BoolValue()
			}
			bot.setPlayChannel(i.Data.Options[0].RoleValue(bot.dg, i.GuildID).ID, isPlayChannel, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"suggest": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.suggestCmd(i.Data.Options[0].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"mark": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.markCmd(i.Data.Options[0].StringValue(), i.Data.Options[1].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"image": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.imageCmd(i.Data.Options[0].StringValue(), i.Data.Options[1].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
		},
	}
)
