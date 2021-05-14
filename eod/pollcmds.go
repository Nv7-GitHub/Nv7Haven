package eod

import (
	"fmt"
	"strings"
)

func (b *EoD) suggestCmd(suggestion string, autocapitalize bool, m msg, rsp rsp) {
	suggestion = strings.Replace(suggestion, "+", "", -1)
	suggestion = strings.Replace(suggestion, "\\", "", -1)
	suggestion = strings.Replace(suggestion, "|", "", -1)
	suggestion = strings.TrimSpace(suggestion)
	if len(suggestion) > 1 && suggestion[0] == '#' {
		suggestion = suggestion[1:]
	}
	if len(suggestion) == 0 {
		rsp.Resp("You need to suggest something!")
		return
	}
	suggestIllegalTest = strings.Replace(suggestion, "<@", "", -1)
	if len(suggestIllegalTest) != len(suggestion) {
		rsp.Resp("Element names cannot contain pings.")
		return
	}

	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	suggestIllegalTest = strings.Replace(suggestIllegalTest, "@everyone", "", -1)
	suggestIllegalTest = strings.Replace(suggestIllegalTest, "@here", "", -1)
	if autocapitalize || len(suggestIllegalTest) != len(suggestion) {
		suggestion = strings.Title(suggestion)
	}
	if dat.combCache == nil {
		dat.combCache = make(map[string]comb)
	}
	comb, exists := dat.combCache[m.Author.ID]
	if !exists {
		rsp.ErrorMessage("You haven't combined anything!")
		return
	}

	data := elems2txt(comb.elems)
	row := b.db.QueryRow("SELECT COUNT(1) FROM eod_combos WHERE guild=? AND elems=?", m.GuildID, data)
	var count int
	err := row.Scan(&count)
	if rsp.Error(err) {
		return
	}
	if count != 0 {
		rsp.ErrorMessage("That combo already has a result!")
		return
	}

	err = b.createPoll(poll{
		Channel:   dat.votingChannel,
		Guild:     m.GuildID,
		Kind:      pollCombo,
		Value3:    suggestion,
		Value4:    m.Author.ID,
		Data:      map[string]interface{}{"elems": comb.elems},
		Upvotes:   0,
		Downvotes: 0,
	})
	if rsp.Error(err) {
		return
	}
	txt := "Suggested "
	for _, val := range comb.elems {
		txt += dat.elemCache[strings.ToLower(val)].Name + " + "
	}
	txt = txt[:len(txt)-3]
	if len(comb.elems) == 1 {
		txt += " + " + dat.elemCache[strings.ToLower(comb.elems[0])].Name
	}
	txt += " = " + suggestion + " ✨"
	rsp.Resp(txt)
}

func (b *EoD) markCmd(elem string, mark string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	el, exists := dat.elemCache[strings.ToLower(elem)]
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Element %s doesn't exist!", elem))
		return
	}

	if el.Creator == m.Author.ID {
		b.mark(m.GuildID, elem, mark, "")
		rsp.Resp(fmt.Sprintf("You have signed **%s**! 🖋️", el.Name))
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
	rsp.Resp(fmt.Sprintf("Suggested a note for **%s** 🖊️", el.Name))
}

func (b *EoD) imageCmd(elem string, image string, m msg, rsp rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	el, exists := dat.elemCache[strings.ToLower(elem)]
	if !exists {
		rsp.ErrorMessage(fmt.Sprintf("Element %s doesn't exist!", elem))
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
	rsp.Resp(fmt.Sprintf("Suggested an image for **%s** 📷", el.Name))
}
