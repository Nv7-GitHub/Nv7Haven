package eod

import (
	"fmt"
	"strconv"

	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "set",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Updates server data!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "votes",
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
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "polls",
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
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "playchannel",
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
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "votingchannel",
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
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "newschannel",
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
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "modrole",
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
			},
		},
		{
			Name:        "suggest",
			Type:        discordgo.ChatApplicationCommand,
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
			Type:        discordgo.ChatApplicationCommand,
			Description: "Suggest a mark, or add a mark to an element you created!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "element",
					Description: "The name of the element to add a mark to!",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "mark",
					Description: "What the new mark should be!",
					Required:    true,
				},
			},
		},
		{
			Name:        "image",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Add an image to an element or a category!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "element",
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
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "category",
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
							Description: "URL of an image to add to the category! You can also upload an image and then put the link here.",
							Required:    true,
						},
					},
				},
			},
		},
		{
			Name:        "inv",
			Type:        discordgo.ChatApplicationCommand,
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
					Description: "How to sort the inventory!",
					Required:    false,
					Choices:     util.SortChoices,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "filter",
					Description: "How to filter the inventory!",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "None",
							Value: "none",
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
			Type:        discordgo.ChatApplicationCommand,
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
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User to view the leaderboard from the POV of!",
					Required:    false,
				},
			},
		},
		{
			Name:        "addcat",
			Type:        discordgo.ChatApplicationCommand,
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
			Type:        discordgo.ChatApplicationCommand,
			Description: "Get info on a category!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "category",
					Description: "Name of the category",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "sort",
					Description: "How to sort the elements of the category!",
					Required:    false,
					Choices: append([]*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Found",
							Value: "catfound",
						},
						{
							Name:  "Not Found",
							Value: "catnotfound",
						},
						{
							Name:  "Element Count",
							Value: "catelemcount",
						},
					}, util.SortChoices...),
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User's inventory to compare",
					Required:    false,
				},
			},
		},
		{
			Name:        "hint",
			Type:        discordgo.ChatApplicationCommand,
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
			Type:        discordgo.ChatApplicationCommand,
			Description: "Get your server's stats!",
		},
		{
			Name:        "resetinv",
			Type:        discordgo.ChatApplicationCommand,
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
			Type:        discordgo.ChatApplicationCommand,
			Description: "Give elements to a user!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "element",
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
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "cat",
					Description: "Give a user all the elements in a category, and optionally give the tree!",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "category",
							Description: "Name of the category!",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionBoolean,
							Name:        "givetree",
							Description: "Give all the elements required to make the elements?",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionUser,
							Name:        "user",
							Description: "User to give the elemenst (and maybe the elements required) to!",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "all",
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
			},
		},
		{
			Name:        "path",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Calculate paths!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "element",
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
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "category",
					Description: "Calculate the path of a category!",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "category",
							Description: "Name of the category!",
							Required:    true,
						},
					},
				},
			},
		},
		{
			Name:        "elemsort",
			Type:        discordgo.ChatApplicationCommand,
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
			Name:        "help",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Get help and learn about the bot!",
		},
		{
			Name:        "rmcat",
			Type:        discordgo.ChatApplicationCommand,
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
			Type:        discordgo.ChatApplicationCommand,
			Description: "Get a random unused combination!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "count",
					Description: "Number of random unused elements in the combination!",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "category",
					Description: "Use a category for the elements to choose from!",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "element",
					Description: "Require an element to be in the idea!",
					Required:    false,
				},
			},
		},
		{
			Name:        "download",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Download an inventory or category!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "inv",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "Download a user's inventory!",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionUser,
							Name:        "user",
							Description: "Optionally, download the inventory of another user!",
							Required:    false,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "sortby",
							Description: "How to sort the inventory!",
							Required:    false,
							Choices:     util.SortChoices,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "filter",
							Description: "How to filter the inventory!",
							Required:    false,
							Choices: []*discordgo.ApplicationCommandOptionChoice{
								{
									Name:  "None",
									Value: "none",
								},
								{
									Name:  "Made By",
									Value: "madeby",
								},
							},
						},
						{
							Type:        discordgo.ApplicationCommandOptionBoolean,
							Name:        "postfix",
							Description: "Whether to put the value after the element!",
							Required:    false,
						},
					},
				},
				{
					Name:        "cat",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "Download a category!",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "category",
							Description: "Which category to download!",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "sortby",
							Description: "How to sort the category!",
							Required:    false,
							Choices:     util.SortChoices,
						},
						{
							Type:        discordgo.ApplicationCommandOptionBoolean,
							Name:        "postfix",
							Description: "Whether to put the value after the element!",
							Required:    false,
						},
					},
				},
			},
		},
		{
			Name:        "breakdown",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Breakdown elements!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "element",
					Description: "Get an element's breakdown!",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "element",
							Description: "Name of the element!",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionBoolean,
							Name:        "calctree",
							Description: "Whether to include the tree of that element!",
							Required:    false,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "category",
					Description: "Get the breakdown of a category!",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "category",
							Description: "Name of the category!",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionBoolean,
							Name:        "calctree",
							Description: "Whether to include the tree of the elements in the category!",
							Required:    false,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "inv",
					Description: "Get the breakdown of a user's inventory!",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionUser,
							Name:        "user",
							Description: "Which user's inventory to breakdown!",
							Required:    false,
						},
						{
							Type:        discordgo.ApplicationCommandOptionBoolean,
							Name:        "calctree",
							Description: "Whether to include the tree of the elements in the category!",
							Required:    false,
						},
					},
				},
			},
		},
		{
			Name:        "graph",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Graph element trees!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "element",
					Description: "Create a graph of an element's tree!",
					Options: append([]*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "element",
							Description: "Name of the element!",
							Required:    true,
						},
					}, trees.GraphOpts...),
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "category",
					Description: "Create a graph of an element's tree!",
					Options: append([]*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "category",
							Description: "Name of the category!",
							Required:    true,
						},
					}, trees.GraphOpts...),
				},
			},
		},
		{
			Name:        "get",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Get a value of an element!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "found",
					Description: "See the user's who have found an element!",
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
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "categories",
					Description: "See the categories an element is in!",
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
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "info",
					Description: "Get the info of an element!",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "element",
							Description: "Name of the element!",
							Required:    false,
						},
						{
							Type:        discordgo.ApplicationCommandOptionInteger,
							Name:        "id",
							Description: "ID of the element!",
							Required:    false,
						},
					},
				},
			},
		},
		{
			Name:        "setcolor",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Set your embed color! If you don't provide a color, it will reset your color.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "color",
					Description: "Hex code to set your embed color too",
					Required:    false,
				},
			},
		},
		{
			Name:        "invhint",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Get the inverse hint of an element!",
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
			Name:        "elemsearch",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Search for an element by name!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "query",
					Description: "The query to search with!",
					Required:    true,
				},
			},
		},
		{
			Name: "View Inventory",
			//Description: "View the user's inventory!",
			Type: discordgo.UserApplicationCommand,
		},
		{
			Name: "View Info",
			//Description: "View the info of the element in a message!",
			Type: discordgo.MessageApplicationCommand,
		},
		{
			Name: "Get Hint",
			//Description: "Get the hint of the element in a message!",
			Type: discordgo.MessageApplicationCommand,
		},
		{
			Name: "Get Inverse Hint",
			//Description: "Get the inverse hint of the element in a message!",
			Type: discordgo.MessageApplicationCommand,
		},
		{
			Name: "Get Color",
			//Description: "Get a user's embed color!",
			Type: discordgo.UserApplicationCommand,
		},
		{
			Name: "View Leaderboard",
			//Description: "View the leaderboard from the user's point of view!",
			Type: discordgo.UserApplicationCommand,
		},
		{
			Name: "View Inventory Breakdown",
			//Description: "View a user's inventory breakdown!",
			Type: discordgo.UserApplicationCommand,
		},
		{
			Name:        "notation",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Calculate notations!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "element",
					Description: "Calculate the notation of an element!",
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
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "category",
					Description: "Calculate the notation of a category!",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "category",
							Description: "Name of the category!",
							Required:    true,
						},
					},
				},
			},
		},
		{
			Name:        "color",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Set the color of an element or category!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "element",
					Description: "Suggest the color for an element, or set the color of an element you created!",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "element",
							Description: "The name of the element to set the color of!",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "color",
							Description: "The new hex color of the element.",
							Required:    true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "category",
					Description: "Add an image to a category!",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "category",
							Description: "The name of the category to set the color of!",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "color",
							Description: "The new hex color of the element.",
							Required:    false,
						},
					},
				},
			},
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"set": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData().Options[0]
			switch resp.Name {
			case "newschannel":
				bot.setNewsChannel(resp.Options[0].ChannelValue(bot.dg).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
			case "votingchannel":
				bot.setVotingChannel(resp.Options[0].ChannelValue(bot.dg).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
			case "votes":
				bot.setVoteCount(int(resp.Options[0].IntValue()), bot.newMsgSlash(i), bot.newRespSlash(i))
			case "polls":
				bot.setPollCount(int(resp.Options[0].IntValue()), bot.newMsgSlash(i), bot.newRespSlash(i))
			case "modrole":
				bot.setModRole(resp.Options[0].RoleValue(bot.dg, i.GuildID).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
			case "playchannel":
				isPlayChannel := true
				if len(resp.Options) > 1 {
					isPlayChannel = resp.Options[1].BoolValue()
				}
				bot.setPlayChannel(resp.Options[0].ChannelValue(bot.dg).ID, isPlayChannel, bot.newMsgSlash(i), bot.newRespSlash(i))
			}
		},
		"suggest": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			autocapitalize := true
			if len(resp.Options) > 1 {
				autocapitalize = resp.Options[1].BoolValue()
			}
			bot.suggestCmd(resp.Options[0].StringValue(), autocapitalize, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"mark": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			bot.markCmd(resp.Options[0].StringValue(), resp.Options[1].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"image": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData().Options[0]
			switch resp.Name {
			case "element":
				bot.imageCmd(resp.Options[0].StringValue(), resp.Options[1].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))

			case "category":
				bot.catImgCmd(resp.Options[0].StringValue(), resp.Options[1].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
			}
		},
		"inv": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			sortby := "name"
			filter := "none"
			id := i.Member.User.ID
			for _, val := range resp.Options {
				if val.Name == "sortby" {
					sortby = val.StringValue()
				}

				if val.Name == "filter" {
					filter = val.StringValue()
				}

				if val.Name == "user" {
					id = val.UserValue(bot.dg).ID
				}
			}
			bot.invCmd(id, bot.newMsgSlash(i), bot.newRespSlash(i), sortby, filter)
		},
		"lb": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			sort := "count"
			user := i.Member.User.ID
			for _, opt := range resp.Options {
				if opt.Name == "sortby" {
					sort = resp.Options[0].StringValue()
				}

				if opt.Name == "user" {
					user = opt.UserValue(bot.dg).ID
				}
			}
			bot.lbCmd(bot.newMsgSlash(i), bot.newRespSlash(i), sort, user)
		},
		"addcat": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			suggestAdd := []string{resp.Options[1].StringValue()}
			if len(resp.Options) > 2 {
				for _, val := range resp.Options[2:] {
					suggestAdd = append(suggestAdd, val.StringValue())
				}
			}
			bot.categoryCmd(suggestAdd, resp.Options[0].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"cat": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			isAll := true
			sort := "name"
			catName := ""
			hasUser := false
			var user string
			for _, val := range resp.Options {
				if val.Name == "category" {
					isAll = false
					catName = val.StringValue()
				}

				if val.Name == "sort" {
					sort = val.StringValue()
				}

				if val.Name == "user" {
					hasUser = true
					user = val.UserValue(bot.dg).ID
				}
			}

			if isAll {
				bot.allCatCmd(sort, hasUser, user, bot.newMsgSlash(i), bot.newRespSlash(i))
				return
			}

			bot.catCmd(catName, sort, hasUser, user, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"hint": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			hasElem := false
			var elem string
			for _, opt := range resp.Options {
				if opt.Name == "element" {
					hasElem = true
					elem = opt.StringValue()
				}
			}

			bot.hintCmd(elem, hasElem, false, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"stats": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.statsCmd(bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"resetinv": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			bot.resetInvCmd(resp.Options[0].UserValue(bot.dg).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"give": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData().Options[0]
			switch resp.Name {
			case "element":
				bot.giveCmd(resp.Options[0].StringValue(), resp.Options[1].BoolValue(), resp.Options[2].UserValue(bot.dg).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
			case "cat":
				bot.giveCatCmd(resp.Options[0].StringValue(), resp.Options[1].BoolValue(), resp.Options[2].UserValue(bot.dg).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
			case "all":
				bot.giveAllCmd(resp.Options[0].UserValue(bot.dg).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
			}

		},
		"path": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData().Options[0]
			switch resp.Name {
			case "element":
				bot.calcTreeCmd(resp.Options[0].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))

			case "category":
				bot.calcTreeCatCmd(resp.Options[0].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
			}
		},
		"elemsort": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			bot.sortCmd(resp.Options[0].StringValue(), resp.Options[1].StringValue() == "1", bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"help": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.helpCmd(bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"rmcat": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			suggestAdd := []string{resp.Options[1].StringValue()}
			if len(resp.Options) > 2 {
				for _, val := range resp.Options[2:] {
					suggestAdd = append(suggestAdd, val.StringValue())
				}
			}
			bot.rmCategoryCmd(suggestAdd, resp.Options[0].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"idea": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			count := 2
			hasCat := false
			catName := ""
			hasEl := false
			elName := ""
			for _, opt := range resp.Options {
				if opt.Name == "count" {
					count = int(opt.IntValue())
				}

				if opt.Name == "category" {
					hasCat = true
					catName = opt.StringValue()
				}

				if opt.Name == "element" {
					hasEl = true
					elName = opt.StringValue()
				}
			}
			bot.ideaCmd(count, catName, hasCat, elName, hasEl, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"download": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()

			switch resp.Options[0].Name {
			case "inv":
				opts := resp.Options[0]
				sortby := "name"
				filter := "none"
				id := i.Member.User.ID
				postfix := false
				for _, val := range opts.Options {
					if val.Name == "sortby" {
						sortby = val.StringValue()
					}

					if val.Name == "filter" {
						filter = val.StringValue()
					}

					if val.Name == "user" {
						id = val.UserValue(bot.dg).ID
					}

					if val.Name == "postfix" {
						postfix = val.BoolValue()
					}
				}
				bot.downloadInvCmd(id, sortby, filter, postfix, bot.newMsgSlash(i), bot.newRespSlash(i))
				return

			case "cat":
				opts := resp.Options[0]
				sortby := "name"
				catName := ""
				postfix := false
				for _, val := range opts.Options {
					if val.Name == "category" {
						catName = val.StringValue()
					}

					if val.Name == "sortby" {
						sortby = val.StringValue()
					}

					if val.Name == "postfix" {
						postfix = val.BoolValue()
					}
				}
				bot.downloadCatCmd(catName, sortby, postfix, bot.newMsgSlash(i), bot.newRespSlash(i))
				return
			}
		},
		"breakdown": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData().Options[0]
			switch resp.Name {
			case "element":
				calctree := true
				if len(resp.Options) > 1 {
					calctree = resp.Options[1].BoolValue()
				}
				bot.elemBreakdownCmd(resp.Options[0].StringValue(), calctree, bot.newMsgSlash(i), bot.newRespSlash(i))

			case "category":
				calctree := false
				if len(resp.Options) > 1 {
					calctree = resp.Options[1].BoolValue()
				}
				bot.catBreakdownCmd(resp.Options[0].StringValue(), calctree, bot.newMsgSlash(i), bot.newRespSlash(i))

			case "inv":
				user := i.Member.User.ID
				calcTree := true
				for _, opt := range resp.Options {
					if opt.Name == "user" {
						user = opt.UserValue(bot.dg).ID
					}

					if opt.Name == "calctree" {
						calcTree = opt.BoolValue()
					}
				}
				bot.invBreakdownCmd(user, calcTree, bot.newMsgSlash(i), bot.newRespSlash(i))
			}
		},
		"graph": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData().Options[0]

			switch resp.Name {
			case "element":
				var elem string
				outputType := ""
				layout := ""
				special := false
				for _, opt := range resp.Options {
					if opt.Name == "element" {
						elem = opt.StringValue()
					}

					if opt.Name == "output_type" {
						outputType = opt.StringValue()
					}

					if opt.Name == "layout" {
						layout = opt.StringValue()
					}

					if opt.Name == "distinct" {
						special = opt.BoolValue()
					}
				}
				bot.elemGraphCmd(elem, layout, outputType, special, bot.newMsgSlash(i), bot.newRespSlash(i))

			case "category":
				var catName string
				outputType := ""
				layout := ""
				special := false
				for _, opt := range resp.Options {
					if opt.Name == "category" {
						catName = opt.StringValue()
					}

					if opt.Name == "output_type" {
						outputType = opt.StringValue()
					}

					if opt.Name == "layout" {
						layout = opt.StringValue()
					}

					if opt.Name == "distinct" {
						special = opt.BoolValue()
					}
				}
				bot.catGraphCmd(catName, layout, outputType, special, bot.newMsgSlash(i), bot.newRespSlash(i))
			}

		},
		"get": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData().Options[0]
			switch resp.Name {
			case "found":
				bot.foundCmd(resp.Options[0].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))

			case "categories":
				bot.categoriesCmd(resp.Options[0].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))

			case "info":
				elem := ""
				var id int
				isID := false
				for _, opt := range resp.Options {
					if opt.Name == "element" {
						elem = opt.StringValue()
					}

					if opt.Name == "id" {
						isID = true
						id = int(opt.IntValue())
					}
				}
				rsp := bot.newRespSlash(i)
				if !isID && elem == "" {
					rsp.ErrorMessage("You must input an element or an element's ID!")
					return
				}
				if isID && elem != "" {
					rsp.ErrorMessage("You can't input an element and an element's ID!")
					return
				}
				bot.info(elem, id, isID, bot.newMsgSlash(i), rsp)
			}
		},
		"setcolor": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			color := ""
			rmColor := true
			for _, opt := range resp.Options {
				if opt.Name == "color" {
					rmColor = false
					color = opt.StringValue()
				}
			}
			bot.setUserColor(color, rmColor, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"invhint": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			bot.hintCmd(resp.Options[0].StringValue(), true, true, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"elemsearch": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			bot.elemSearchCmd(resp.Options[0].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"View Inventory": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			bot.invCmd(resp.TargetID, bot.newMsgSlash(i), bot.newRespSlash(i), "name", "none")
		},
		"View Info": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			rsp := bot.newRespSlash(i)
			res, suc := bot.getMessageElem(resp.TargetID, i.GuildID)
			if !suc {
				rsp.ErrorMessage(res)
				return
			}
			bot.infoCmd(res, bot.newMsgSlash(i), rsp)
		},
		"Get Hint": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			rsp := bot.newRespSlash(i)
			res, suc := bot.getMessageElem(resp.TargetID, i.GuildID)
			if !suc {
				rsp.ErrorMessage(res)
				return
			}
			bot.hintCmd(res, true, false, bot.newMsgSlash(i), rsp)
		},
		"Get Inverse Hint": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			rsp := bot.newRespSlash(i)
			res, suc := bot.getMessageElem(resp.TargetID, i.GuildID)
			if !suc {
				rsp.ErrorMessage(res)
				return
			}
			bot.hintCmd(res, true, true, bot.newMsgSlash(i), rsp)
		},
		"Get Color": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			rsp := bot.newRespSlash(i)
			color, err := bot.getColor(i.GuildID, resp.TargetID)
			if rsp.Error(err) {
				return
			}
			hex := strconv.FormatInt(int64(color), 16)
			rsp.Message(fmt.Sprintf("https://singlecolorimage.com/get/%s/100x100", hex))
		},
		"View Leaderboard": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			bot.lbCmd(bot.newMsgSlash(i), bot.newRespSlash(i), "count", resp.TargetID)
		},
		"notation": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData().Options[0]
			switch resp.Name {
			case "element":
				bot.notationCmd(resp.Options[0].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))

			case "category":
				bot.catNotationCmd(resp.Options[0].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
			}
		},
		"View Inventory Breakdown": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.invBreakdownCmd(i.ApplicationCommandData().TargetID, false, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"color": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData().Options[0]
			rsp := bot.newRespSlash(i)
			switch resp.Name {
			case "element":
				color, err := strconv.ParseInt(resp.Options[1].StringValue(), 16, 64)
				if rsp.Error(err) {
					return
				}
				bot.colorCmd(resp.Options[0].StringValue(), int(color), bot.newMsgSlash(i), rsp)

			case "category":
				var color int64 = 0
				if len(resp.Options) > 1 {
					var err error
					color, err = strconv.ParseInt(resp.Options[1].StringValue(), 16, 64)
					if rsp.Error(err) {
						return
					}
				}
				bot.catColorCmd(resp.Options[0].StringValue(), int(color), bot.newMsgSlash(i), bot.newRespSlash(i))
			}
		},
	}
)
