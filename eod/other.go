package eod

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

func (b *EoD) giveAllCmd(user string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	dat.Lock.RLock()
	inv, exists := dat.InvCache[user]
	dat.Lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}
	dat.Lock.RLock()
	for k := range dat.ElemCache {
		inv[k] = types.Empty{}
	}
	dat.Lock.RUnlock()

	dat.Lock.Lock()
	dat.InvCache[user] = inv
	dat.Lock.Unlock()

	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()
	b.saveInv(m.GuildID, user, true, true)
	rsp.Resp("Successfully gave every element to <@" + user + ">!")
}

func (b *EoD) resetInvCmd(user string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	inv := make(map[string]types.Empty)
	for _, v := range starterElements {
		inv[strings.ToLower(v.Name)] = types.Empty{}
	}

	dat.Lock.Lock()
	dat.InvCache[user] = inv
	dat.Lock.Unlock()

	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()
	b.saveInv(m.GuildID, user, true, true)
	rsp.Resp("Successfully reset <@" + user + ">'s inventory!")
}

func (b *EoD) downloadInvCmd(user string, sorter string, m types.Msg, rsp types.Rsp) {
	rsp.Acknowledge()

	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	inv, exists := dat.InvCache[user]
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
	dat.Lock.RLock()
	for k := range inv {
		items[i] = dat.ElemCache[k].Name
		i++
	}
	dat.Lock.RUnlock()

	switch sorter {
	case "id":
		sort.Slice(items, func(i, j int) bool {
			dat.Lock.RLock()
			elem1, exists := dat.ElemCache[strings.ToLower(items[i])]
			if !exists {
				return false
			}

			elem2, exists := dat.ElemCache[strings.ToLower(items[j])]
			if !exists {
				return false
			}
			dat.Lock.RUnlock()
			return elem1.CreatedOn.Before(elem2.CreatedOn)
		})

	case "madeby":
		count := 0
		outs := make([]string, len(items))
		for _, val := range items {
			creator := ""
			dat.Lock.RLock()
			elem, exists := dat.ElemCache[strings.ToLower(val)]
			dat.Lock.RUnlock()
			if exists {
				creator = elem.Creator
			}
			if creator == user {
				outs[count] = val
				count++
			}
		}
		outs = outs[:count]
		sortStrings(outs)
		items = outs

	case "length":
		sort.Slice(items, func(i, j int) bool {
			return len(items[i]) < len(items[j])
		})

	default:
		sortStrings(items)
	}

	out := &strings.Builder{}
	for _, val := range items {
		out.WriteString(val + "\n")
	}
	buf := strings.NewReader(out.String())

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
