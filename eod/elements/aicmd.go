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
			Emoji: discordgo.ComponentEmoji{
				Name:     "ai",
				ID:       "932832511459999844",
				Animated: false,
			},
		},
	},
}

type aiComponent struct {
	b *Elements
}

func (c *aiComponent) Handler(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	res, suc := c.b.genAi(i.GuildID, i.Member.User.ID)
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

func (b *Elements) genAi(guild string, author string) (string, bool) {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return res.Message, false
	}

	tries := 0
	var comb []int
	success := false
	for !success {
		comb = db.AI.PredictCombo()
		_, res := db.GetCombo(comb)
		if !res.Exists {
			success = true
		}
		tries++
		if tries > types.MaxTries {
			return "Failed to generate a valid idea!", false
		}
	}

	text := ""
	for i, el := range comb {
		el, _ := db.GetElement(el)
		text += el.Name
		if i != len(comb)-1 {
			text += " + "
		}
	}

	// Check if you can suggest
	suggest := ""
	canSuggest := true
	inv := db.GetInv(author)
	for _, el := range comb {
		if !inv.Contains(el) {
			canSuggest = false
			break
		}
	}
	if canSuggest {
		data, res := b.GetData(guild)
		if !res.Exists {
			return res.Message, false
		}
		suggest = "\n 	Suggest it by typing **/suggest**"
		data.SetComb(author, types.Comb{
			Elems: comb,
			Elem3: -1,
		})
	}

	return fmt.Sprintf("Your AI generated combination is... **%s**%s", text, suggest), true
}

func (b *Elements) AiCmd(m types.Msg, rsp types.Rsp) {
	res, suc := b.genAi(m.GuildID, m.Author.ID)
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
