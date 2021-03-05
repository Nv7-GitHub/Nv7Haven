package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type reactionMsgType int

const ldbPageSwitcher = 0
const invPageSwitcher = 1
const suggestionReaction = 2

var combs = []string{
	"+",
	",",
}

func (b *Bot) einvPageHandler(r *discordgo.MessageReactionAdd) {
	pg := b.pages[r.MessageID]
	var page int
	if r.Emoji.Name == leftArrow {
		page = pg.Metadata["page"].(int) - 1
	} else {
		page = pg.Metadata["page"].(int) + 1
	}
	inv := pg.Metadata["found"].([]string)
	if ((page * 20) > len(inv)) || (page < 0) {
		b.dg.MessageReactionsRemoveAll(r.ChannelID, r.MessageID)
		b.dg.MessageReactionAdd(r.ChannelID, r.MessageID, leftArrow)
		b.dg.MessageReactionAdd(r.ChannelID, r.MessageID, rightArrow)
		return
	}
	pg.Metadata["page"] = page

	text := ""
	for i := page * 20; i < len(inv); i++ {
		text += inv[i] + "\n"
		if i > 20+(page*20) {
			break
		}
	}

	b.dg.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s's Elemental Inventory", pg.Metadata["name"].(string)),
		Description: text,
	})
	b.dg.MessageReactionsRemoveAll(r.ChannelID, r.MessageID)
	b.dg.MessageReactionAdd(r.ChannelID, r.MessageID, leftArrow)
	b.dg.MessageReactionAdd(r.ChannelID, r.MessageID, rightArrow)
}

func (b *Bot) einvCmd(m msg, rsp rsp) {
	b.checkUser(m, rsp)
	if !b.isLoggedIn(m, rsp) {
		return
	}

	u, suc := b.getUser(m, rsp, m.Author.ID)
	if !suc {
		return
	}

	inv, err := b.e.GetFound(u.Metadata["uid"].(string))
	if rsp.Error(err) {
		return
	}

	text := ""
	for i := 0; i < len(inv); i++ {
		text += inv[i] + "\n"
		if i > 20 {
			break
		}
	}
	id := rsp.Embed(&discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s's Elemental Inventory", m.Author.Username),
		Description: text,
	})
	b.dg.MessageReactionAdd(m.ChannelID, id, leftArrow)
	b.dg.MessageReactionAdd(m.ChannelID, id, rightArrow)
	b.pages[id] = reactionMsg{
		Type: invPageSwitcher,
		Metadata: map[string]interface{}{
			"page":  0,
			"found": inv,
			"name":  m.Author.Username,
		},
		Handler: b.einvPageHandler,
	}
}

func (b *Bot) eloginCmd(username string, m msg, rsp rsp) {
	u, suc := b.getUser(m, rsp, m.Author.ID)
	if !suc {
		return
	}
	_, exists := u.Metadata["uid"]
	if exists {
		rsp.ErrorMessage("You are already logged in!")
		return
	}
	resp := b.e.CreateUser(username, m.Author.ID)
	if !resp.Success {
		rsp.ErrorMessage(resp.Data)
		return
	}
	u.Metadata["uid"] = resp.Data
	u.Metadata["eusername"] = username
	suc = b.updateUser(rsp, u)
	if !suc {
		return
	}
	rsp.Resp("Successfully logged in!")
}

func (b *Bot) comboCmd(elem1 string, elem2 string, m msg, rsp rsp) {
	elem1 = strings.TrimSpace(elem1)
	elem2 = strings.TrimSpace(elem2)
	if elem1 == "" || elem2 == "" {
		return
	}

	b.checkUser(m, rsp)
	if !b.isLoggedIn(m, rsp) {
		return
	}

	elem3, comboExists, err := b.e.GetCombo(elem1, elem2)
	if rsp.Error(err) {
		return
	}

	exts, suc := b.exts(rsp, "elements", "name=?", elem1)
	if !suc {
		return
	}
	if !exts {
		rsp.ErrorMessage(fmt.Sprintf("Element %s doesn't exist!", elem1))
		return
	}
	exts, suc = b.exts(rsp, "elements", "name=?", elem2)
	if !suc {
		return
	}
	if !exts {
		rsp.ErrorMessage(fmt.Sprintf("Element %s doesn't exist!", elem2))
		return
	}

	u, suc := b.getUser(m, rsp, m.Author.ID)
	if !suc {
		return
	}
	uidRaw, exists := u.Metadata["uid"]
	if !exists {
		rsp.ErrorMessage("Not logged in!")
		return
	}
	uid, ok := uidRaw.(string)
	if !ok {
		rsp.ErrorMessage("Invalid UID!")
		return
	}

	hasElem1 := false
	hasElem2 := false
	el1 := strings.ToUpper(elem1)
	el2 := strings.ToUpper(elem2)
	found, err := b.e.GetFound(uid)
	if rsp.Error(err) {
		return
	}
	for _, val := range found {
		v := strings.ToUpper(val)
		if v == el1 {
			hasElem1 = true
			if hasElem1 && hasElem2 {
				break
			}
		}
		if v == el2 {
			hasElem2 = true
			if hasElem1 && hasElem2 {
				break
			}
		}
	}

	if !hasElem1 {
		rsp.Resp(fmt.Sprintf("You haven't found element %s yet!", elem1))
		return
	}
	if !hasElem2 {
		rsp.Resp(fmt.Sprintf("You haven't found element %s yet!", elem2))
		return
	}

	if !comboExists {
		rsp.Resp("Combo doesn't exist, gotta suggest something")
		return
	}

	err = b.e.NewFound(elem3, uid)
	if rsp.Error(err) {
		return
	}
	rsp.Resp(fmt.Sprintf("You made %s!", elem3))
	b.combos[m.Author.ID] = comb{
		elem1: elem1,
		elem2: elem2,
		elem3: elem3,
	}
}

func (b *Bot) elementalHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if b.startsWith(m, "elogin") {
		msg := b.newMsgNormal(m)
		rsp := b.newRespNormal(m)
		b.checkUser(msg, rsp)

		var name string
		_, err := fmt.Sscanf(m.Content, "elogin %s", &name)
		if rsp.Error(err) {
			return
		}
		b.eloginCmd(name, msg, rsp)
		return
	}

	for _, comb := range combs {
		if strings.Contains(m.Content, comb) {
			parts := strings.Split(m.Content, comb)
			if len(parts) != 2 {
				return
			}

			msg := b.newMsgNormal(m)
			rsp := b.newRespNormal(m)
			b.comboCmd(parts[0], parts[1], msg, rsp)
			return
		}
	}

	if strings.HasPrefix(m.Content, "*2") {
		comb, exists := b.combos[m.Author.ID]
		if exists {
			msg := b.newMsgNormal(m)
			rsp := b.newRespNormal(m)
			b.comboCmd(comb.elem3, comb.elem3, msg, rsp)
		}
	}

	if b.startsWith(m, "einv") {
		b.einvCmd(b.newMsgNormal(m), b.newRespNormal(m))
	}
}

func (b *Bot) isLoggedIn(m msg, rsp rsp) bool {
	u, suc := b.getUser(m, rsp, m.Author.ID)
	if !suc {
		return false
	}

	_, exists := u.Metadata["uid"]
	if !exists {
		rsp.ErrorMessage("You need to get an account! Use the `elogin` command to login to Nv7's Elemental!")
		return false
	}

	return true
}
