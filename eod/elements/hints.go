package elements

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/eodsort"
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
			Emoji: discordgo.ComponentEmoji{
				Name:     "hint",
				ID:       "932833490066620457",
				Animated: false,
			},
		},
	},
}

type hintComponent struct {
	b  *Elements
	db *eodb.DB
}

func (h *hintComponent) Handler(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	m := types.Msg{
		Author:    i.Member.User,
		ChannelID: i.ChannelID,
		GuildID:   i.GuildID,
	}
	hint, msg, suc := h.b.getHint(0, h.db, false, i.Member.User.ID, i.GuildID, false, m, nil)
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
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	elem = strings.TrimSpace(elem)
	elem = util.EscapeElement(elem)
	elId := 0
	if hasElem {
		el, res := db.GetElementByName(elem)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
		elId = el.ID
	}

	rspInp := rsp
	if !hasElem {
		if inverse {
			rsp.ErrorMessage("You cannot have an inverse hint without an element!")
			return
		}
		rspInp = nil
	}
	hint, msg, suc := b.getHint(elId, db, hasElem, m.Author.ID, m.GuildID, inverse, m, rspInp)
	if !suc && msg == "" {
		return
	}

	if !suc {
		rsp.ErrorMessage(msg)
		return
	}

	data, _ := b.GetData(m.GuildID)

	if hasElem {
		id := rsp.Embed(hint)
		data.SetMsgElem(id, elId)
		return
	}

	id := rsp.Embed(hint, hintCmp)

	data.AddComponentMsg(id, &hintComponent{b: b, db: db})
}

func (b *Elements) getHint(elem int, db *eodb.DB, hasElem bool, author string, guild string, inverse bool, m types.Msg, rsp types.Rsp) (*discordgo.MessageEmbed, string, bool) {
	rand.Seed(time.Now().UnixNano())
	inv := db.GetInv(author)
	var el types.Element
	if !hasElem {
		hasFound := false
		ids := make([]int, len(db.Elements))
		db.RLock()
		for i, v := range db.Elements {
			ids[i] = v.ID
		}
		db.RUnlock()

		// Shuffle ids
		rand.Shuffle(len(ids), func(i, j int) {
			v := ids[i]
			ids[i] = ids[j]
			ids[j] = v
		})
		for _, id := range ids {
			exists := inv.Contains(id)
			if !exists {
				el, _ = db.GetElement(id)
				hasFound = true
				break
			}
		}
		if !hasFound {
			db.RLock()
			id := rand.Intn(len(db.Elements)-1) + 1
			el, _ = db.GetElement(id, true)
			db.RUnlock()
		}
	} else {
		el, _ = db.GetElement(elem)
	}

	vals := make(map[string]types.Empty)
	if !inverse {
		db.RLock()
		for elems, elem3 := range db.Combos() {
			if elem3 == el.ID {
				vals[elems] = types.Empty{}
			}
		}
		db.RUnlock()
	} else {
		db.RLock()
		for elems, elem3 := range db.Combos() {
			parts := strings.Split(elems, "+")
			for _, part := range parts {
				num, err := strconv.Atoi(part)
				if err != nil {
					continue
				}
				if num == el.ID {
					el, _ := db.GetElement(elem3)
					vals[el.Name] = types.Empty{}
					break
				}
			}
		}
		db.RUnlock()
	}

	out := make([]hintCombo, len(vals))
	length := 0
	i := 0
	for val := range vals {
		txt, ex := getHintText(val, inv, db, inverse)
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
	hasElem = inv.Contains(el.ID)
	if hasElem {
		txt = ""
	}
	footer := fmt.Sprintf("%d Hints • You %sHave This", len(out), txt)

	db.Config.RLock()
	isPlayChannel := db.Config.PlayChannels.Contains(m.ChannelID)
	db.Config.RUnlock()

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

func getHintText(elemTxt string, inv *types.Inventory, db *eodb.DB, inverse bool) (string, int) {
	if !inverse {
		elemDat := strings.Split(elemTxt, "+")
		elems := make([]string, len(elemDat))
		hasElems := true
		for i, val := range elemDat {
			num, err := strconv.Atoi(val)
			if err != nil {
				hasElems = false
				continue
			}
			el, res := db.GetElement(num)
			if res.Exists {
				exists := inv.Contains(el.ID)
				if !exists {
					hasElems = false
				}
			} else {
				hasElems = false
			}
			elems[i] = el.Name
		}
		sort.Slice(elems, func(i, j int) bool {
			return eodsort.CompareStrings(elems[i], elems[j])
		})
		pref := types.X
		ex := 0
		if hasElems {
			pref = types.Check
			ex = 1
		}
		prf := "%s"
		params := make([]interface{}, len(elems))
		i := 0
		db.RLock()
		for _, k := range elems {
			elem, _ := db.GetElementByName(k, true)
			params[i] = interface{}(elem.Name)

			if i == 0 {
				prf += " %s"
			} else {
				prf += " + %s"
			}
			i++
		}
		db.RUnlock()

		params = append([]interface{}{pref}, params...)
		params[len(params)-1] = util.Obscure(params[len(params)-1].(string))
		txt := fmt.Sprintf(prf, params...)
		return txt, ex
	}

	el, _ := db.GetElementByName(elemTxt)
	found := inv.Contains(el.ID)
	txt := types.X
	ex := 0
	if found {
		txt = types.Check
		ex = 1
	}
	txt += " " + el.Name
	return txt, ex
}
