package eod

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const pageLength = 10
const leftArrow = "⬅️"
const rightArrow = "➡️"

const ldbQuery = `
SELECT rw, ` + "user" + `, ` + "count" + `
FROM (
    SELECT 
         ROW_NUMBER() OVER (ORDER BY ` + "count" + ` DESC) AS rw,
         ` + "user" + `, ` + "count" + `
    FROM eod_inv WHERE guild=?
) sub
WHERE sub.user=?
`

func (b *EoD) invPageGetter(p pageSwitcher) (string, int, int, error) {
	length := (len(p.Items) - 1) / pageLength
	if pageLength*p.Page > len(p.Items) {
		return "", 0, length, nil
	}

	if p.Page < 0 {
		return "", length, length, nil
	}

	items := p.Items[pageLength*p.Page:]
	if len(items) > pageLength {
		items = items[:pageLength]
	}
	return strings.Join(items, "\n"), p.Page, length, nil
}

func (b *EoD) ldbPageGetter(p pageSwitcher) (string, int, int, error) {
	cnt := b.db.QueryRow("SELECT COUNT(1) FROM eod_inv WHERE guild=?", p.Guild)
	pos := b.db.QueryRow(ldbQuery, p.Guild, p.User)
	var count int
	var ps int
	var u string
	var ucnt int
	err := pos.Scan(&ps, &u, &ucnt)
	if err != nil {
		return "", 0, 0, err
	}
	cnt.Scan(&count)
	length := (count - 1) / pageLength
	if err != nil {
		return "", 0, 0, err
	}
	if pageLength*p.Page > count {
		return "", 0, length, nil
	}

	if p.Page < 0 {
		return "", length, length, nil
	}

	text := ""
	res, err := b.db.Query("SELECT `count`, `user` FROM eod_inv WHERE guild=? ORDER BY `count` DESC LIMIT ? OFFSET ?", p.Guild, pageLength, p.Page*pageLength)
	if err != nil {
		return "", 0, 0, err
	}
	defer res.Close()
	i := pageLength*p.Page + 1
	var user string
	var ct int
	for res.Next() {
		err = res.Scan(&ct, &user)
		if err != nil {
			return "", 0, 0, err
		}
		text += fmt.Sprintf("%d. <@%s> - %d\n", i, user, ct)
		i++
	}
	if !((pageLength*p.Page <= (ps - 1)) && ((ps - 1) <= (p.Page+1)*pageLength)) {
		text += fmt.Sprintf("%d. <@%s> - %d\n", ps, u, ucnt)
	}
	return text, p.Page, length, nil
}

func (b *EoD) newPageSwitcher(ps pageSwitcher, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	ps.Channel = m.ChannelID
	ps.Guild = m.GuildID
	ps.Page = 0

	cont, _, length, err := ps.PageGetter(ps)
	if rsp.Error(err) {
		return
	}
	rsp.BlankReply()
	msg, _ := b.dg.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       ps.Title,
		Description: cont,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Page %d/%d", ps.Page+1, length+1),
		},
	})
	id := msg.ID
	b.dg.MessageReactionAdd(m.ChannelID, id, leftArrow)
	b.dg.MessageReactionAdd(m.ChannelID, id, rightArrow)
	if dat.pageSwitchers == nil {
		dat.pageSwitchers = make(map[string]pageSwitcher)
	}
	dat.pageSwitchers[id] = ps

	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()
}

func (b *EoD) pageSwitchHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == b.dg.State.User.ID {
		return
	}

	lock.RLock()
	dat, exists := b.dat[r.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	ps, exists := dat.pageSwitchers[r.MessageID]
	if !exists {
		return
	}

	if r.Emoji.Name == rightArrow {
		ps.Page++
	} else if r.Emoji.Name == leftArrow {
		ps.Page--
	} else {
		return
	}

	cont, page, length, err := ps.PageGetter(ps)
	if err != nil {
		return
	}
	if page != ps.Page {
		ps.Page = page
		cont, _, length, err = ps.PageGetter(ps)
		if err != nil {
			return
		}
	}

	b.dg.ChannelMessageEditEmbed(ps.Channel, r.MessageID, &discordgo.MessageEmbed{
		Title:       ps.Title,
		Description: cont,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Page %d/%d", ps.Page+1, length+1),
		},
	})
	b.dg.MessageReactionRemove(ps.Channel, r.MessageID, r.Emoji.Name, r.UserID)
	dat.pageSwitchers[r.MessageID] = ps

	lock.Lock()
	b.dat[r.GuildID] = dat
	lock.Unlock()
}

func (b *EoD) invCmd(m msg, rsp rsp) {
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
	items := make([]string, len(inv))
	i := 0
	for k := range inv {
		items[i] = dat.elemCache[k].Name
		i++
	}

	b.newPageSwitcher(pageSwitcher{
		Kind:       pageSwitchInv,
		Title:      m.Author.Username + "'s Inventory",
		PageGetter: b.invPageGetter,
		Items:      items,
	}, m, rsp)
}

func (b *EoD) ldbCmd(m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	_, exists = dat.invCache[m.Author.ID]
	if !exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}

	b.newPageSwitcher(pageSwitcher{
		Kind:       pageSwitchLdb,
		Title:      "Top Most Elements",
		PageGetter: b.ldbPageGetter,
		User:       m.Author.ID,
	}, m, rsp)
}
