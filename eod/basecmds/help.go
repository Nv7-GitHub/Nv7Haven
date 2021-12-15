package basecmds

import (
	_ "embed"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

//go:embed help/about.txt
var helpAbout string

//go:embed help/basics.txt
var helpBasics string

//go:embed help/advanced.txt
var helpAdvanced string

//go:embed help/setup.txt
var helpSetup string

func makeHelpComponents(selected string) discordgo.ActionsRow {
	return discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.SelectMenu{
				CustomID: "help-select",
				Options: []discordgo.SelectMenuOption{
					{
						Label:       "About",
						Value:       "about",
						Description: "Get basic information about the bot!",
						Default:     selected == "about",
					},
					{
						Label:       "Basics",
						Value:       "basics",
						Description: "Learn the basics about using the bot!",
						Default:     selected == "basics",
					},
					{
						Label:       "Advanced",
						Value:       "advanced",
						Description: "Learn how to use the advanced features of the bot!",
						Default:     selected == "advanced",
					},
					{
						Label:       "Setup",
						Value:       "setup",
						Description: "Learn how to set up your own EoD server!",
						Default:     selected == "setup",
					},
				},
			},
		},
	}
}

type helpComponent struct {
	b *BaseCmds
}

func (h *helpComponent) Handler(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	var txt string

	val := i.MessageComponentData().Values[0]
	switch val {
	case "about":
		txt = helpAbout
	case "basics":
		txt = helpBasics
	case "advanced":
		txt = helpAdvanced
	case "setup":
		txt = helpSetup
	default:
		return
	}

	h.b.dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    txt,
			Components: []discordgo.MessageComponent{makeHelpComponents(val)},
		},
	})
}

func (b *BaseCmds) HelpCmd(m types.Msg, rsp types.Rsp) {
	rsp.Acknowledge()
	id := rsp.Message(helpAbout, makeHelpComponents("about"))

	data, res := b.GetData(m.GuildID)
	if !res.Exists {
		return
	}

	data.AddComponentMsg(id, &helpComponent{
		b: b,
	})
}
