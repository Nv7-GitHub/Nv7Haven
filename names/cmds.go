package names

import "github.com/bwmarrin/discordgo"

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "set",
			Description: "Sets a user's name!",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "user",
					Type:        discordgo.ApplicationCommandOptionUser,
					Description: "The user to set the name of!",
					Required:    true,
				},
				{
					Name:        "name",
					Type:        discordgo.ApplicationCommandOptionString,
					Description: "The name of the user!",
					Required:    true,
				},
			},
		},
		{
			Name:        "get",
			Description: "Gets a user's name!",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "user",
					Type:        discordgo.ApplicationCommandOptionUser,
					Description: "The user to set the name of!",
					Required:    true,
				},
			},
		},
		{
			Name:        "search",
			Description: "Searches for a user based on their name!",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:         "name",
					Type:         discordgo.ApplicationCommandOptionString,
					Description:  "The name of the user!",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		{
			Name: "View Name",
			Type: discordgo.UserApplicationCommand,
		},
	}
)
