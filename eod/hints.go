package eod

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var hintCmp = discordgo.ActionsRow{
	Components: []discordgo.MessageComponent{
		discordgo.Button{
			Label:    "New Hint",
			CustomID: "hint-new",
			Style:    discordgo.SuccessButton,
		},
	},
}

type hintComponent struct {
	b *EoD
}

func (h *hintComponent) handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	hint, msg, suc := h.b.getHint("", false, i.Member.User.ID, i.GuildID)
	if !suc {
		h.b.dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content:    msg,
				Components: []discordgo.MessageComponent{hintCmp},
			},
		})
		return
	}

	h.b.dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{hint},
			Components: []discordgo.MessageComponent{hintCmp},
		},
	})
}

type hintCombo struct {
	exists int
	text   string
}

var noObscure = map[rune]empty{
	' ': {},
	'.': {},
	'-': {},
	'_': {},
}

func obscure(val string) string {
	out := make([]rune, len([]rune(val)))
	i := 0
	for _, char := range val {
		_, exists := noObscure[char]
		if exists {
			out[i] = char
		} else {
			out[i] = '?'
		}
		i++
	}
	return string(out)
}

func (b *EoD) hintCmd(elem string, hasElem bool, m msg, rsp rsp) {
	rsp.Acknowledge()

	hint, msg, suc := b.getHint(elem, hasElem, m.Author.ID, m.GuildID)
	if !suc {
		rsp.ErrorMessage(msg)
		return
	}

	if hasElem {
		rsp.Embed(hint)
		return
	}

	id := rsp.Embed(hint, hintCmp)

	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	dat.lock.Lock()
	dat.componentMsgs[id] = &hintComponent{b: b}
	dat.lock.Unlock()

	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()
}

func (b *EoD) getHint(elem string, hasElem bool, author string, guild string) (*discordgo.MessageEmbed, string, bool) {
	lock.RLock()
	dat, exists := b.dat[guild]
	lock.RUnlock()
	if !exists {
		return nil, "Guild not found", false
	}

	dat.lock.RLock()
	inv, exists := dat.invCache[author]
	dat.lock.RUnlock()
	if !exists {
		return nil, "You don't have an inventory!", false
	}
	var el element
	if hasElem {
		dat.lock.RLock()
		el, exists = dat.elemCache[strings.ToLower(elem)]
		dat.lock.RUnlock()
		if !exists {
			return nil, fmt.Sprintf("No hints were found for **%s**!", elem), false
		}
	}
	if !hasElem {
		hasFound := false
		dat.lock.RLock()
		for _, v := range dat.elemCache {
			_, exists := inv[strings.ToLower(v.Name)]
			if !exists {
				el = v
				elem = v.Name
				hasFound = true
				break
			}
		}
		dat.lock.RUnlock()
		if !hasFound {
			dat.lock.RLock()
			for _, v := range dat.elemCache {
				el = v
				elem = v.Name
				hasFound = true
				break
			}
			dat.lock.RUnlock()
		}
	}

	var combs *sql.Rows
	var err error
	query := "SELECT elems FROM eod_combos WHERE elem3 LIKE ? AND guild=?"
	if isASCII(elem) {
		query = "SELECT elems FROM eod_combos WHERE CONVERT(elem3 USING utf8mb4) LIKE CONVERT(? USING utf8mb4) AND guild=CONVERT(? USING utf8mb4) COLLATE utf8mb4_general_ci"
	}

	if isWildcard(elem) {
		query = strings.ReplaceAll(query, " LIKE ", "=")
	}

	combs, err = b.db.Query(query, elem, guild)
	if err != nil {
		return nil, err.Error(), false
	}
	defer combs.Close()
	out := make([]hintCombo, 0)
	var elemTxt string

	length := 0
	for combs.Next() {
		err = combs.Scan(&elemTxt)
		if err != nil {
			return nil, err.Error(), false
		}
		elems := strings.Split(elemTxt, "+")

		txt, ex := getHintText(elems, inv, dat)
		out = append(out, hintCombo{
			exists: ex,
			text:   txt,
		})
		length += len(txt)
	}

	if len(out) == 0 {
		dat.lock.RLock()
		element := dat.elemCache[strings.ToLower(elem)]
		dat.lock.RUnlock()

		txt, ex := getHintText(element.Parents, inv, dat)
		out = append(out, hintCombo{
			exists: ex,
			text:   txt,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].exists > out[j].exists
	})

	text := ""
	for _, val := range out {
		text += val.text + "\n"
	}

	txt := "Don't "
	_, hasElem = inv[strings.ToLower(el.Name)]
	if hasElem {
		txt = ""
	}

	if len(text) > 2000 {
		lines := strings.Split(text, "\n")
		text = ""
		for _, line := range lines {
			if len(text)+len(line) < 2000 {
				text += line + "\n"
			}
		}
	}

	return &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Hints for %s", el.Name),
		Description: text,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: el.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%d Hints • You %sHave This", len(out), txt),
		},
	}, "", true
}

func getHintText(elems []string, inv map[string]empty, dat serverData) (string, int) {
	hasElems := true
	for _, val := range elems {
		_, exists := inv[strings.ToLower(val)]
		if !exists {
			hasElems = false
		}
	}
	pref := x
	ex := 0
	if hasElems {
		pref = check
		ex = 1
	}
	prf := "%s"
	params := make([]interface{}, len(elems))
	i := 0
	for _, k := range elems {
		dat.lock.RLock()
		params[i] = interface{}(dat.elemCache[strings.ToLower(k)].Name)
		dat.lock.RUnlock()

		if i == 0 {
			prf += " %s"
		} else {
			prf += " + %s"
		}
		i++
	}

	params = append([]interface{}{pref}, params...)
	params[len(params)-1] = obscure(params[len(params)-1].(string))
	txt := fmt.Sprintf(prf, params...)
	return txt, ex
}
