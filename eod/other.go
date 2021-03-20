package eod

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// guild, guild, user, guild
const ideaQuery = `SELECT ` + "eod_elements" + `.name, nm2.name FROM ` + "eod_elements" + `, (SELECT name FROM eod_elements WHERE guild=? ORDER BY RAND() LIMIT 1) nm2, (SELECT inv FROM eod_inv WHERE guild=? AND ` + "user" + `=?) inv WHERE (SELECT COUNT(1) FROM eod_combos WHERE (elem1=` + "eod_elements" + `.name AND elem2=nm2.name) OR (elem1=nm2.name AND elem2=` + "eod_elements" + `.name))=0 AND guild=? AND (JSON_EXTRACT(inv.inv, CONCAT('$."', LOWER(nm2.name), '"')) IS NOT NULL) AND (JSON_EXTRACT(inv.inv, CONCAT('$."', LOWER(` + "eod_elements" + `.name), '"')) IS NOT NULL) ORDER BY RAND() LIMIT 1`

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
	combs, err = b.db.Query("SELECT elem1, elem2 FROM eod_combos WHERE elem3=? AND guild=?", elem, m.GuildID)
	if rsp.Error(err) {
		return
	}
	defer combs.Close()
	var elem1 string
	var elem2 string
	out := make([]hintCombo, 0)
	for combs.Next() {
		err = combs.Scan(&elem1, &elem2)
		if rsp.Error(err) {
			return
		}

		_, haselem1 := inv[strings.ToLower(elem1)]
		_, haselem2 := inv[strings.ToLower(elem2)]
		pref := x
		ex := 0
		if haselem1 && haselem2 {
			pref = check
			ex = 1
		}
		txt := fmt.Sprintf("%s %s + %s", pref, elem1, obscure(elem2))
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

func (b *EoD) ideaCmd(m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	var elem1 string
	var elem2 string
	var count int
	for i := 0; i < 10; i++ {
		row := b.db.QueryRow(ideaQuery, m.GuildID, m.GuildID, m.Author.ID, m.GuildID)
		err := row.Scan(&elem1, &elem2)
		if err != nil {
			continue
		}
		el1, exists := dat.elemCache[strings.ToLower(elem1)]
		if !exists {
			continue
		}
		el2, exists := dat.elemCache[strings.ToLower(elem2)]
		if !exists {
			continue
		}
		row = b.db.QueryRow("SELECT COUNT(1) FROM eod_combos WHERE guild=? AND ((elem1=? AND elem2=?) OR (elem1=? AND elem2=?))", m.GuildID, elem1, elem2, elem2, elem1)
		err = row.Scan(&count)
		if err != nil {
			continue
		}
		if count != 0 {
			continue
		}
		if dat.combCache == nil {
			dat.combCache = make(map[string]comb)
		}
		dat.combCache[m.Author.ID] = comb{
			elem1: el1.Name,
			elem2: el2.Name,
			elem3: "",
		}
		lock.Lock()
		b.dat[m.GuildID] = dat
		lock.Unlock()

		rsp.Resp(fmt.Sprintf("Your random unused combination is... **%s** + **%s**\n 	Suggest it by typing **/suggest**", el1.Name, el2.Name))
		return
	}
	rsp.ErrorMessage("Couldn't find a random unused combo! Maybe try again?")
}
