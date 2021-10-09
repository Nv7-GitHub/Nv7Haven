package eod

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/bwmarrin/discordgo"
)

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

	dat.SetInv(user, inv)

	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()
	b.base.SaveInv(m.GuildID, user, true, true)
	rsp.Resp("Successfully reset <@" + user + ">'s inventory!")
}

func (b *EoD) downloadInvCmd(user string, sorter string, filter string, postfix bool, m types.Msg, rsp types.Rsp) {
	rsp.Acknowledge()

	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	inv, res := dat.GetInv(user, user == m.Author.ID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	items := make([]string, len(inv))
	i := 0
	dat.Lock.RLock()
	for k := range inv {
		el, _ := dat.GetElement(k, true)
		items[i] = el.Name
		i++
	}
	dat.Lock.RUnlock()

	switch filter {
	case "madeby":
		count := 0
		outs := make([]string, len(items))
		for _, val := range items {
			creator := ""
			elem, res := dat.GetElement(val, true)
			if res.Exists {
				creator = elem.Creator
			}
			if creator == user {
				outs[count] = val
				count++
			}
		}
		outs = outs[:count]
		items = outs
	}

	if postfix {
		util.SortElemList(items, sorter, dat)
	} else {
		util.SortElemList(items, sorter, dat, true)
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
