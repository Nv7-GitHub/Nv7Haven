package treecmds

import (
	"bytes"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
	"github.com/goccy/go-graphviz"
)

var maxSizes = map[string]int{
	"Dot":   115, // This is 7 * 3 multiplied by a number nice
	"Twopi": 504, // 7 * 3^2 * 2^3 also very cool
}

var outputTypes = map[string]types.Empty{
	"PNG":  {},
	"SVG":  {},
	"Text": {},
	"DOT":  {},
}

func (b *TreeCmds) graphCmd(elems map[int]types.Empty, db *eodb.DB, m types.Msg, layout string, outputType string, name string, distinctPrimary bool, rsp types.Rsp) {
	// Create graph
	graph, err := trees.NewGraph(db)
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
			rsp.ErrorMessage(db.Config.LangProperty("GraphLayoutInvalid", layout))
			return
		}

		if maxSizes[layout] > 0 && graph.NodeCount() > maxSizes[layout] {
			rsp.ErrorMessage(db.Config.LangProperty("GraphTooBigForLayout", layout))
			return
		}
	}

	_, exists := outputTypes[outputType]
	if !exists {
		rsp.ErrorMessage(db.Config.LangProperty("GraphOutputInvalid", outputType))
		return
	}

	// Create Output
	var file *discordgo.File
	txt := db.Config.LangProperty("SentGraphToDMs", nil)

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
		txt = db.Config.LangProperty("GraphNotRendered", nil)
		name := "graph.dot"
		if outputType == "Text" {
			name = "graph.txt"
		}
		splines := "ortho"
		if layout == "Twopi" {
			splines = "false"
		}
		file = &discordgo.File{
			Name:        name,
			ContentType: "text/plain",
			Reader:      strings.NewReader(graph.String(distinctPrimary, splines)),
		}
	}

	id := rsp.Message(txt)
	if len(elems) == 1 {
		var elem int
		for k := range elems {
			elem = k
			break
		}
		dat, res := b.GetData(m.GuildID)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}

		dat.SetMsgElem(id, elem)
	}

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}

	b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: db.Config.LangProperty("NameGraphElem", name),
		Files:   []*discordgo.File{file},
	})
}

func (b *TreeCmds) ElemGraphCmd(elem string, layout string, outputType string, distinctPrimary bool, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	rsp.Acknowledge()

	el, res := db.GetElementByName(elem)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	inv := db.GetInv(m.Author.ID)
	if !inv.Contains(el.ID) {
		rsp.ErrorMessage(db.Config.LangProperty("MustHaveElemForPath", el.Name))
		return
	}

	b.graphCmd(map[int]types.Empty{el.ID: {}}, db, m, layout, outputType, elem, distinctPrimary, rsp)
}

func (b *TreeCmds) CatGraphCmd(catName, layout, outputType string, distinctPrimary bool, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	rsp.Acknowledge()

	var els map[int]types.Empty
	cat, res := db.GetCat(catName)
	if !res.Exists {
		vcat, res := db.GetVCat(catName)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
		els, res = b.base.CalcVCat(vcat, db, true)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
		catName = vcat.Name
	} else {
		cat.Lock.RLock()
		els = make(map[int]types.Empty, len(cat.Elements))
		for el := range cat.Elements {
			els[el] = types.Empty{}
		}
		cat.Lock.RUnlock()
		catName = cat.Name
	}

	inv := db.GetInv(m.Author.ID)
	for k := range els {
		if !inv.Contains(k) {
			rsp.ErrorMessage(db.Config.LangProperty("MustHaveCatForPath", catName))
			return
		}
	}

	b.graphCmd(els, db, m, layout, outputType, catName, distinctPrimary, rsp)
}
