package treecmds

import (
	"bytes"
	"image/png"

	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

var (
	WCMinWidth  float64 = 1
	WCMinHeight float64 = 1
	WCMaxWidth  float64 = 4096
	WCMaxHeight float64 = 4096
)

func (b *TreeCmds) WordCloudCmd(name string, elems map[int]types.Empty, calcTree bool, width, height int, m types.Msg, rsp types.Rsp) {
	rsp.Acknowledge()

	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	tree := trees.NewWordTree(db)
	tree.CalcTree = calcTree
	for elem := range elems {
		suc, msg := tree.AddElem(elem)
		if !suc {
			rsp.ErrorMessage(msg)
			return
		}
	}

	out := bytes.NewBuffer(nil)
	im := tree.Render(width, height)
	err := png.Encode(out, im)
	if rsp.Error(err) {
		return
	}

	rsp.Message(db.Config.LangProperty("SentWordCloud", nil))
	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}

	b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: db.Config.LangProperty("WordCloudElem", name),
		Files: []*discordgo.File{
			{
				Name:        "wordcloud.png",
				ContentType: "image/png",
				Reader:      out,
			},
		},
	})
}

func (b *TreeCmds) ElemWordCloudCmd(elem string, calcTree bool, width, height int, m types.Msg, rsp types.Rsp) {
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
	b.WordCloudCmd(el.Name, map[int]types.Empty{el.ID: {}}, calcTree, width, height, m, rsp)
}

func (b *TreeCmds) CatWordCloudCmd(catName string, calcTree bool, width, height int, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	rsp.Acknowledge()

	var els map[int]types.Empty
	catv, res := db.GetCat(catName)
	if !res.Exists {
		vcat, res := db.GetVCat(catName)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
		catName = vcat.Name
		els, res = b.base.CalcVCat(vcat, db)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
	} else {
		catv.Lock.RLock()
		els = make(map[int]types.Empty, len(catv.Elements))
		for el := range catv.Elements {
			els[el] = types.Empty{}
		}
		catv.Lock.RUnlock()
		catName = catv.Name
	}
	b.WordCloudCmd(catName, els, calcTree, width, height, m, rsp)
}

func (b *TreeCmds) InvWordCloudCmd(user string, calcTree bool, width, height int, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	rsp.Acknowledge()

	inv := db.GetInv(user)
	inv.Lock.RLock()
	defer inv.Lock.RUnlock()

	usr, err := b.dg.User(user)
	if rsp.Error(err) {
		return
	}
	b.WordCloudCmd(usr.Username, inv.Elements, calcTree, width, height, m, rsp)
}
