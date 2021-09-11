package eod

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/types"

	"github.com/bwmarrin/discordgo"
)

func (b *EoD) giveCmd(elem string, giveTree bool, user string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	inv, res := dat.GetInv(user, true)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	el, res := dat.GetElement(elem)
	if !res.Exists {
		rsp.Resp(res.Message)
		return
	}

	msg, suc := giveElem(dat, giveTree, elem, &inv)
	if !suc {
		rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", msg))
		return
	}

	dat.SetInv(user, inv)

	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()
	b.saveInv(m.GuildID, user, true, true)

	rsp.Resp("Successfully gave element **" + el.Name + "**!")
}

func (b *EoD) giveCatCmd(catName string, giveTree bool, user string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	inv, res := dat.GetInv(user, true)
	if !exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	cat, res := dat.GetCategory(catName)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	for elem := range cat.Elements {
		_, res := dat.GetElement(elem)
		if !res.Exists {
			rsp.Resp(fmt.Sprintf("Element **%s** doesn't exist!", elem))
			return
		}

		msg, suc := giveElem(dat, giveTree, elem, &inv)
		if !suc {
			rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", msg))
			return
		}
	}

	dat.SetInv(user, inv)

	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()
	b.saveInv(m.GuildID, user, true, true)

	rsp.Resp("Successfully gave all elements in category **" + cat.Name + "**!")
}

func giveElem(dat types.ServerData, giveTree bool, elem string, out *types.Container) (string, bool) {
	el, res := dat.GetElement(elem)
	if !res.Exists {
		return elem, false
	}
	if giveTree {
		for _, parent := range el.Parents {
			if len(strings.TrimSpace(parent)) == 0 {
				continue
			}
			_, exists := (*out)[strings.ToLower(parent)]
			if !exists {
				msg, suc := giveElem(dat, giveTree, parent, out)
				if !suc {
					return msg, false
				}
			}
		}
	}
	(*out)[strings.ToLower(el.Name)] = types.Empty{}
	return "", true
}

func (b *EoD) calcTreeCmd(elem string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
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
		lock.Lock()
		b.dat[m.GuildID] = dat
		lock.Unlock()

		rsp.DM(txt)
		return
	}
	id := rsp.Message("The path was too long! Sending it as a file in DMs!")

	dat.SetMsgElem(id, elem)
	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()

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

func (b *EoD) calcTreeCatCmd(catName string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
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

func (b *EoD) notationCmd(elem string, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	rsp.Acknowledge()
	tree := trees.NewNotationTree(dat)

	dat.Lock.RLock()
	txt, suc := tree.GetNotation(elem)
	dat.Lock.RUnlock()
	if !suc {
		rsp.ErrorMessage(txt)
		return
	}

	if len(txt) <= 2000 {
		id := rsp.Message("Sent notation in DMs!")

		dat.SetMsgElem(id, elem)
		lock.Lock()
		b.dat[m.GuildID] = dat
		lock.Unlock()

		rsp.DM(txt)
		return
	}
	id := rsp.Message("The path was too long! Sending it as a file in DMs!")

	dat.SetMsgElem(id, elem)
	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}
	buf := strings.NewReader(txt)
	el, _ := dat.GetElement(elem)
	b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Notation for **%s**:", el.Name),
		Files: []*discordgo.File{
			{
				Name:        "notation.txt",
				ContentType: "text/plain",
				Reader:      buf,
			},
		},
	})
}
