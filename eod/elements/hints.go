package elements

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
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
	b *Elements
}

func (h *hintComponent) Handler(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	m := types.Msg{
		Author:    i.Member.User,
		ChannelID: i.ChannelID,
		GuildID:   i.GuildID,
	}
	hint, msg, suc := h.b.getHint("", false, i.Member.User.ID, i.GuildID, false, m, nil)
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

func (b *Elements) HintCmd(elem string, hasElem bool, inverse bool, m types.Msg, rsp types.Rsp) {
	rsp.Acknowledge()
	elem = strings.TrimSpace(elem)
	elem = util.EscapeElement(elem)

	rspInp := rsp
	if !hasElem {
		if inverse {
			rsp.ErrorMessage("You cannot have an inverse hint without an element!")
			return
		}
		rspInp = nil
	}
	hint, msg, suc := b.getHint(elem, hasElem, m.Author.ID, m.GuildID, inverse, m, rspInp)
	if !suc && msg == "" {
		return
	}

	if !suc {
		rsp.ErrorMessage(msg)
		return
	}

	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		return
	}

	if hasElem {
		id := rsp.Embed(hint)
		dat.SetMsgElem(id, elem)
		return
	}

	id := rsp.Embed(hint, hintCmp)

	dat.AddComponentMsg(id, &hintComponent{b: b})

	b.lock.Lock()
	b.dat[m.GuildID] = dat
	b.lock.Unlock()
}

func (b *Elements) getHint(elem string, hasElem bool, author string, guild string, inverse bool, m types.Msg, rsp types.Rsp) (*discordgo.MessageEmbed, string, bool) {
	b.lock.RLock()
	dat, exists := b.dat[guild]
	b.lock.RUnlock()
	if !exists {
		return nil, "Guild not found", false
	}

	inv, res := dat.GetInv(author, true)
	if !exists {
		return nil, res.Message, false
	}
	var el types.Element
	if hasElem {
		el, res = dat.GetElement(elem)
		if !res.Exists {
			return nil, fmt.Sprintf("No hints were found for **%s**!", elem), false
		}
	}
	if !hasElem {
		hasFound := false
		dat.Lock.RLock()
		for _, v := range dat.Elements {
			_, exists := inv[strings.ToLower(v.Name)]
			if !exists {
				el = v
				elem = v.Name
				hasFound = true
				break
			}
		}
		dat.Lock.RUnlock()
		if !hasFound {
			dat.Lock.RLock()
			for _, v := range dat.Elements {
				el = v
				elem = v.Name
				hasFound = true
				break
			}
			dat.Lock.RUnlock()
		}
	}

	vals := make(map[string]types.Empty)
	if !inverse {
		dat.Lock.RLock()
		for elems, elem3 := range dat.Combos {
			if strings.EqualFold(elem3, elem) {
				vals[elems] = types.Empty{}
			}
		}
		dat.Lock.RUnlock()
	} else {
		for elems, elem3 := range dat.Combos {
			parts := strings.Split(elems, "+")
			for _, part := range parts {
				if strings.EqualFold(part, elem) {
					vals[elem3] = types.Empty{}
					break
				}
			}
		}
	}

	out := make([]hintCombo, len(vals))
	length := 0
	i := 0
	for val := range vals {
		txt, ex := getHintText(val, inv, dat, inverse)
		out[i] = hintCombo{
			exists: ex,
			text:   txt,
		}

		length += len(txt)
		i++
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].exists > out[j].exists
	})

	inverseTitle := ""
	if inverse {
		inverseTitle = "Inverse "
	}
	title := fmt.Sprintf("%sHints for %s", inverseTitle, el.Name)

	text := &strings.Builder{}
	for _, val := range out {
		text.WriteString(val.text)
		text.WriteString("\n")
	}
	val := text.String()

	txt := "Don't "
	_, hasElem = inv[strings.ToLower(el.Name)]
	if hasElem {
		txt = ""
	}
	footer := fmt.Sprintf("%d Hints • You %sHave This", len(out), txt)

	isPlayChannel := dat.PlayChannels.Contains(m.ChannelID)

	if len(val) > 2000 && rsp == nil {
		// If can't do page switcher, shorten it (cant do pageswitcher if its in a random hint)
		text = &strings.Builder{}
		length := 0
		for _, val := range out {
			if length+len(val.text) > 2000 {
				break
			}

			length += len(val.text)
			text.WriteString(val.text)
		}
		val = text.String()
	} else if (len(val) > 2000) || (!isPlayChannel && len(out) > base.DefaultPageLength) || (isPlayChannel && len(out) > base.PlayPageLength) {
		vals := make([]string, len(out))
		for i, v := range out {
			vals[i] = v.text
		}
		b.base.NewPageSwitcher(types.PageSwitcher{
			Kind:       types.PageSwitchInv,
			Title:      title,
			PageGetter: b.base.InvPageGetter,
			Items:      vals,
			Thumbnail:  el.Image,
			Footer:     footer,
		}, m, rsp)
		return nil, "", false
	}

	return &discordgo.MessageEmbed{
		Title:       title,
		Description: val,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: el.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: footer,
		},
	}, "", true
}

func getHintText(elemTxt string, inv types.Container, dat types.ServerData, inverse bool) (string, int) {
	if !inverse {
		elems := strings.Split(elemTxt, "+")
		hasElems := true
		for _, val := range elems {
			_, exists := inv[strings.ToLower(val)]
			if !exists {
				hasElems = false
			}
		}
		pref := types.X
		ex := 0
		if hasElems {
			pref = types.Check
			ex = 1
		}
		prf := "%s"
		params := make([]interface{}, len(elems))
		i := 0
		dat.Lock.RLock()
		for _, k := range elems {
			elem, _ := dat.GetElement(k, true)
			params[i] = interface{}(elem.Name)

			if i == 0 {
				prf += " %s"
			} else {
				prf += " + %s"
			}
			i++
		}
		dat.Lock.RUnlock()

		params = append([]interface{}{pref}, params...)
		params[len(params)-1] = util.Obscure(params[len(params)-1].(string))
		txt := fmt.Sprintf(prf, params...)
		return txt, ex
	}

	found := inv.Contains(elemTxt)
	txt := types.X
	ex := 0
	if found {
		txt = types.Check
		ex = 1
	}
	el, _ := dat.GetElement(elemTxt)
	txt += " " + el.Name
	return txt, ex
}
