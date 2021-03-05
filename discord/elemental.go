package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var combs = []string{
	"+",
	",",
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

	b.checkUser(m, rsp)
	if !b.isLoggedIn(m, rsp) {
		return
	}

	elem3, exists, err := b.e.GetCombo(elem1, elem2)
	if rsp.Error(err) {
		return
	}
	if !exists {
		rsp.ErrorMessage("Combo doesn't exist!")
		return
	}

	exts, suc := b.exts(rsp, "elements", "name=?", elem1)
	if !suc {
		return
	}
	if !exts {
		rsp.ErrorMessage(fmt.Sprintf("Element %s doesn't exist!", elem1))
	}
	exts, suc = b.exts(rsp, "elements", "name=?", elem2)
	if !suc {
		return
	}
	if !exts {
		rsp.ErrorMessage(fmt.Sprintf("Element %s doesn't exist!", elem2))
	}

	rsp.Message("this combo makes " + elem3)
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
