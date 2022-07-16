package treecmds

import (
	"strings"
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

func (b *TreeCmds) NotationCmd(elem string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	rsp.Acknowledge()
	tree := trees.NewNotationTree(db)

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

	db.RLock()
	msg, suc := tree.AddElem(el.ID)
	db.RUnlock()
	if !suc {
		rsp.ErrorMessage(msg)
		return
	}

	txt := tree.String()
	data, res := b.GetData(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	if len(txt) <= 2000 {
		id := rsp.Message(db.Config.LangProperty("SentNotationToDMs", nil))
		data.SetMsgElem(id, el.ID)
		rsp.DM(txt)
		return
	}

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}
	buf := strings.NewReader(txt)
	_, err = b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: db.Config.LangProperty("NameNotationElem", el.Name),
		Files: []*discordgo.File{
			b.base.PrepareFile(&discordgo.File{
				Name:        "notation.txt",
				ContentType: "text/plain",
				Reader:      buf,
			}, len(txt)),
		},
	})
	if rsp.Error(err) {
		return
	}
	id := rsp.Message(db.Config.LangProperty("NotationTooLong", nil))

	data.SetMsgElem(id, el.ID)
}

func (b *TreeCmds) CatNotationCmd(catName string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	rsp.Acknowledge()
	tree := trees.NewNotationTree(db)

	var els map[int]types.Empty
	var lock *sync.RWMutex
	catv, res := db.GetCat(catName)
	if !res.Exists {
		vcat, res := db.GetVCat(catName)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
		catName = vcat.Name
		els, res = b.base.CalcVCat(vcat, db, true)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
	} else {
		lock = catv.Lock
		els = catv.Elements
		catName = catv.Name
	}

	inv := db.GetInv(m.Author.ID)
	for k := range els {
		if !inv.Contains(k) {
			rsp.ErrorMessage(db.Config.LangProperty("MustHaveCatForPath", catName))
			return
		}
	}

	db.RLock()
	if lock != nil {
		lock.RLock()
	}
	for elem := range els {
		msg, suc := tree.AddElem(elem)
		if !suc {
			db.RUnlock()
			rsp.ErrorMessage(msg)
			return
		}
	}
	if lock != nil {
		lock.RUnlock()
	}
	db.RUnlock()

	txt := tree.String()

	if len(txt) <= 2000 {
		rsp.Message(db.Config.LangProperty("SentNotationToDMs", nil))

		rsp.DM(txt)
		return
	}

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}
	buf := strings.NewReader(txt)
	_, err = b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: db.Config.LangProperty("NameNotationCat", catName),
		Files: []*discordgo.File{
			b.base.PrepareFile(&discordgo.File{
				Name:        "notation.txt",
				ContentType: "text/plain",
				Reader:      buf,
			}, len(txt)),
		},
	})
	if rsp.Error(err) {
		return
	}
	rsp.Message(db.Config.LangProperty("NotationTooLong", nil))
}
