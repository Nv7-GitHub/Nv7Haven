package treecmds

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/types"

	"github.com/bwmarrin/discordgo"
)

func (b *TreeCmds) CalcTreeCmd(elem string, m types.Msg, rsp types.Rsp) {
	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		return
	}
	rsp.Acknowledge()
	txt, suc, msg := trees.CalcTree(dat, elem)
	if !suc {
		rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", msg))
		return
	}
	if len(txt) <= 2000 {
		id := rsp.Message("Sent path in DMs!")

		dat.SetMsgElem(id, elem)
		b.lock.Lock()
		b.dat[m.GuildID] = dat
		b.lock.Unlock()

		rsp.DM(txt)
		return
	}
	id := rsp.Message("The path was too long! Sending it as a file in DMs!")

	dat.SetMsgElem(id, elem)
	b.lock.Lock()
	b.dat[m.GuildID] = dat
	b.lock.Unlock()

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}
	buf := strings.NewReader(txt)
	el, _ := dat.GetElement(elem)
	b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Path for **%s**:", el.Name),
		Files: []*discordgo.File{
			{
				Name:        "path.txt",
				ContentType: "text/plain",
				Reader:      buf,
			},
		},
	})
}

func (b *TreeCmds) CalcTreeCatCmd(catName string, m types.Msg, rsp types.Rsp) {
	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		return
	}
	rsp.Acknowledge()

	cat, res := dat.GetCategory(catName)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	txt, suc, msg := trees.CalcTreeCat(dat, cat.Elements)
	if !suc {
		rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", msg))
		return
	}
	if len(txt) <= 2000 {
		rsp.Message("Sent path in DMs!")
		rsp.DM(txt)
		return
	}
	rsp.Message("The path was too long! Sending it as a file in DMs!")

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}
	buf := strings.NewReader(txt)
	b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Path for category **%s**:", cat.Name),
		Files: []*discordgo.File{
			{
				Name:        "path.txt",
				ContentType: "text/plain",
				Reader:      buf,
			},
		},
	})
}
