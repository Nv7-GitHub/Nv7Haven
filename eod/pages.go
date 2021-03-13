package eod

import (
	"fmt"
	"math"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const pageLength = 10
const leftArrow = "⬅️"
const rightArrow = "➡️"

const ldbQuery = `
SELECT rw
FROM (
    SELECT 
         ROW_NUMBER() OVER (ORDER BY ` + "count" + ` DESC) AS rw,
         ` + "user" + `
    FROM eod_inv WHERE guild=?
) sub
WHERE sub.user=?
`

func (b *EoD) invPageGetter(p pageSwitcher) (string, int, int, error) {
	length := len(p.Items) / pageLength
	if pageLength*p.Page > len(p.Items) {
		return "", 0, length, nil
	}

	if p.Page < 0 {
		return "", int(math.Floor(float64(p.Page*pageLength) / float64(len(p.Items)))), length, nil
	}

	items := p.Items[pageLength*p.Page:]
	if len(items) > pageLength {
		items = items[:pageLength]
	}
	return strings.Join(items, "\n"), p.Page, length, nil
}

func (b *EoD) newPageSwitcher(ps pageSwitcher, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	cont, _, length, err := ps.PageGetter(ps)
	if rsp.Error(err) {
		return
	}
	id := rsp.Embed(&discordgo.MessageEmbed{
		Title:       ps.Title,
		Description: cont,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Page %d/%d", ps.Page, length),
		},
	})
	b.dg.MessageReactionAdd(m.ChannelID, id, leftArrow)
	b.dg.MessageReactionAdd(m.ChannelID, id, rightArrow)
	ps.Channel = m.ChannelID
	ps.Guild = m.GuildID
	ps.Page = 0
	dat.pageSwitchers[id] = ps

	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()
}

func (b *EoD) pageSwitchHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
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
			Text: fmt.Sprintf("Page %d/%d", ps.Page, length),
		},
	})
	b.dg.MessageReactionsRemoveEmoji(ps.Channel, r.MessageID, r.Emoji.Name)
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
		items[i] = k
		i++
	}

	b.newPageSwitcher(pageSwitcher{
		Kind:       pageSwitchInv,
		Title:      m.Author.Username + "'s Inventory",
		PageGetter: b.invPageGetter,
		Items:      items,
	}, m, rsp)
}
