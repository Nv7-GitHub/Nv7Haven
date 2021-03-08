package discord

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/elemental"
	"github.com/bwmarrin/discordgo"
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

	_, err = b.db.Exec("UPDATE elements SET comment=? WHERE name=?", mark, id)
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
		b.createCmd(comb.elem1, comb.elem2, u.Metadata["eusername"].(string), name, u.Metadata["uid"].(string), m, rsp)
		return
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
	upperName := strings.ToUpper(name)
	for _, val := range combs {
		if strings.ToUpper(val) == upperName {
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
		b.createCmd(comb.elem1, comb.elem2, u.Metadata["eusername"].(string), name, u.Metadata["uid"].(string), m, rsp)
		return
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
	upperName := strings.ToUpper(name)
	for _, val := range combs {
		if strings.ToUpper(val) == upperName {
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
		return
	}

	err := b.e.NewFound(id, uid)
	if rsp.Error(err) {
		return
	}

	rsp.Resp(fmt.Sprintf("Succesfully created element %s! You can use the `mark` command to add a creator mark!", id))
}

func (b *Bot) randomCmd(getter func(uid string) ([]string, error), m msg, rsp rsp) {
	u, suc := b.getUser(m, rsp, m.Author.ID)
	if !suc {
		return
	}
	uid := u.Metadata["uid"].(string)
	for i := 0; i < 10; i++ {
		combs, err := getter(uid)
		if err != nil {
			return
		}
		if len(combs) == 2 {
			elem1 := combs[0]
			elem2 := combs[1]

			b.combos[m.Author.ID] = comb{
				elem1: elem1,
				elem2: elem2,
				elem3: "",
			}

			combs, err := b.e.GetSuggestions(elem1, elem2)
			if rsp.Error(err) {
				return
			}
			rsp.Embed(&discordgo.MessageEmbed{
				Title:       fmt.Sprintf("Suggestions for %s+%s", elem1, elem2),
				Description: strings.Join(combs, "\n"),
			})
			return
		}
	}
	rsp.ErrorMessage("No available random lonely suggestions right now! Check back later!")
}
