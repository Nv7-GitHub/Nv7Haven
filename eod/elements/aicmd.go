package elements

import (
	"fmt"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

func newAiCmp(db *eodb.DB) discordgo.ActionsRow {
	return discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    db.Config.LangProperty("NewAIIdea", nil),
				Style:    discordgo.SuccessButton,
				CustomID: "idea",
				Emoji: discordgo.ComponentEmoji{
					Name:     "ai",
					ID:       "932832481768517672",
					Animated: false,
				},
			},
		},
	}
}

type aiComponent struct {
	b *Elements
}

func (c *aiComponent) Handler(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	res, suc, db := c.b.genAi(i.GuildID, i.Member.User.ID)
	components := []discordgo.MessageComponent{}
	if !suc {
		res += " " + types.RedCircle
	} else {
		components = []discordgo.MessageComponent{newAiCmp(db)}
	}
	err := c.b.dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    res,
			Components: components,
		},
	})
	if err != nil {
		fmt.Println("Failed to send message:", err)
	}
}

func (b *Elements) genAi(guild string, author string) (string, bool, *eodb.DB) {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return res.Message, false, nil
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
			return db.Config.LangProperty("FailedAIGenerate", nil), false, db
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
			return res.Message, false, db
		}
		suggest = db.Config.LangProperty("SuggestIdea", nil)
		data.SetComb(author, types.Comb{
			Elems: comb,
			Elem3: -1,
		})
	}

	return db.Config.LangProperty("YourAIIdea", map[string]any{
		"Idea":        text,
		"SuggestText": suggest,
	}), true, db
}

func (b *Elements) AiCmd(m types.Msg, rsp types.Rsp) {
	res, suc, db := b.genAi(m.GuildID, m.Author.ID)
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

	id := rsp.Message(res, newAiCmp(db))

	data.AddComponentMsg(id, &aiComponent{
		b: b,
	})
}
