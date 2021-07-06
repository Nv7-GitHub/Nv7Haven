package eod

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	deadlock "github.com/sasha-s/go-deadlock"
)

func (b *EoD) giveCmd(elem string, giveTree bool, user string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	dat.lock.RLock()
	inv, exists := dat.invCache[user]
	dat.lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}

	dat.lock.RLock()
	el, exists := dat.elemCache[strings.ToLower(elem)]
	dat.lock.RUnlock()
	if !exists {
		rsp.Resp(fmt.Sprintf("Element **%s** doesn't exist!", elem))
		return
	}

	msg, suc := giveElem(dat.elemCache, giveTree, elem, &inv, dat.lock)
	if !suc {
		rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", msg))
		return
	}

	dat.lock.Lock()
	dat.invCache[user] = inv
	dat.lock.Unlock()

	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()
	b.saveInv(m.GuildID, user, true, true)

	rsp.Resp("Successfully gave element **" + el.Name + "**!")
}

func (b *EoD) giveCatCmd(catName string, giveTree bool, user string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	dat.lock.RLock()
	inv, exists := dat.invCache[user]
	dat.lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}
	cat, exists := dat.catCache[strings.ToLower(catName)]
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Category **%s** doesn't exist!", catName))
		return
	}

	for elem := range cat.Elements {
		dat.lock.RLock()
		_, exists := dat.elemCache[strings.ToLower(elem)]
		dat.lock.RUnlock()
		if !exists {
			rsp.Resp(fmt.Sprintf("Element **%s** doesn't exist!", elem))
			return
		}

		msg, suc := giveElem(dat.elemCache, giveTree, elem, &inv, dat.lock)
		if !suc {
			rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", msg))
			return
		}
	}

	dat.lock.Lock()
	dat.invCache[user] = inv
	dat.lock.Unlock()

	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()
	b.saveInv(m.GuildID, user, true, true)

	rsp.Resp("Successfully gave all elements in category **" + cat.Name + "**!")
}

func giveElem(elemCache map[string]element, giveTree bool, elem string, out *map[string]empty, lock *deadlock.RWMutex) (string, bool) {
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
	(*out)[strings.ToLower(el.Name)] = empty{}
	return "", true
}

func (b *EoD) calcTreeCmd(elem string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	rsp.Acknowledge()
	txt, suc, msg := calcTree(dat.elemCache, elem, dat.lock)
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
	dat.lock.RLock()
	name := dat.elemCache[strings.ToLower(elem)].Name
	dat.lock.RUnlock()
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

func (b *EoD) calcTreeCatCmd(catName string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	rsp.Acknowledge()
	cat, exists := dat.catCache[strings.ToLower(catName)]
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Category **%s** doesn't exist!", catName))
		return
	}

	txt, suc, msg := calcTreeCat(dat.elemCache, cat.Elements, dat.lock)
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
