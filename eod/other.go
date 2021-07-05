package eod

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

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

	categorized := 0
	for _, val := range dat.catCache {
		categorized += len(val.Elements)
	}

	dat.lock.RLock()
	rsp.Message(fmt.Sprintf("Element Count: **%s**\nCombination Count: **%s**\nMember Count: **%s**\nElements Found: **%s**\nElements Categorized: **%s**", formatInt(len(dat.elemCache)), formatInt(cnt), formatInt(gd.MemberCount), formatInt(found), formatInt(categorized)))
	dat.lock.RUnlock()
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
	dat.lock.RLock()
	for k := range dat.elemCache {
		inv[k] = empty{}
	}
	dat.lock.RUnlock()
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
	dat.lock.RLock()
	for k := range inv {
		items[i] = dat.elemCache[k].Name
		i++
	}
	dat.lock.RUnlock()

	switch sorter {
	case "id":
		sort.Slice(items, func(i, j int) bool {
			dat.lock.RLock()
			elem1, exists := dat.elemCache[strings.ToLower(items[i])]
			if !exists {
				return false
			}

			elem2, exists := dat.elemCache[strings.ToLower(items[j])]
			if !exists {
				return false
			}
			dat.lock.RUnlock()
			return elem1.CreatedOn.Before(elem2.CreatedOn)
		})

	case "madeby":
		count := 0
		outs := make([]string, len(items))
		for _, val := range items {
			creator := ""
			dat.lock.RLock()
			elem, exists := dat.elemCache[strings.ToLower(val)]
			dat.lock.RUnlock()
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
		Content: fmt.Sprintf("**%s**'s Inventory in **%s**:", usr.Username, gld.Name),
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
