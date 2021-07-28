package eod

import (
	"fmt"
	"strings"
)

func (b *EoD) markCmd(elem string, mark string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	dat.lock.RLock()
	el, exists := dat.elemCache[strings.ToLower(elem)]
	dat.lock.RUnlock()
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", elem))
		return
	}
	dat.lock.RLock()
	inv, exists := dat.invCache[m.Author.ID]
	dat.lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}
	_, exists = inv[strings.ToLower(el.Name)]
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Element **%s** is not in your inventory!", el.Name))
		return
	}
	if len(mark) >= 2400 {
		rsp.ErrorMessage("Creator marks must be under 2400 characters!")
		return
	}

	if el.Creator == m.Author.ID {
		b.mark(m.GuildID, elem, mark, "")
		rsp.Message(fmt.Sprintf("You have signed **%s**! 🖋️", el.Name))
		return
	}

	err := b.createPoll(poll{
		Channel: dat.votingChannel,
		Guild:   m.GuildID,
		Kind:    pollSign,
		Value1:  el.Name,
		Value2:  mark,
		Value3:  el.Comment,
		Value4:  m.Author.ID,
	})
	if rsp.Error(err) {
		return
	}
	rsp.Message(fmt.Sprintf("Suggested a note for **%s** 🖊️", el.Name))
}

func (b *EoD) imageCmd(elem string, image string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	dat.lock.RLock()
	el, exists := dat.elemCache[strings.ToLower(elem)]
	dat.lock.RUnlock()
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", elem))
		return
	}

	dat.lock.RLock()
	inv, exists := dat.invCache[m.Author.ID]
	dat.lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("You don't have an inventory!")
		return
	}
	_, exists = inv[strings.ToLower(el.Name)]
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Element **%s** is not in your inventory!", el.Name))
		return
	}

	if el.Creator == m.Author.ID {
		b.image(m.GuildID, elem, image, "")
		rsp.Resp(fmt.Sprintf("You added an image to **%s**! 📷", el.Name))
		return
	}

	err := b.createPoll(poll{
		Channel: dat.votingChannel,
		Guild:   m.GuildID,
		Kind:    pollImage,
		Value1:  el.Name,
		Value2:  image,
		Value3:  el.Image,
		Value4:  m.Author.ID,
	})
	if rsp.Error(err) {
		return
	}
	rsp.Message(fmt.Sprintf("Suggested an image for **%s** 📷", el.Name))
}
