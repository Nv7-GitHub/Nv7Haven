package elements

import (
	"fmt"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

var aiCmp = discordgo.ActionsRow{
	Components: []discordgo.MessageComponent{
		discordgo.Button{
			Label:    "New AI Generated Idea",
			Style:    discordgo.SuccessButton,
			CustomID: "idea",
		},
	},
}

type aiComponent struct {
	b *Elements
}

func (c *aiComponent) Handler(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	res, suc := c.b.genAi(i.GuildID)
	if !suc {
		res += " " + types.RedCircle
	}
	err := c.b.dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    res,
			Components: []discordgo.MessageComponent{aiCmp},
		},
	})
	if err != nil {
		fmt.Println("Failed to send message:", err)
	}
}

func (b *Elements) genAi(guild string) (string, bool) {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return res.Message, false
	}

	comb := db.AI.PredictCombo()

	text := ""
	for i, el := range comb {
		el, _ := db.GetElement(el)
		text += el.Name
		if i != len(comb)-1 {
			text += " + "
		}
	}

	return fmt.Sprintf("Your AI generated combination is... **%s**", text), true
}

func (b *Elements) AiCmd(m types.Msg, rsp types.Rsp) {
	res, suc := b.genAi(m.GuildID)
	if !suc {
		rsp.ErrorMessage(res)
		return
	}
	rsp.Acknowledge()

	data, ex := b.GetData(m.GuildID)
	if !ex.Exists {
		rsp.ErrorMessage(ex.Message)
		return
	}

	id := rsp.Message(res, aiCmp)

	data.AddComponentMsg(id, &aiComponent{
		b: b,
	})
}
