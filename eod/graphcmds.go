package eod

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

const renderDifficulty = 30

func (b *EoD) graphCmd(elem string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	rsp.Acknowledge()

	dat.Lock.RLock()
	el, exists := dat.ElemCache[strings.ToLower(elem)]
	dat.Lock.RUnlock()
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", elem))
		return
	}

	graph, err := trees.NewGraph(dat)
	if rsp.Error(err) {
		return
	}

	msg, suc := graph.AddElem(elem, el.Difficulty >= renderDifficulty)
	if !suc {
		rsp.ErrorMessage(msg)
		return
	}

	var file *discordgo.File
	txt := "Sent graph in DMs!"
	if graph.NodeCount() < 200 {
		out, err := graph.RenderPNG()
		if rsp.Error(err) {
			return
		}
		file = &discordgo.File{
			Name:        "graph.png",
			ContentType: "image/png",
			Reader:      out,
		}
	} else {
		file = &discordgo.File{
			Name:        "graph.dot",
			ContentType: "text/plain",
			Reader:      strings.NewReader(graph.String()),
		}
		txt = "The graph was to big to render server-side! Check out https://graphviz.org/ to render it on your computer!"
	}

	rsp.Message(txt)

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}

	dat.Lock.RLock()
	name := dat.ElemCache[strings.ToLower(elem)].Name
	dat.Lock.RUnlock()
	b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Graph for **%s**:", name),
		Files:   []*discordgo.File{file},
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

	compl := 0
	dat.Lock.RLock()
	for elem := range cat.Elements {
		el, exists := dat.ElemCache[strings.ToLower(elem)]
		if !exists {
			rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", elem))
			return
		}
		compl += el.Difficulty
	}
	dat.Lock.RUnlock()

	graph, err := trees.NewGraph(dat)
	if rsp.Error(err) {
		return
	}

	for elem := range cat.Elements {
		msg, suc := graph.AddElem(elem, compl >= renderDifficulty)
		if !suc {
			rsp.ErrorMessage(msg)
			return
		}
	}

	var file *discordgo.File
	txt := "Sent graph in DMs!"
	if graph.NodeCount() < 200 {
		out, err := graph.RenderPNG()
		if rsp.Error(err) {
			return
		}
		file = &discordgo.File{
			Name:        "graph.png",
			ContentType: "image/png",
			Reader:      out,
		}
	} else {
		file = &discordgo.File{
			Name:        "graph.dot",
			ContentType: "text/plain",
			Reader:      strings.NewReader(graph.String()),
		}
		txt = "The graph was to big to render server-side! Check out https://graphviz.org/ to render it on your computer!"
	}

	rsp.Message(txt)

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}

	b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Graph for category **%s**:", cat.Name),
		Files:   []*discordgo.File{file},
	})
}
