package eod

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

func (b *EoD) notationCmd(elem string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	rsp.Acknowledge()
	tree := trees.NewNotationTree(dat)

	dat.Lock.RLock()
	msg, suc := tree.AddElem(elem)
	dat.Lock.RUnlock()
	if !suc {
		rsp.ErrorMessage(msg)
		return
	}

	txt := tree.String()

	if len(txt) <= 2000 {
		id := rsp.Message("Sent notation in DMs!")

		dat.SetMsgElem(id, elem)
		lock.Lock()
		b.dat[m.GuildID] = dat
		lock.Unlock()

		rsp.DM(txt)
		return
	}
	id := rsp.Message("The notation was too long! Sending it as a file in DMs!")

	dat.SetMsgElem(id, elem)
	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}
	buf := strings.NewReader(txt)
	el, _ := dat.GetElement(elem)
	b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Notation for **%s**:", el.Name),
		Files: []*discordgo.File{
			{
				Name:        "notation.txt",
				ContentType: "text/plain",
				Reader:      buf,
			},
		},
	})
}

func (b *EoD) catNotationCmd(catName string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	rsp.Acknowledge()
	tree := trees.NewNotationTree(dat)

	cat, res := dat.GetCategory(catName)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
	}

	dat.Lock.RLock()
	for elem := range cat.Elements {
		msg, suc := tree.AddElem(elem)
		if !suc {
			dat.Lock.RUnlock()
			rsp.ErrorMessage(msg)
			return
		}
	}
	dat.Lock.RUnlock()

	txt := tree.String()

	if len(txt) <= 2000 {
		rsp.Message("Sent notation in DMs!")

		rsp.DM(txt)
		return
	}
	rsp.Message("The notation was too long! Sending it as a file in DMs!")

	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}
	buf := strings.NewReader(txt)
	b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Notation for category **%s**:", cat.Name),
		Files: []*discordgo.File{
			{
				Name:        "notation.txt",
				ContentType: "text/plain",
				Reader:      buf,
			},
		},
	})
}
