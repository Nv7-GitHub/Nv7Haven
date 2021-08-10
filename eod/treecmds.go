package eod

import (
	"fmt"
	"strings"
	"sync"

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
	dat.Lock.RLock()
	inv, exists := dat.InvCache[user]
	dat.Lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}

	dat.Lock.RLock()
	el, exists := dat.ElemCache[strings.ToLower(elem)]
	dat.Lock.RUnlock()
	if !exists {
		rsp.Resp(fmt.Sprintf("Element **%s** doesn't exist!", elem))
		return
	}

	msg, suc := giveElem(dat.ElemCache, giveTree, elem, &inv, dat.Lock)
	if !suc {
		rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", msg))
		return
	}

	dat.Lock.Lock()
	dat.InvCache[user] = inv
	dat.Lock.Unlock()

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
	dat.Lock.RLock()
	inv, exists := dat.InvCache[user]
	dat.Lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}
	cat, exists := dat.CatCache[strings.ToLower(catName)]
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Category **%s** doesn't exist!", catName))
		return
	}

	for elem := range cat.Elements {
		dat.Lock.RLock()
		_, exists := dat.ElemCache[strings.ToLower(elem)]
		dat.Lock.RUnlock()
		if !exists {
			rsp.Resp(fmt.Sprintf("Element **%s** doesn't exist!", elem))
			return
		}

		msg, suc := giveElem(dat.ElemCache, giveTree, elem, &inv, dat.Lock)
		if !suc {
			rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", msg))
			return
		}
	}

	dat.Lock.Lock()
	dat.InvCache[user] = inv
	dat.Lock.Unlock()

	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()
	b.saveInv(m.GuildID, user, true, true)

	rsp.Resp("Successfully gave all elements in category **" + cat.Name + "**!")
}

func giveElem(elemCache map[string]types.Element, giveTree bool, elem string, out *map[string]types.Empty, lock *sync.RWMutex) (string, bool) {
	lock.RLock()
	el, exists := elemCache[strings.ToLower(elem)]
	lock.RUnlock()
	if !exists {
		return elem, false
	}
	if giveTree {
		for _, parent := range el.Parents {
			if len(strings.TrimSpace(parent)) == 0 {
				continue
			}
			_, exists := (*out)[strings.ToLower(parent)]
			if !exists {
				msg, suc := giveElem(elemCache, giveTree, parent, out, lock)
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
	txt, suc, msg := calcTree(dat.ElemCache, elem, dat.Lock)
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
	dat.Lock.RLock()
	name := dat.ElemCache[strings.ToLower(elem)].Name
	dat.Lock.RUnlock()
	b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Path for **%s**:", name),
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
	cat, exists := dat.CatCache[strings.ToLower(catName)]
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Category **%s** doesn't exist!", catName))
		return
	}

	txt, suc, msg := calcTreeCat(dat.ElemCache, cat.Elements, dat.Lock)
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
