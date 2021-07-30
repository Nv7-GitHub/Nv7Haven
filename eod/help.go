package eod

import "github.com/bwmarrin/discordgo"

//go:embed help/about.txt
var helpAbout string

//go:embed help/basics.txt
var helpBasics string

//go:embed help/advanced.txt
var helpAdvanced string

//go:embed help/setup.txt
var helpSetup string

var helpComponents = discordgo.ActionsRow{
	Components: []discordgo.MessageComponent{
		discordgo.SelectMenu{
			CustomID: "help-select",
			Options: []discordgo.SelectMenuOption{
				{
					Label:       "About",
					Value:       "about",
					Description: "Get basic information about the bot!",
					Default:     true,
				},
				{
					Label:       "Basics",
					Value:       "basics",
					Description: "Learn the basics about using the bot!",
					Default:     false,
				},
				{
					Label:       "Advanced",
					Value:       "advanced",
					Description: "Learn how to use the advanced features of the bot!",
					Default:     false,
				},
				{
					Label:       "Setup",
					Value:       "setup",
					Description: "Learn how to set up your own EoD server!",
					Default:     false,
				},
			},
		},
	},
}

type helpComponent struct {
	b *EoD
}

func (h *helpComponent) handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var txt string

	switch i.MessageComponentData().Values[0] {
	case "about":
		txt = helpAbout
	case "basics":
		txt = helpBasics
	case "advanced":
		txt = helpAdvanced
	case "setup":
		txt = helpSetup
	}

	h.b.dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    txt,
			Components: []discordgo.MessageComponent{helpComponents},
		},
	})
}

func (b *EoD) helpCmd(rsp rsp) {
	rsp.Resp(helpAbout, helpComponents)
}
