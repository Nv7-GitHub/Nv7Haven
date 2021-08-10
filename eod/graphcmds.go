package eod

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

func (b *EoD) graphCmd(elem string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	rsp.Acknowledge()

	graph, err := trees.NewGraph(dat)
	if rsp.Error(err) {
		return
	}
	msg, suc := graph.AddElem(elem)
	if !suc {
		rsp.ErrorMessage(msg)
		return
	}
	out, err := graph.RenderPNG()
	if rsp.Error(err) {
		return
	}

	rsp.Message("Sent graph in DMs!")

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}

	dat.Lock.RLock()
	name := dat.ElemCache[strings.ToLower(elem)].Name
	dat.Lock.RUnlock()
	b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Graph for **%s**:", name),
		Files: []*discordgo.File{
			{
				Name:        "graph.png",
				ContentType: "image/png",
				Reader:      out,
			},
		},
	})
}

func (b *EoD) catGraphCmd(catName string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	rsp.Acknowledge()
	cat, exists := dat.CatCache[strings.ToLower(catName)]
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Category **%s** doesn't exist!", catName))
		return
	}

	graph, err := trees.NewGraph(dat)
	if rsp.Error(err) {
		return
	}

	for elem := range cat.Elements {
		msg, suc := graph.AddElem(elem)
		if !suc {
			rsp.ErrorMessage(msg)
			return
		}
	}

	out, err := graph.RenderPNG()
	if rsp.Error(err) {
		return
	}

	rsp.Message("Sent graph in DMs!")

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}

	b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Graph for category **%s**:", cat.Name),
		Files: []*discordgo.File{
			{
				Name:        "graph.png",
				ContentType: "image/png",
				Reader:      out,
			},
		},
	})
}
