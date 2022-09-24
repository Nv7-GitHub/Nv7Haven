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

func newHintCmp(db *eodb.DB) discordgo.ActionsRow {
	return discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    db.Config.LangProperty("NewHint", nil),
				CustomID: "hint-new",
				Style:    discordgo.SuccessButton,
				Emoji: discordgo.ComponentEmoji{
					Name:     "hint",
					ID:       "932833472396025908",
					Animated: false,
				},
			},
		},
	}
}

type hintComponent struct {
	b       *Elements
	db      *eodb.DB
	hasCat  bool
	catName string
}

func (h *hintComponent) Handler(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	m := types.Msg{
		Author:    i.Member.User,
		ChannelID: i.ChannelID,
		GuildID:   i.GuildID,
	}
	hint, msg, suc := h.b.getHint(0, h.db, false, h.hasCat, h.catName, i.Member.User.ID, i.GuildID, false, m, nil)
	if !suc {
		h.b.dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content:    msg,
				Components: []discordgo.MessageComponent{newHintCmp(h.db)},
			},
		})
		return
	}

	h.b.dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{hint},
			Components: []discordgo.MessageComponent{newHintCmp(h.db)},
		},
	})
}

type hintCombo struct {
	exists int
	text   string
}

func (b *Elements) HintCmd(elem string, hasElem bool, hasCat bool, inverse bool, m types.Msg, rsp types.Rsp) {
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

	catName := ""
	if hasCat {
		cat, res := db.GetCat(elem)
		if !res.Exists {
			vcat, res := db.GetVCat(elem)
			if !res.Exists {
				rsp.ErrorMessage(res.Message)
				return
			} else {
				catName = vcat.Name
			}
		} else {
			catName = cat.Name
		}
	}

	rspInp := rsp
	if !hasElem {
		if inverse {
			rsp.ErrorMessage(db.Config.LangProperty("InvHintNoElem", nil))
			return
		}
		rspInp = nil
	}
	hint, msg, suc := b.getHint(elId, db, hasElem, hasCat, catName, m.Author.ID, m.GuildID, inverse, m, rspInp)
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

	id := rsp.Embed(hint, newHintCmp(db))

	data.AddComponentMsg(id, &hintComponent{b: b, db: db, hasCat: hasCat, catName: catName})
}

func (b *Elements) getHint(elem int, db *eodb.DB, hasElem bool, hasCat bool, catName string, author string, guild string, inverse bool, m types.Msg, rsp types.Rsp) (*discordgo.MessageEmbed, string, bool) {
	rand.Seed(time.Now().UnixNano())
	inv := db.GetInv(author)
	var el types.Element
	if !hasElem {
		hasFound := false
		ids := make([]int, len(db.Elements))
		if !hasCat { // Use all elements if no cat
			db.RLock()
			for i, v := range db.Elements {
				ids[i] = v.ID
			}
			db.RUnlock()
		} else {
			cat, res := db.GetCat(catName) // Can ignore errors since checked in HintCmd
			if !res.Exists {
				vcat, _ := db.GetVCat(catName)
				els, res := b.base.CalcVCat(vcat, db, true)
				if !res.Exists {
					return nil, res.Message, false
				}
				ids = make([]int, len(els))
				i := 0
				for v := range els {
					ids[i] = v
					i++
				}
			} else {
				i := 0
				ids = make([]int, len(cat.Elements))
				cat.Lock.RLock()
				for v := range cat.Elements {
					ids[i] = v
					i++
				}
				cat.Lock.RUnlock()
			}
		}

		if len(ids) == 0 {
			return nil, "eod: cannot calculate hint", false
		}

		// Shuffle ids
		rand.Shuffle(len(ids), func(i, j int) {
			v := ids[i]
			ids[i] = ids[j]
			ids[j] = v
		})
		for _, id := range ids {
			exists := inv.Contains(id)
			if !exists {
				// Check if you can make the element
				db.RLock()
				canMake := false
				for elems, elem3 := range db.Combos() {
					if elem3 == id {
						ok := true
						parts := strings.Split(elems, "+")
						for _, part := range parts {
							v, _ := strconv.Atoi(part)
							if !inv.Contains(v) {
								ok = false
								break
							}
						}

						if ok {
							canMake = true
							break
						}
					}
				}
				db.RUnlock()
				if canMake {
					el, _ = db.GetElement(id)
					hasFound = true
					break
				}
			}
		}
		if !hasFound {
			db.RLock()
			ind := rand.Intn(len(ids)) + 1
			el, _ = db.GetElement(ids[ind], true)
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
					el, _ := db.GetElement(elem3, true)
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

	title := db.Config.LangProperty("HintElem", el.Name)
	if inverse {
		title = db.Config.LangProperty("InvHintElem", el.Name)
	}

	text := &strings.Builder{}
	for _, val := range out {
		text.WriteString(val.text)
		text.WriteString("\n")
	}
	val := text.String()

	footer := db.Config.LangProperty("HintCountNoHasElem", len(out))
	hasElem = inv.Contains(el.ID)
	if hasElem {
		footer = db.Config.LangProperty("HintCountHasElem", len(out))
	}

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
		params := make([]any, len(elems))
		i := 0
		db.RLock()
		for _, k := range elems {
			elem, _ := db.GetElementByName(k, true)
			params[i] = any(elem.Name)

			if i == 0 {
				prf += " %s"
			} else {
				prf += " + %s"
			}
			i++
		}
		db.RUnlock()

		params = append([]any{pref}, params...)
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
