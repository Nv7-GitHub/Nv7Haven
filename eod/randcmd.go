package eod

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

var ideaCmp = discordgo.ActionsRow{
	Components: []discordgo.MessageComponent{
		discordgo.Button{
			Label:    "New Idea",
			Style:    discordgo.SuccessButton,
			CustomID: "idea",
		},
	},
}

type ideaComponent struct {
	catName  string
	hasCat   bool
	elemName string
	hasEl    bool
	count    int
	b        *EoD
}

func (c *ideaComponent) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	res, suc := c.b.genIdea(c.count, c.catName, c.hasCat, c.elemName, c.hasEl, i.GuildID, i.Member.User.ID)
	if !suc {
		res += " " + redCircle
	}
	err := c.b.dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    res,
			Components: []discordgo.MessageComponent{ideaCmp},
		},
	})
	if err != nil {
		fmt.Println("Failed to send message:", err)
	}
}

func (b *EoD) genIdea(count int, catName string, hasCat bool, elemName string, hasEl bool, guild string, author string) (string, bool) {
	if count > maxComboLength {
		return fmt.Sprintf("You can only combine up to %d elements!", maxComboLength), false
	}

	if count < 2 {
		return "You must combine at least 2 elements!", false
	}

	lock.RLock()
	dat, exists := b.dat[guild]
	lock.RUnlock()
	if !exists {
		return "Guild not found", false
	}

	dat.Lock.RLock()
	inv, exists := dat.InvCache[author]
	dat.Lock.RUnlock()
	if !exists {
		return "You don't have an inventory!", false
	}

	if hasEl {
		elName := strings.ToLower(elemName)
		dat.Lock.RLock()
		el, exists := dat.ElemCache[elName]
		dat.Lock.RUnlock()
		if !exists {
			return fmt.Sprintf("Element **%s** doesn't exist!", elemName), false
		} else {
			elemName = elName
			count--
		}

		dat.Lock.RLock()
		_, exists = inv[elemName]
		dat.Lock.RUnlock()
		if !exists {
			return fmt.Sprintf("Element **%s** is not in your inventory!", el.Name), false
		}
	}

	els := inv
	if hasCat {
		cat, exists := dat.CatCache[strings.ToLower(catName)]
		if !exists {
			return fmt.Sprintf("Category **%s** doesn't exist!", catName), false
		}
		els = make(map[string]types.Empty)

		for el := range cat.Elements {
			l := strings.ToLower(el)
			_, exists := inv[l]
			if exists {
				els[l] = types.Empty{}
			}
		}

		if len(els) == 0 {
			return fmt.Sprintf("You don't have any elements in category **%s**!", cat.Name), false
		}
	}

	var elem3 string
	cont := true
	var elems []string
	tries := 0
	for cont {
		elems = make([]string, count)
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
			elems = append([]string{elemName}, elems...)
		}

		query := "SELECT elem3 FROM eod_combos WHERE elems LIKE ? AND guild=?"
		els := elems2txt(elems)
		if isASCII(els) {
			query = "SELECT elem3 FROM eod_combos WHERE CONVERT(elems USING utf8mb4) LIKE CONVERT(? USING utf8mb4) AND guild=CONVERT(? USING utf8mb4) COLLATE utf8mb4_general_ci"
		}
		row := b.db.QueryRow(query, els, guild)
		err := row.Scan(&elem3)
		if err != nil {
			cont = false
		}
		tries++

		if tries > 10 {
			return "Couldn't find a random unused combination, maybe try again later?", false
		}
	}

	text := ""
	for i, el := range elems {
		dat.Lock.RLock()
		text += dat.ElemCache[strings.ToLower(el)].Name
		dat.Lock.RUnlock()
		if i != len(elems)-1 {
			text += " + "
		}
	}

	if dat.CombCache == nil {
		dat.CombCache = make(map[string]types.Comb)
	}
	dat.Lock.Lock()
	dat.CombCache[author] = types.Comb{
		Elems: elems,
		Elem3: "",
	}
	dat.Lock.Unlock()

	lock.Lock()
	b.dat[guild] = dat
	lock.Unlock()

	return fmt.Sprintf("Your random unused combination is... **%s**\n 	Suggest it by typing **/suggest**", text), true
}
func (b *EoD) ideaCmd(count int, catName string, hasCat bool, elemName string, hasEl bool, m types.Msg, rsp types.Rsp) {
	res, suc := b.genIdea(count, catName, hasCat, elemName, hasEl, m.GuildID, m.Author.ID)
	if !suc {
		rsp.ErrorMessage(res)
		return
	}
	rsp.Acknowledge()

	lock.Lock()
	dat, exists := b.dat[m.GuildID]
	lock.Unlock()
	if !exists {
		rsp.ErrorMessage("Guild not found")
		return
	}

	id := rsp.Message(res, ideaCmp)

	dat.Lock.Lock()
	dat.ComponentMsgs[id] = &ideaComponent{
		catName:  catName,
		count:    count,
		hasCat:   hasCat,
		elemName: elemName,
		hasEl:    hasEl,
		b:        b,
	}
	dat.Lock.Unlock()

	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()
}
