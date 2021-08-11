package eod

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
	"github.com/goccy/go-graphviz"
)

var maxSizes = map[string]int{
	"Dot":   100,
	"Twopi": 1000,
}

var outputTypes = map[string]types.Empty{
	"PNG":  {},
	"SVG":  {},
	"Text": {},
	"DOT":  {},
}

func (b *EoD) graphCmd(elems map[string]types.Empty, dat types.ServerData, m types.Msg, layout string, outputType string, name string, distinctPrimary bool, rsp types.Rsp) {
	// Create graph
	graph, err := trees.NewGraph(dat)
	if rsp.Error(err) {
		return
	}

	for elem := range elems {
		msg, suc := graph.AddElem(elem, true)
		if !suc {
			rsp.ErrorMessage(msg)
			return
		}
	}

	// Automatically Select best layout and output type
	if outputType == "" {
		if layout != "" {
			outputType = "PNG"
		} else if graph.NodeCount() > maxSizes["Twopi"] {
			outputType = "DOT"
		} else if graph.NodeCount() > maxSizes["Dot"] {
			layout = "Twopi"
			outputType = "PNG"
		} else {
			layout = "Dot"
			outputType = "PNG"
		}
	} else if (outputType == "SVG" || outputType == "PNG") && layout == "" {
		if graph.NodeCount() > maxSizes["Dot"] {
			layout = "Twopi"
		} else {
			layout = "Dot"
		}
	}

	// Check input
	if !(outputType == "Text" || outputType == "DOT") {
		_, exists := maxSizes[layout]
		if !exists {
			rsp.ErrorMessage(fmt.Sprintf("Layout **%s** is invalid!", layout))
			return
		}

		if maxSizes[layout] > 0 && graph.NodeCount() > maxSizes[layout] {
			rsp.ErrorMessage(fmt.Sprintf("Graph is too big for layout **%s**!", layout))
		}
	}

	_, exists := outputTypes[outputType]
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Output type **%s** is invalid!", outputType))
		return
	}

	// Create Output
	var file *discordgo.File
	txt := "Sent graph in DMs!"

	switch outputType {
	case "PNG", "SVG":
		var out *bytes.Buffer
		var err error

		format := graphviz.PNG
		if outputType == "SVG" {
			format = graphviz.SVG
		}

		switch layout {
		case "Dot":
			out, err = graph.Render(false, graphviz.DOT, format)
		case "Twopi":
			out, err = graph.Render(false, graphviz.TWOPI, format)
		}

		if rsp.Error(err) {
			return
		}

		file = &discordgo.File{
			Name:        "graph.png",
			ContentType: "image/png",
			Reader:      out,
		}

		if outputType == "SVG" {
			file = &discordgo.File{
				Name:        "graph.svg",
				ContentType: "image/svg+xml",
				Reader:      out,
			}

		}
	case "Text", "DOT":
		txt = "The graph was not rendered server-side! Check out https://graphviz.org/ to render it on your computer!"
		name := "graph.dot"
		if outputType == "Text" {
			name = "graph.txt"
		}
		file = &discordgo.File{
			Name:        name,
			ContentType: "text/plain",
			Reader:      strings.NewReader(graph.String(distinctPrimary)),
		}
	}

	rsp.Message(txt)

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}

	b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Graph for **%s**:", name),
		Files:   []*discordgo.File{file},
	})
}

func (b *EoD) elemGraphCmd(elem string, layout string, outputType string, distinctPrimary bool, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	rsp.Acknowledge()

	dat.Lock.RLock()
	name := dat.ElemCache[strings.ToLower(elem)].Name
	dat.Lock.RUnlock()

	b.graphCmd(map[string]types.Empty{elem: {}}, dat, m, layout, outputType, name, distinctPrimary, rsp)
}

func (b *EoD) catGraphCmd(catName, layout, outputType string, distinctPrimary bool, m types.Msg, rsp types.Rsp) {
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

	b.graphCmd(cat.Elements, dat, m, layout, outputType, catName, distinctPrimary, rsp)
}
