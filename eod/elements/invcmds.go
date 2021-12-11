package elements

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/eodsort"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

func (b *Elements) ResetInvCmd(user string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}
	inv := db.GetInv(user)
	inv.Lock.Lock()
	inv.Elements = make(map[int]types.Empty)
	for _, el := range base.StarterElements {
		inv.Elements[el.ID] = types.Empty{}
	}
	inv.Lock.Unlock()

	err := db.SaveInv(inv, true)
	if rsp.Error(err) {
		return
	}
	rsp.Resp("Successfully reset <@" + user + ">'s inventory!")
}

func (b *Elements) DownloadInvCmd(user string, sorter string, filter string, postfix bool, m types.Msg, rsp types.Rsp) {
	rsp.Acknowledge()

	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}
	inv := db.GetInv(user)
	items := make([]int, len(inv.Elements))
	i := 0
	db.RLock()
	for k := range inv.Elements {
		el, _ := db.GetElement(k, true)
		items[i] = el.ID
		i++
	}

	switch filter {
	case "madeby":
		count := 0
		outs := make([]int, len(items))
		for _, val := range items {
			creator := ""
			elem, res := db.GetElement(val, true)
			if res.Exists {
				creator = elem.Creator
			}
			if creator == user {
				outs[count] = elem.ID
				count++
			}
		}
		outs = outs[:count]
		items = outs
	}

	if postfix {
		eodsort.SortElemList(items, sorter, db)
	} else {
		eodsort.SortElemList(items, sorter, db, true)
	}

	out := &strings.Builder{}
	db.Lock()
	for _, val := range items {
		elem, res := db.GetElement(val, true)
		if !res.Exists {
			continue
		}
		out.WriteString(elem.Name + "\n")
	}
	db.RUnlock()
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
