package discord

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/elemental"
)

func (b *Bot) markCmd(mark string, id string, m msg, rsp rsp) {
	b.checkUser(m, rsp)
	if !b.isLoggedIn(m, rsp) {
		return
	}

	exts, suc := b.exts(rsp, "elements", "name=?", id)
	if !suc {
		return
	}
	if !exts {
		rsp.ErrorMessage(fmt.Sprintf("Element %s doesn't exist!", id))
		return
	}

	elem, err := b.e.GetElement(id)
	if rsp.Error(err) {
		return
	}

	u, suc := b.getUser(m, rsp, m.Author.ID)
	if !suc {
		return
	}

	if elem.Comment != "None" {
		rsp.ErrorMessage("The element already has a creator mark!")
		return
	}
	if elem.Pioneer != u.Metadata["eusername"].(string) {
		rsp.ErrorMessage("You didn't make this element!")
		return
	}

	_, err = b.db.Exec("UPDATE elemnts SET comment=? WHERE name=?", id)
	if rsp.Error(err) {
		return
	}

	rsp.Resp("Succesfully added creator mark!")
}

func (b *Bot) suggestCmd(name string, color string, m msg, rsp rsp) {
	b.checkUser(m, rsp)
	if !b.isLoggedIn(m, rsp) {
		return
	}

	name = strings.TrimSpace(name)
	comb, exists := b.combos[m.Author.ID]
	if !exists {
		rsp.ErrorMessage("You haven't combined any elements!")
	}

	_, comboExists, err := b.e.GetCombo(comb.elem1, comb.elem2)
	if rsp.Error(err) {
		return
	}

	if comboExists {
		rsp.ErrorMessage("Combo already exists!")
		return
	}

	combs, err := b.e.GetSuggestions(comb.elem1, comb.elem2)
	if rsp.Error(err) {
		return
	}

	for _, val := range combs {
		if val == name {
			rsp.ErrorMessage("Someone's already suggested that! Use the `upvote` command to upvote a suggestion!")
			return
		}
	}

	u, suc := b.getUser(m, rsp, m.Author.ID)
	if !suc {
		return
	}

	create, err := b.e.NewSuggestion(comb.elem1, comb.elem2, elemental.Suggestion{
		Name:    name,
		Creator: u.Metadata["eusername"].(string),
		Color: elemental.Color{
			Base: color,
		},
		Votes: 0,
		Voted: []string{u.Metadata["uid"].(string)},
	})
	if rsp.Error(err) {
		return
	}
	if create {
		b.createCmd(comb.elem1, comb.elem2, u.Metadata["eusername"].(string), u.Metadata["uid"].(string), name, m, rsp)
	}
	rsp.Resp("Succesfully created suggestion!")
}

func (b *Bot) upvoteCmd(name string, m msg, rsp rsp) {
	b.checkUser(m, rsp)
	if !b.isLoggedIn(m, rsp) {
		return
	}

	name = strings.TrimSpace(name)
	comb, exists := b.combos[m.Author.ID]
	if !exists {
		rsp.ErrorMessage("You haven't combined any elements!")
		return
	}

	_, comboExists, err := b.e.GetCombo(comb.elem1, comb.elem2)
	if rsp.Error(err) {
		return
	}

	if comboExists {
		rsp.ErrorMessage("Combo already exists!")
		return
	}

	combs, err := b.e.GetSuggestions(comb.elem1, comb.elem2)
	if rsp.Error(err) {
		return
	}

	isIn := false
	for _, val := range combs {
		if val == name {
			isIn = true
			break
		}
	}
	if !isIn {
		rsp.ErrorMessage("Suggestion doesn't exist! Use the `suggest` command to suggest something!")
		return
	}

	u, suc := b.getUser(m, rsp, m.Author.ID)
	if !suc {
		return
	}

	create, suc, msg := b.e.UpvoteSuggestion(name, u.Metadata["uid"].(string))
	if !suc {
		rsp.ErrorMessage(msg)
		return
	}
	if create {
		b.createCmd(comb.elem1, comb.elem2, u.Metadata["eusername"].(string), u.Metadata["uid"].(string), name, m, rsp)
	}
	rsp.Resp("Succesfully upvoted suggestion!")
}

func (b *Bot) downvoteCmd(name string, m msg, rsp rsp) {
	b.checkUser(m, rsp)
	if !b.isLoggedIn(m, rsp) {
		return
	}

	name = strings.TrimSpace(name)
	comb, exists := b.combos[m.Author.ID]
	if !exists {
		rsp.ErrorMessage("You haven't combined any elements!")
		return
	}

	_, comboExists, err := b.e.GetCombo(comb.elem1, comb.elem2)
	if rsp.Error(err) {
		return
	}

	if comboExists {
		rsp.ErrorMessage("Combo already exists!")
		return
	}

	combs, err := b.e.GetSuggestions(comb.elem1, comb.elem2)
	if rsp.Error(err) {
		return
	}

	isIn := false
	for _, val := range combs {
		if val == name {
			isIn = true
			break
		}
	}
	if !isIn {
		rsp.ErrorMessage("Suggestion doesn't exist! Use the `suggest` command to suggest something!")
		return
	}

	u, suc := b.getUser(m, rsp, m.Author.ID)
	if !suc {
		return
	}

	suc, msg := b.e.DownvoteSuggestion(name, u.Metadata["uid"].(string))
	if !suc {
		rsp.ErrorMessage(msg)
		return
	}
	rsp.Resp("Succesfully downvoted suggestion!")
}

func (b *Bot) createCmd(elem1 string, elem2 string, username string, id string, uid string, m msg, rsp rsp) {
	suc, msg := b.e.CreateSuggestion("None", username, elem1, elem2, id)
	if !suc {
		rsp.ErrorMessage(msg)
	}

	err := b.e.NewFound(id, uid)
	if rsp.Error(err) {
		return
	}

	rsp.Resp(fmt.Sprintf("Succesfully created element %s! You can use the `mark` command to add a creator mark!", id))
}
