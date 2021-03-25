package eod

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type hintCombo struct {
	exists int
	text   string
}

func obscure(val string) string {
	question := []byte("?")[0]
	out := make([]byte, len(val))
	for i := range out {
		out[i] = question
	}
	return string(out)
}

func (b *EoD) hintCmd(elem string, hasElem bool, m msg, rsp rsp) {
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
		for _, v := range dat.elemCache {
			el = v
			elem = v.Name
			break
		}
	}

	var combs *sql.Rows
	var err error
	combs, err = b.db.Query("SELECT elems FROM eod_combos WHERE elem3=? AND guild=?", elem, m.GuildID)
	if rsp.Error(err) {
		return
	}
	defer combs.Close()
	var elemDat string
	var elemCombos map[string]empty
	out := make([]hintCombo, 0)
	for combs.Next() {
		err = combs.Scan(&elemDat)
		if rsp.Error(err) {
			return
		}
		err = json.Unmarshal([]byte(elemDat), &elemCombos)
		if rsp.Error(err) {
			return
		}

		hasElems := true
		for k := range elemCombos {
			_, exists := inv[strings.ToLower(k)]
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
		prf := "%s %s + %s"
		params := make([]interface{}, len(elemCombos))
		i := 0
		for k := range elemCombos {
			params[i] = interface{}(k)
			prf += " + %s"
			i++
		}
		params = append([]interface{}{pref}, params)
		params[len(params)-1] = obscure(params[len(params)-1].(string))
		txt := fmt.Sprintf(prf, params...)
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
	rsp.Resp(fmt.Sprintf("Element Count: %d\nCombination Count: %d\nMember Count: %d", len(dat.elemCache), cnt, gd.MemberCount))
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
	b.saveInv(m.GuildID, user)
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
	b.saveInv(m.GuildID, user)
	rsp.Resp("Successfully reset <@" + user + ">'s inventory!")
}
