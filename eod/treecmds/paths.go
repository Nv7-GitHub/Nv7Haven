package treecmds

import (
	"fmt"
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
	txt, suc, msg := trees.CalcTree(db, el.ID)
	if !suc {
		rsp.ErrorMessage(fmt.Sprintf(db.Config.LangProperty("DoesntExist"), msg))
		return
	}
	data, res := b.GetData(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	if len(txt) <= 2000 {
		id := rsp.Message(db.Config.LangProperty("SentPathToDMs"))

		data.SetMsgElem(id, el.ID)

		rsp.DM(txt)
		return
	}
	id := rsp.Message(db.Config.LangProperty("PathTooLong"))

	data.SetMsgElem(id, el.ID)

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}
	buf := strings.NewReader(txt)

	b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf(db.Config.LangProperty("NamePathElem"), el.Name),
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

	cat, res := db.GetCat(catName)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	txt, suc, msg := trees.CalcTreeCat(db, cat.Elements)
	if !suc {
		rsp.ErrorMessage(fmt.Sprintf(db.Config.LangProperty("DoesntExist"), msg))
		return
	}
	if len(txt) <= 2000 {
		rsp.Message(db.Config.LangProperty("SentPathToDMs"))
		rsp.DM(txt)
		return
	}
	rsp.Message(db.Config.LangProperty("PathTooLong"))

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}
	buf := strings.NewReader(txt)
	b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf(db.Config.LangProperty("NamePathCat"), cat.Name),
		Files: []*discordgo.File{
			{
				Name:        "path.txt",
				ContentType: "text/plain",
				Reader:      buf,
			},
		},
	})
}
