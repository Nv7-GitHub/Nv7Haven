package names

import "github.com/bwmarrin/discordgo"

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "setname",
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
	}
)
