package elements

import (
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/eodsort"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

func (b *Elements) ResetInvCmd(user string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}
	inv := db.GetInv(user)
	inv.Lock.Lock()
	inv.Elements = make(map[int]types.Empty)
	for _, el := range types.StarterElements {
		inv.Elements[el.ID] = types.Empty{}
	}
	inv.Lock.Unlock()

	err := db.SaveInv(inv, true)
	if rsp.Error(err) {
		return
	}
	rsp.Resp(db.Config.LangProperty("ResetUserInv", user))
}

func (b *Elements) DownloadInvCmd(user string, sorter string, filter string, postfix bool, m types.Msg, rsp types.Rsp) {
	rsp.Acknowledge()

	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}
	inv := db.GetInv(user)
	type invItem struct {
		name string
		id   int
	}
	items := make([]invItem, len(inv.Elements))
	i := 0
	db.RLock()
	for k := range inv.Elements {
		el, _ := db.GetElement(k, true)
		items[i] = invItem{el.Name, el.ID}
		i++
	}

	switch filter {
	case "madeby":
		count := 0
		outs := make([]invItem, len(items))
		for _, val := range items {
			creator := ""
			elem, res := db.GetElement(val.id, true)
			if res.Exists {
				creator = elem.Creator
			}
			if creator == user {
				outs[count] = invItem{elem.Name, elem.ID}
				count++
			}
		}
		outs = outs[:count]
		items = outs
	}

	eodsort.Sort(items, len(items), func(index int) int {
		return items[index].id
	}, func(index int) string {
		return items[index].name
	}, func(index int, val string) {
		items[index].name = val
	}, sorter, m.Author.ID, db, postfix)

	out := &strings.Builder{}
	for _, val := range items {
		out.WriteString(val.name + "\n")
	}
	db.RUnlock()
	buf := strings.NewReader(out.String())

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}

	usr, err := b.dg.User(user)
	if rsp.Error(err) {
		return
	}
	gld, err := b.dg.Guild(m.GuildID)
	if rsp.Error(err) {
		return
	}

	_, err = b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: db.Config.LangProperty("DownloadedInvUserServer", map[string]any{
			"Username": usr.Username,
			"Server":   gld.Name,
		}),
		Files: []*discordgo.File{
			b.base.PrepareFile(&discordgo.File{
				Name:        "inv.txt",
				ContentType: "text/plain",
				Reader:      buf,
			}, out.Len()),
		},
	})
	if rsp.Error(err) {
		return
	}
	rsp.Message(db.Config.LangProperty("SentInvToDMs", nil))
}
