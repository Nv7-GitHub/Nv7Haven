package eod

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type hintCombo struct {
	exists int
	text   string
}

var noObscure = map[byte]empty{
	' ': {},
	'.': {},
	'-': {},
	'_': {},
}

func obscure(val string) string {
	question := []byte("?")[0]
	out := make([]byte, len(val))
	for i, char := range []byte(val) {
		_, exists := noObscure[char]
		if exists {
			out[i] = char
		} else {
			out[i] = question
		}
	}
	return string(out)
}

func (b *EoD) hintCmd(elem string, hasElem bool, m msg, rsp rsp) {
	rsp.Acknowledge()

	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	inv, exists := dat.invCache[m.Author.ID]
	if !exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}
	var el element
	if hasElem {
		el, exists = dat.elemCache[strings.ToLower(elem)]
		if !exists {
			rsp.ErrorMessage(fmt.Sprintf("No hints were found for **%s**!", elem))
			return
		}
	}
	if !hasElem {
		hasFound := false
		for _, v := range dat.elemCache {
			_, exists := inv[strings.ToLower(v.Name)]
			if !exists {
				el = v
				elem = v.Name
				hasFound = true
				break
			}
		}
		if !hasFound {
			for _, v := range dat.elemCache {
				el = v
				elem = v.Name
				hasFound = true
				break
			}
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

	combs, err = b.db.Query(query, elem, m.GuildID)
	if rsp.Error(err) {
		return
	}
	defer combs.Close()
	out := make([]hintCombo, 0)
	var elemTxt string

	length := 0
	for combs.Next() {
		err = combs.Scan(&elemTxt)
		if rsp.Error(err) {
			return
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
		element := dat.elemCache[strings.ToLower(elem)]

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

	rsp.Embed(&discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Hints for %s", el.Name),
		Description: text,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: el.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("%d Hints • You %sHave This", len(out), txt),
		},
	})
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
		params[i] = interface{}(dat.elemCache[strings.ToLower(k)].Name)
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

func (b *EoD) statsCmd(m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	gd, err := b.dg.State.Guild(m.GuildID)
	if rsp.Error(err) {
		return
	}
	var cnt int
	row := b.db.QueryRow("SELECT COUNT(1) FROM eod_combos WHERE guild=?", m.GuildID)
	err = row.Scan(&cnt)
	if rsp.Error(err) {
		return
	}
	found := 0
	for _, val := range dat.invCache {
		found += len(val)
	}
	rsp.Message(fmt.Sprintf("Element Count: **%s**\nCombination Count: **%s**\nMember Count: **%s**\nElements Found: **%s**", formatInt(len(dat.elemCache)), formatInt(cnt), formatInt(gd.MemberCount), formatInt(found)))
}

func (b *EoD) giveAllCmd(user string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	inv, exists := dat.invCache[user]
	if !exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}
	for k := range dat.elemCache {
		inv[k] = empty{}
	}
	dat.invCache[user] = inv

	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()
	b.saveInv(m.GuildID, user, true, true)
	rsp.Resp("Successfully gave every element to <@" + user + ">!")
}

func (b *EoD) resetInvCmd(user string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	inv := make(map[string]empty)
	for _, v := range starterElements {
		inv[strings.ToLower(v.Name)] = empty{}
	}
	dat.invCache[user] = inv

	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()
	b.saveInv(m.GuildID, user, true, true)
	rsp.Resp("Successfully reset <@" + user + ">'s inventory!")
}

func (b *EoD) downloadInvCmd(user string, sorter string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	inv, exists := dat.invCache[user]
	if !exists {
		if user == m.Author.ID {
			rsp.ErrorMessage("You don't have an inventory!")
		} else {
			rsp.ErrorMessage(fmt.Sprintf("User <@%s> doesn't have an inventory!", user))
		}
		return
	}
	items := make([]string, len(inv))
	i := 0
	for k := range inv {
		items[i] = dat.elemCache[k].Name
		i++
	}

	switch sorter {
	case "id":
		sort.Slice(items, func(i, j int) bool {
			elem1, exists := dat.elemCache[strings.ToLower(items[i])]
			if !exists {
				return false
			}

			elem2, exists := dat.elemCache[strings.ToLower(items[j])]
			if !exists {
				return false
			}
			return elem1.CreatedOn.Before(elem2.CreatedOn)
		})

	case "madeby":
		count := 0
		outs := make([]string, len(items))
		for _, val := range items {
			creator := ""
			elem, exists := dat.elemCache[strings.ToLower(val)]
			if exists {
				creator = elem.Creator
			}
			if creator == user {
				outs[count] = val
				count++
			}
		}
		outs = outs[:count]
		sort.Strings(outs)
		items = outs

	default:
		sort.Strings(items)
	}

	txt := strings.Join(items, "\n")
	buf := strings.NewReader(txt)

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}

	usr, err := b.dg.User(user)
	if rsp.Error(err) {
		return
	}
	gld, err := b.dg.Guild(m.GuildID)
	if rsp.Error(err) {
		return
	}

	b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Inv for **%s** in **%s**:", usr.Username, gld.Name),
		Files: []*discordgo.File{
			{
				Name:        "inv.txt",
				ContentType: "text/plain",
				Reader:      buf,
			},
		},
	})
	rsp.Message("Sent inv in DMs!")
}
