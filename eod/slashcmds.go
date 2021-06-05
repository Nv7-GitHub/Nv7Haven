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
			Name:        "setpolls",
			Description: "Sets the maximum amount of polls a user can make",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "pollcount",
					Description: "The maximum number of polls a user can make",
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
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "autocapitalize",
					Description: "Should the bot autocapitalize? Default: true",
					Required:    false,
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
		{
			Name:        "inv",
			Description: "See your elements!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "Optionally, get the inventory of another user!",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "sortby",
					Description: "How to sort the inventory",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Name",
							Value: "name",
						},
						{
							Name:  "Element ID",
							Value: "id",
						},
						{
							Name:  "Made By",
							Value: "madeby",
						},
					},
				},
			},
		},
		{
			Name:        "lb",
			Description: "See the leaderboard!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "sortby",
					Description: "What to sort the leaderboard by!",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Elements Found",
							Value: "count",
						},
						{
							Name:  "Elements Made",
							Value: "made",
						},
					},
				},
			},
		},
		{
			Name:        "addcat",
			Description: "Suggest or add an element to a category!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "category",
					Description: "The name of the category to add the element to!",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "elem",
					Description: "What element to add!",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "elem2",
					Description: "Another element to add to the category!",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "elem3",
					Description: "Another element to add to the category!",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "elem4",
					Description: "Another element to add to the category!",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "elem5",
					Description: "Another element to add to the category!",
					Required:    false,
				},
			},
		},
		{
			Name:        "cat",
			Description: "Get info on a category!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "category",
					Description: "Name of the category",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "sort",
					Description: "How to sort the elements of the category!",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Alphabetical",
							Value: catSortAlphabetical,
						},
						{
							Name:  "Found",
							Value: catSortByFound,
						},
						{
							Name:  "Not Found",
							Value: catSortByNotFound,
						},
					},
				},
			},
		},
		{
			Name:        "hint",
			Description: "Get a hint on an element!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "element",
					Description: "Name of the element!",
					Required:    false,
				},
			},
		},
		{
			Name:        "stats",
			Description: "Get your server's stats!",
		},
		{
			Name:        "giveall",
			Description: "Give a user every element!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to give every element to!",
					Required:    true,
				},
			},
		},
		{
			Name:        "resetinv",
			Description: "Reset a user's inventory!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The user to reset the inventory of!",
					Required:    true,
				},
			},
		},
		{
			Name:        "give",
			Description: "Give a user an element, and choose whether to give all the elements required to make that element!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "element",
					Description: "Name of the element!",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "givetree",
					Description: "Give all the elements required to make that element?",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User to give the element (and maybe the elements required) to!",
					Required:    true,
				},
			},
		},
		{
			Name:        "path",
			Description: "Calculate the path of an element!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "element",
					Description: "Name of the element!",
					Required:    true,
				},
			},
		},
		{
			Name:        "elemsort",
			Description: "Sort all the elements in this server!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "sortby",
					Description: "How to sort the elements",
					Required:    true,
					Choices:     infoChoices,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "order",
					Description: "The order to sort the elements!",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Descending",
							Value: "0",
						},
						{
							Name:  "Ascending",
							Value: "1",
						},
					},
				},
			},
		},
		{
			Name:        "about",
			Description: "Get help and learn about the bot!",
		},
		{
			Name:        "setmodrole",
			Description: "Set a role to be a role for moderators!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role",
					Description: "Role to be set as moderator role",
					Required:    true,
				},
			},
		},
		{
			Name:        "rmcat",
			Description: "Suggest or remove an element from a category!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "category",
					Description: "The name of the category to add the element to!",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "elem",
					Description: "What element to remove!",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "elem2",
					Description: "Another element to remove from the category!",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "elem3",
					Description: "Another element to remove from the category!",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "elem4",
					Description: "Another element to remove from the category!",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "elem5",
					Description: "Another element to remove from the category!",
					Required:    false,
				},
			},
		},
		{
			Name:        "idea",
			Description: "Get a random unused combination!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "count",
					Description: "Number of random unused elements in the combination!",
					Required:    false,
				},
			},
		},
		{
			Name:        "catimg",
			Description: "Add an image to a category!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "category",
					Description: "The name of the category to add the image to!",
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
			bot.setNewsChannel(i.Data.Options[0].ChannelValue(bot.dg).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"setvotingchannel": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.setVotingChannel(i.Data.Options[0].ChannelValue(bot.dg).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"setvotes": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.setVoteCount(int(i.Data.Options[0].IntValue()), bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"setpolls": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.setPollCount(int(i.Data.Options[0].IntValue()), bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"setplaychannel": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			isPlayChannel := true
			if len(i.Data.Options) > 1 {
				isPlayChannel = i.Data.Options[1].BoolValue()
			}
			bot.setPlayChannel(i.Data.Options[0].ChannelValue(bot.dg).ID, isPlayChannel, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"suggest": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			autocapitalize := true
			if len(i.Data.Options) > 1 {
				autocapitalize = i.Data.Options[1].BoolValue()
			}
			bot.suggestCmd(i.Data.Options[0].StringValue(), autocapitalize, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"mark": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.markCmd(i.Data.Options[0].StringValue(), i.Data.Options[1].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"image": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.imageCmd(i.Data.Options[0].StringValue(), i.Data.Options[1].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"inv": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			sortby := "name"
			id := i.Member.User.ID
			for _, val := range i.Data.Options {
				if val.Name == "sortby" {
					sortby = val.StringValue()
				}

				if val.Name == "user" {
					id = val.UserValue(bot.dg).ID
				}
			}
			bot.invCmd(id, bot.newMsgSlash(i), bot.newRespSlash(i), sortby)
		},
		"lb": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			sort := "count"
			if len(i.Data.Options) > 0 {
				sort = i.Data.Options[0].StringValue()
			}
			bot.lbCmd(bot.newMsgSlash(i), bot.newRespSlash(i), sort)
		},
		"addcat": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			suggestAdd := []string{i.Data.Options[1].StringValue()}
			if len(i.Data.Options) > 2 {
				for _, val := range i.Data.Options[2:] {
					suggestAdd = append(suggestAdd, val.StringValue())
				}
			}
			bot.categoryCmd(suggestAdd, i.Data.Options[0].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"cat": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if len(i.Data.Options) == 0 {
				bot.allCatCmd(bot.newMsgSlash(i), bot.newRespSlash(i))
				return
			}
			if len(i.Data.Options) == 1 {
				rsp := bot.newRespSlash(i)
				if i.Data.Options[0].Name == "sort" {
					rsp.ErrorMessage("You must specify a category name!")
					return
				}
				bot.catCmd(i.Data.Options[0].StringValue(), catSortAlphabetical, bot.newMsgSlash(i), rsp)
				return
			}

			bot.catCmd(i.Data.Options[0].StringValue(), int(i.Data.Options[1].IntValue()), bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"hint": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if len(i.Data.Options) == 0 {
				bot.hintCmd("", false, bot.newMsgSlash(i), bot.newRespSlash(i))
				return
			}
			bot.hintCmd(i.Data.Options[0].StringValue(), true, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"stats": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.statsCmd(bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"giveall": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.giveAllCmd(i.Data.Options[0].UserValue(bot.dg).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"resetinv": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.resetInvCmd(i.Data.Options[0].UserValue(bot.dg).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"give": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.giveCmd(i.Data.Options[0].StringValue(), i.Data.Options[1].BoolValue(), i.Data.Options[2].UserValue(bot.dg).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"path": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.calcTreeCmd(i.Data.Options[0].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"elemsort": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.sortCmd(i.Data.Options[0].StringValue(), i.Data.Options[1].StringValue() == "1", bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"about": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.aboutCmd(bot.newRespSlash(i))
		},
		"setmodrole": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.setModRole(i.Data.Options[0].RoleValue(bot.dg, i.GuildID).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"rmcat": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			suggestAdd := []string{i.Data.Options[1].StringValue()}
			if len(i.Data.Options) > 2 {
				for _, val := range i.Data.Options[2:] {
					suggestAdd = append(suggestAdd, val.StringValue())
				}
			}
			bot.rmCategoryCmd(suggestAdd, i.Data.Options[0].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"idea": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			count := 2
			if len(i.Data.Options) > 0 {
				count = int(i.Data.Options[0].IntValue())
			}
			bot.ideaCmd(count, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"catimg": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.catImgCmd(i.Data.Options[0].StringValue(), i.Data.Options[1].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
		},
	}
)
