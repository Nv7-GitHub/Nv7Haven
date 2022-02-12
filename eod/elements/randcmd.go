package elements

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

func newIdeaCmp(db *eodb.DB) discordgo.ActionsRow {
	return discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    db.Config.LangProperty("NewIdea", nil),
				Style:    discordgo.SuccessButton,
				CustomID: "idea",
				Emoji: discordgo.ComponentEmoji{
					Name:     "idea",
					ID:       "932832178847502386",
					Animated: false,
				},
			},
		},
	}
}

type ideaComponent struct {
	catName  string
	hasCat   bool
	elemName string
	hasEl    bool
	count    int
	b        *Elements
	db       *eodb.DB
}

func (c *ideaComponent) Handler(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	res, suc := c.b.genIdea(c.count, c.catName, c.hasCat, c.elemName, c.hasEl, i.GuildID, i.Member.User.ID)
	if !suc {
		res += " " + types.RedCircle
	}
	err := c.b.dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    res,
			Components: []discordgo.MessageComponent{newIdeaCmp(c.db)},
		},
	})
	if err != nil {
		fmt.Println("Failed to send message:", err)
	}
}

func (b *Elements) genIdea(count int, catName string, hasCat bool, elemName string, hasEl bool, guild string, author string) (string, bool) {
	db, res := b.GetDB(guild)
	if !res.Exists {
		return res.Message, false
	}

	if count > types.MaxComboLength {
		return db.Config.LangProperty("MaxCombine", types.MaxComboLength), false
	}

	if count < 2 {
		return db.Config.LangProperty("MustCombine", 2), false
	}
	inv := db.GetInv(author)

	var elID int
	if hasEl {
		elName := strings.ToLower(elemName)

		el, res := db.GetElementByName(elName)
		if !res.Exists {
			return res.Message, false
		} else {
			count--
		}
		elID = el.ID

		exists := inv.Contains(el.ID)
		if !exists {
			return db.Config.LangProperty("DontHave", el.Name), false
		}
	}

	els := inv.Elements
	if hasCat {
		cat, res := db.GetCat(catName)
		if !res.Exists {
			return res.Message, false
		}
		els = make(map[int]types.Empty)

		for el := range cat.Elements {
			exists := inv.Contains(el)
			if exists {
				els[el] = types.Empty{}
			}
		}

		if len(els) == 0 {
			return db.Config.LangProperty("HaveNoElemsInCat", cat.Name), false
		}
	}

	res = types.GetResponse{Exists: true}
	var elems []int
	tries := 0
	for res.Exists {
		elems = make([]int, count)
		for i := range elems {
			cnt := rand.Intn(len(els))
			j := 0
			for k := range els {
				if j == cnt {
					elems[i] = k
					break
				}
				j++
			}
		}
		if hasEl {
			elems = append([]int{elID}, elems...)
		}

		_, res = db.GetCombo(elems)
		tries++

		if tries > types.MaxTries {
			return db.Config.LangProperty("FailedIdea", nil), false
		}
	}

	text := ""
	for i, el := range elems {
		el, _ := db.GetElement(el)
		text += el.Name
		if i != len(elems)-1 {
			text += " + "
		}
	}

	data, _ := b.GetData(guild)
	data.SetComb(author, types.Comb{
		Elems: elems,
		Elem3: -1,
	})

	return db.Config.LangProperty("YourIdea", map[string]interface{}{
		"Combo":       text,
		"SuggestText": db.Config.LangProperty("SuggestIdea", nil),
	}), true
}

func (b *Elements) IdeaCmd(count int, catName string, hasCat bool, elemName string, hasEl bool, m types.Msg, rsp types.Rsp) {
	res, suc := b.genIdea(count, catName, hasCat, elemName, hasEl, m.GuildID, m.Author.ID)
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
	db, _ := b.GetDB(m.GuildID)

	id := rsp.Message(res, newIdeaCmp(db))

	data.AddComponentMsg(id, &ideaComponent{
		catName:  catName,
		count:    count,
		hasCat:   hasCat,
		elemName: elemName,
		hasEl:    hasEl,
		b:        b,
		db:       db,
	})
}
