package treecmds

import (
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/types"

	"github.com/bwmarrin/discordgo"
)

func (b *TreeCmds) CalcTreeCmd(elem string, m types.Msg, rsp types.Rsp) {
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

	// Check if can
	inv := db.GetInv(m.Author.ID)
	if !inv.Contains(el.ID) {
		rsp.ErrorMessage(db.Config.LangProperty("MustHaveElemForPath", el.Name))
		return
	}

	txt, suc, msg := trees.CalcTree(db, el.ID)
	if !suc {
		rsp.ErrorMessage(msg)
		return
	}
	data, res := b.GetData(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	if len(txt) <= 2000 {
		id := rsp.Message(db.Config.LangProperty("SentPathToDMs", nil))

		data.SetMsgElem(id, el.ID)

		rsp.DM(txt)
		return
	}
	id := rsp.Message(db.Config.LangProperty("PathTooLong", nil))

	data.SetMsgElem(id, el.ID)

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}
	buf := strings.NewReader(txt)

	b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: db.Config.LangProperty("NamePathElem", el.Name),
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
		els = make(map[int]types.Empty, len(catv.Elements))
		catv.Lock.RLock()
		for k := range catv.Elements {
			els[k] = types.Empty{}
		}
		catv.Lock.RUnlock()
		catName = catv.Name
	}

	// Check if can
	inv := db.GetInv(m.Author.ID)
	for k := range els {
		if !inv.Contains(k) {
			rsp.ErrorMessage(db.Config.LangProperty("MustHaveCatForPath", catName))
			return
		}
	}

	txt, suc, msg := trees.CalcTreeCat(db, els)
	if !suc {
		rsp.ErrorMessage(msg)
		return
	}
	if len(txt) <= 2000 {
		rsp.Message(db.Config.LangProperty("SentPathToDMs", nil))
		rsp.DM(txt)
		return
	}
	rsp.Message(db.Config.LangProperty("PathTooLong", nil))

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}
	buf := strings.NewReader(txt)
	b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: db.Config.LangProperty("NamePathCat", catName),
		Files: []*discordgo.File{
			{
				Name:        "path.txt",
				ContentType: "text/plain",
				Reader:      buf,
			},
		},
	})
}
