package eod

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/eodsort"
	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
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
							ChannelTypes: []discordgo.ChannelType{
								discordgo.ChannelTypeGuildText,
							},
							Required: true,
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
							ChannelTypes: []discordgo.ChannelType{
								discordgo.ChannelTypeGuildText,
							},
							Required: true,
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
							ChannelTypes: []discordgo.ChannelType{
								discordgo.ChannelTypeGuildText,
							},
							Required: true,
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
					Choices:     eodsort.SortChoices,
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
					Description: "Whether to postfix!",
					Required:    false,
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
						{
							Name:  "Elements Signed",
							Value: "signed",
						},
						{
							Name:  "Elements Imaged",
							Value: "imaged",
						},
						{
							Name:  "Elements Colored",
							Value: "colored",
						},
						{
							Name:  "Categories Imaged",
							Value: "catimaged",
						},
						{
							Name:  "Categories Colored",
							Value: "catcolored",
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
			Name:        "lbimage",
			Type:        discordgo.ChatApplicationCommand,
			Description: "See the leaderboard, as a chart!",
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
			Type:        discordgo.ChatApplicationCommand,
			Description: "Suggest or add an element to a category!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "category",
					Description:  "The name of the category to add the element to!",
					Required:     true,
					Autocomplete: true,
				},
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "elem",
					Description:  "What element to add!",
					Required:     true,
					Autocomplete: true,
				},
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "elem2",
					Description:  "Another element to add to the category!",
					Required:     false,
					Autocomplete: true,
				},
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "elem3",
					Description:  "Another element to add to the category!",
					Required:     false,
					Autocomplete: true,
				},
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "elem4",
					Description:  "Another element to add to the category!",
					Required:     false,
					Autocomplete: true,
				},
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "elem5",
					Description:  "Another element to add to the category!",
					Required:     false,
					Autocomplete: true,
				},
			},
		},
		{
			Name:        "cat",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Get info on a category!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "category",
					Description:  "Name of the category",
					Required:     false,
					Autocomplete: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "sort",
					Description: "How to sort the elements of the category!",
					Required:    false,
					Choices: append([]*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Element Count",
							Value: "catelemcount",
						},
					}, eodsort.SortChoices...),
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User's inventory to compare",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "postfix",
					Description: "Whether to postfix!",
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
					Choices:     eodsort.SortChoices,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "postfix",
					Description: "Whether to postfix or not?",
					Required:    false,
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
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "category",
					Description:  "The name of the category to add the element to!",
					Required:     true,
					Autocomplete: true,
				},
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "elem",
					Description:  "What element to remove!",
					Required:     true,
					Autocomplete: true,
				},
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "elem2",
					Description:  "Another element to remove from the category!",
					Required:     false,
					Autocomplete: true,
				},
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "elem3",
					Description:  "Another element to remove from the category!",
					Required:     false,
					Autocomplete: true,
				},
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "elem4",
					Description:  "Another element to remove from the category!",
					Required:     false,
					Autocomplete: true,
				},
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "elem5",
					Description:  "Another element to remove from the category!",
					Required:     false,
					Autocomplete: true,
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
			Name:        "ai_idea",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Get an AI generated combination!",
			Options:     []*discordgo.ApplicationCommandOption{},
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
							Choices:     eodsort.SortChoices,
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
							Choices:     eodsort.SortChoices,
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
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "element",
							Description:  "Name of the element!",
							Required:     true,
							Autocomplete: true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "categories",
					Description: "See the categories an element is in!",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "element",
							Description:  "Name of the element!",
							Required:     true,
							Autocomplete: true,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "info",
					Description: "Get the info of an element!",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "element",
							Description:  "Name of the element!",
							Autocomplete: true,
							Required:     false,
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
			Name:        "search",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Search for an element by name!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "elements",
					Description: "Search for an element by name!",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "query",
							Description: "The query to search with!",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "sort",
							Description: "How to sort the results!",
							Choices:     eodsort.SortChoices,
							Required:    false,
						},
						{
							Type:        discordgo.ApplicationCommandOptionBoolean,
							Name:        "regex",
							Description: "Whether to use a RegEx!",
							Required:    false,
						},
						{
							Type:        discordgo.ApplicationCommandOptionBoolean,
							Name:        "postfix",
							Description: "Whether to postfix!",
							Required:    false,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "inventory",
					Description: "Search for an element by name in an inventory!",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "query",
							Description: "The query to search with!",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionUser,
							Name:        "user",
							Description: "The user who's inventory to search!",
							Required:    false,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "sort",
							Description: "How to sort the results!",
							Choices:     eodsort.SortChoices,
							Required:    false,
						},
						{
							Type:        discordgo.ApplicationCommandOptionBoolean,
							Name:        "regex",
							Description: "Whether to use a RegEx!",
							Required:    false,
						},
						{
							Type:        discordgo.ApplicationCommandOptionBoolean,
							Name:        "postfix",
							Description: "Whether to postfix!",
							Required:    false,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "category",
					Description: "Search for an element by name in a category!",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "query",
							Description: "The query to search with!",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "category",
							Description: "The category to search in!",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "sort",
							Description: "How to sort the results!",
							Choices:     eodsort.SortChoices,
							Required:    false,
						},
						{
							Type:        discordgo.ApplicationCommandOptionBoolean,
							Name:        "regex",
							Description: "Whether to use a RegEx!",
							Required:    false,
						},
						{
							Type:        discordgo.ApplicationCommandOptionBoolean,
							Name:        "postfix",
							Description: "Whether to postfix!",
							Required:    false,
						},
					},
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
					Description: "Set the color of a category!",
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
		{
			Name:        "resetpolls",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Reset the polls!",
			Options:     []*discordgo.ApplicationCommandOption{},
		},
		{
			Name:        "ping",
			Type:        discordgo.ChatApplicationCommand,
			Description: "Check latency!",
			Options:     []*discordgo.ApplicationCommandOption{},
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"set": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData().Options[0]
			switch resp.Name {
			case "newschannel":
				bot.basecmds.SetNewsChannel(resp.Options[0].ChannelValue(bot.dg).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
			case "votingchannel":
				bot.basecmds.SetVotingChannel(resp.Options[0].ChannelValue(bot.dg).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
			case "votes":
				bot.basecmds.SetVoteCount(int(resp.Options[0].IntValue()), bot.newMsgSlash(i), bot.newRespSlash(i))
			case "polls":
				bot.basecmds.SetPollCount(int(resp.Options[0].IntValue()), bot.newMsgSlash(i), bot.newRespSlash(i))
			case "modrole":
				bot.basecmds.SetModRole(resp.Options[0].RoleValue(bot.dg, i.GuildID).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
			case "playchannel":
				isPlayChannel := true
				if len(resp.Options) > 1 {
					isPlayChannel = resp.Options[1].BoolValue()
				}
				bot.basecmds.SetPlayChannel(resp.Options[0].ChannelValue(bot.dg).ID, isPlayChannel, bot.newMsgSlash(i), bot.newRespSlash(i))
			}
		},
		"suggest": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			autocapitalize := true
			if len(resp.Options) > 1 {
				autocapitalize = resp.Options[1].BoolValue()
			}
			bot.elements.SuggestCmd(resp.Options[0].StringValue(), autocapitalize, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"mark": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			bot.polls.MarkCmd(resp.Options[0].StringValue(), resp.Options[1].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"image": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData().Options[0]
			switch resp.Name {
			case "element":
				bot.polls.ImageCmd(resp.Options[0].StringValue(), resp.Options[1].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))

			case "category":
				bot.polls.CatImgCmd(resp.Options[0].StringValue(), resp.Options[1].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
			}
		},
		"inv": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			sortby := "name"
			filter := "none"
			id := i.Member.User.ID

			defaul := true
			postfix := false
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

				if val.Name == "postfix" {
					postfix = val.BoolValue()
					defaul = false
				}
			}
			bot.elements.InvCmd(id, bot.newMsgSlash(i), bot.newRespSlash(i), sortby, filter, postfix, defaul)
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
			bot.elements.LbCmd(bot.newMsgSlash(i), bot.newRespSlash(i), sort, user)
		},
		"lbimage": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			sort := "count"
			for _, opt := range resp.Options {
				if opt.Name == "sortby" {
					sort = resp.Options[0].StringValue()
				}
			}
			bot.elements.LbImageCmd(bot.newMsgSlash(i), bot.newRespSlash(i), sort)
		},
		"addcat": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			suggestAdd := []string{resp.Options[1].StringValue()}
			if len(resp.Options) > 2 {
				for _, val := range resp.Options[2:] {
					suggestAdd = append(suggestAdd, val.StringValue())
				}
			}
			bot.categories.CategoryCmd(suggestAdd, resp.Options[0].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"cat": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			isAll := true
			sort := "name"
			catName := ""
			hasUser := false
			postfix := false
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

				if val.Name == "postfix" {
					postfix = val.BoolValue()
				}
			}

			if isAll {
				bot.categories.AllCatCmd(sort, hasUser, user, bot.newMsgSlash(i), bot.newRespSlash(i))
				return
			}

			bot.categories.CatCmd(catName, sort, hasUser, user, postfix, bot.newMsgSlash(i), bot.newRespSlash(i))
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

			bot.elements.HintCmd(elem, hasElem, false, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"stats": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.basecmds.StatsCmd(bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"resetinv": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			bot.elements.ResetInvCmd(resp.Options[0].UserValue(bot.dg).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"give": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData().Options[0]
			switch resp.Name {
			case "element":
				bot.treecmds.GiveCmd(resp.Options[0].StringValue(), resp.Options[1].BoolValue(), resp.Options[2].UserValue(bot.dg).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
			case "cat":
				bot.treecmds.GiveCatCmd(resp.Options[0].StringValue(), resp.Options[1].BoolValue(), resp.Options[2].UserValue(bot.dg).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
			case "all":
				bot.treecmds.GiveAllCmd(resp.Options[0].UserValue(bot.dg).ID, bot.newMsgSlash(i), bot.newRespSlash(i))
			}

		},
		"path": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData().Options[0]
			switch resp.Name {
			case "element":
				bot.treecmds.CalcTreeCmd(resp.Options[0].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))

			case "category":
				bot.treecmds.CalcTreeCatCmd(resp.Options[0].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
			}
		},
		"elemsort": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			postfix := false
			if len(resp.Options) > 1 {
				postfix = resp.Options[1].BoolValue()
			}
			bot.elements.SortCmd(resp.Options[0].StringValue(), postfix, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"help": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.basecmds.HelpCmd(bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"rmcat": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			suggestAdd := []string{resp.Options[1].StringValue()}
			if len(resp.Options) > 2 {
				for _, val := range resp.Options[2:] {
					suggestAdd = append(suggestAdd, val.StringValue())
				}
			}
			bot.categories.RmCategoryCmd(suggestAdd, resp.Options[0].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
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
			bot.elements.IdeaCmd(count, catName, hasCat, elName, hasEl, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"ai_idea": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.elements.AiCmd(bot.newMsgSlash(i), bot.newRespSlash(i))
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
				bot.elements.DownloadInvCmd(id, sortby, filter, postfix, bot.newMsgSlash(i), bot.newRespSlash(i))
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
				bot.categories.DownloadCatCmd(catName, sortby, postfix, bot.newMsgSlash(i), bot.newRespSlash(i))
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
				bot.treecmds.ElemBreakdownCmd(resp.Options[0].StringValue(), calctree, bot.newMsgSlash(i), bot.newRespSlash(i))

			case "category":
				calctree := false
				if len(resp.Options) > 1 {
					calctree = resp.Options[1].BoolValue()
				}
				bot.treecmds.CatBreakdownCmd(resp.Options[0].StringValue(), calctree, bot.newMsgSlash(i), bot.newRespSlash(i))

			case "inv":
				user := i.Member.User.ID
				calcTree := false
				for _, opt := range resp.Options {
					if opt.Name == "user" {
						user = opt.UserValue(bot.dg).ID
					}

					if opt.Name == "calctree" {
						calcTree = opt.BoolValue()
					}
				}
				bot.treecmds.InvBreakdownCmd(user, calcTree, bot.newMsgSlash(i), bot.newRespSlash(i))
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
				bot.treecmds.ElemGraphCmd(elem, layout, outputType, special, bot.newMsgSlash(i), bot.newRespSlash(i))

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
				bot.treecmds.CatGraphCmd(catName, layout, outputType, special, bot.newMsgSlash(i), bot.newRespSlash(i))
			}

		},
		"get": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData().Options[0]
			switch resp.Name {
			case "found":
				bot.elements.FoundCmd(resp.Options[0].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))

			case "categories":
				bot.categories.CategoriesCmd(resp.Options[0].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))

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
				bot.elements.Info(elem, id, isID, bot.newMsgSlash(i), rsp)
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
			bot.basecmds.SetUserColor(color, rmColor, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"invhint": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			bot.elements.HintCmd(resp.Options[0].StringValue(), true, true, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"search": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData().Options[0]
			switch resp.Name {
			case "elements":
				regex := false
				postfix := true
				sort := "name"
				for _, opt := range resp.Options {
					if opt.Name == "regex" {
						regex = opt.BoolValue()
					}

					if opt.Name == "sort" {
						sort = opt.StringValue()
					}

					if opt.Name == "postfix" {
						postfix = opt.BoolValue()
					}
				}
				bot.elements.SearchCmd(resp.Options[0].StringValue(), sort, "elements", "", regex, postfix, bot.newMsgSlash(i), bot.newRespSlash(i))

			case "inventory":
				regex := false
				sort := "name"
				postfix := true
				m := bot.newMsgSlash(i)
				user := m.Author.ID
				for _, opt := range resp.Options {
					if opt.Name == "regex" {
						regex = opt.BoolValue()
					}

					if opt.Name == "sort" {
						sort = opt.StringValue()
					}

					if opt.Name == "user" {
						user = opt.UserValue(bot.dg).ID
					}

					if opt.Name == "postfix" {
						postfix = opt.BoolValue()
					}
				}
				bot.elements.SearchCmd(resp.Options[0].StringValue(), sort, "inventory", user, regex, postfix, m, bot.newRespSlash(i))

			case "category":
				regex := false
				sort := "name"
				postfix := true
				m := bot.newMsgSlash(i)
				var category string
				for _, opt := range resp.Options {
					if opt.Name == "regex" {
						regex = opt.BoolValue()
					}

					if opt.Name == "sort" {
						sort = opt.StringValue()
					}

					if opt.Name == "category" {
						category = opt.StringValue()
					}

					if opt.Name == "postfix" {
						postfix = opt.BoolValue()
					}
				}
				bot.elements.SearchCmd(resp.Options[0].StringValue(), sort, "category", category, regex, postfix, m, bot.newRespSlash(i))
			}
		},
		"View Inventory": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			bot.elements.InvCmd(resp.TargetID, bot.newMsgSlash(i), bot.newRespSlash(i), "name", "none", false, true)
		},
		"View Info": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			rsp := bot.newRespSlash(i)
			id, res, suc := bot.getMessageElem(resp.TargetID, i.GuildID)
			if !suc {
				rsp.ErrorMessage(res)
				return
			}
			db, r := bot.GetDB(i.GuildID)
			if !r.Exists {
				return
			}
			elem, r := db.GetElement(id)
			if !r.Exists {
				return
			}
			bot.elements.InfoCmd(elem.Name, bot.newMsgSlash(i), rsp)
		},
		"Get Hint": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			rsp := bot.newRespSlash(i)
			id, res, suc := bot.getMessageElem(resp.TargetID, i.GuildID)
			if !suc {
				rsp.ErrorMessage(res)
				return
			}
			db, r := bot.GetDB(i.GuildID)
			if !r.Exists {
				return
			}
			elem, r := db.GetElement(id)
			if !r.Exists {
				return
			}
			bot.elements.HintCmd(elem.Name, true, false, bot.newMsgSlash(i), rsp)
		},
		"Get Inverse Hint": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			rsp := bot.newRespSlash(i)
			id, res, suc := bot.getMessageElem(resp.TargetID, i.GuildID)
			if !suc {
				rsp.ErrorMessage(res)
				return
			}
			db, r := bot.GetDB(i.GuildID)
			if !r.Exists {
				return
			}
			elem, r := db.GetElement(id)
			if !r.Exists {
				return
			}
			bot.elements.HintCmd(elem.Name, true, true, bot.newMsgSlash(i), rsp)
		},
		"Get Color": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			rsp := bot.newRespSlash(i)
			color, err := bot.base.GetColor(i.GuildID, resp.TargetID)
			if rsp.Error(err) {
				return
			}
			hex := strconv.FormatInt(int64(color), 16)
			rsp.Message(fmt.Sprintf("https://singlecolorimage.com/get/%s/100x100", hex))
		},
		"View Leaderboard": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData()
			bot.elements.LbCmd(bot.newMsgSlash(i), bot.newRespSlash(i), "count", resp.TargetID)
		},
		"notation": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			resp := i.ApplicationCommandData().Options[0]
			switch resp.Name {
			case "element":
				bot.treecmds.NotationCmd(resp.Options[0].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))

			case "category":
				bot.treecmds.CatNotationCmd(resp.Options[0].StringValue(), bot.newMsgSlash(i), bot.newRespSlash(i))
			}
		},
		"View Inventory Breakdown": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.treecmds.InvBreakdownCmd(i.ApplicationCommandData().TargetID, false, bot.newMsgSlash(i), bot.newRespSlash(i))
		},
		"resetpolls": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			bot.polls.ResetPolls(bot.newMsgSlash(i), bot.newRespSlash(i))
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
				bot.polls.ColorCmd(resp.Options[0].StringValue(), int(color), bot.newMsgSlash(i), rsp)

			case "category":
				var color int64 = 0
				if len(resp.Options) > 1 {
					var err error
					color, err = strconv.ParseInt(resp.Options[1].StringValue(), 16, 64)
					if rsp.Error(err) {
						return
					}
				}
				bot.polls.CatColorCmd(resp.Options[0].StringValue(), int(color), bot.newMsgSlash(i), bot.newRespSlash(i))
			}
		},
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			rsp := bot.newRespSlash(i)
			start := time.Now()
			rsp.Acknowledge()
			rsp.Message(fmt.Sprintf("🏓 Pong! Latency: **%s** [Slash Command]", time.Since(start)))
		},
	}
	autocompleteHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"get": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			data := i.ApplicationCommandData().Options[0]
			// autocomplete element names
			names, res := bot.elements.Autocomplete(bot.newMsgSlash(i), data.Options[0].StringValue())
			if !res.Exists {
				return
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{
					Choices: stringsToAutocomplete(names),
				},
			})
		},
		"cat": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			data := i.ApplicationCommandData()

			// autocomplete element names
			if len(data.Options) < 1 {
				return
			}
			// Check for focused and being named category
			focusedInd := -1
			for i, opt := range data.Options {
				if opt.Focused {
					focusedInd = i
					if opt.Name != "category" {
						return
					}
					break
				}
			}
			if focusedInd == -1 {
				return
			}

			names, res := bot.categories.Autocomplete(bot.newMsgSlash(i), data.Options[focusedInd].StringValue())
			if !res.Exists {
				return
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{
					Choices: stringsToAutocomplete(names),
				},
			})
		},
		"addcat": catChangeAutocomplete,
		"rmcat":  catChangeAutocomplete,
	}
)

func catChangeAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	// autocomplete element names
	if len(data.Options) < 1 {
		return
	}
	// Check for focused and being named category
	focusedInd := -1
	isElem := false
	for i, opt := range data.Options {
		if opt.Focused {
			focusedInd = i
			if opt.Name != "category" {
				isElem = true
			}
			break
		}
	}
	if focusedInd == -1 {
		return
	}

	var names []string
	var res types.GetResponse
	if isElem {
		names, res = bot.elements.Autocomplete(bot.newMsgSlash(i), data.Options[focusedInd].StringValue())
		if !res.Exists {
			return
		}
	} else {
		names, res = bot.categories.Autocomplete(bot.newMsgSlash(i), data.Options[focusedInd].StringValue())
		if !res.Exists {
			return
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: stringsToAutocomplete(names),
		},
	})
}
