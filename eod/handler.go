package eod

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

var combs = []string{
	"+",
	",",
}

func (b *EoD) cmdHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := b.newMsgNormal(m)
	rsp := b.newRespNormal(m)

	if strings.HasPrefix(m.Content, "*2") {
		b.checkServer(msg, rsp)
		lock.RLock()
		dat, exists := b.dat[msg.GuildID]
		lock.RUnlock()
		if !exists {
			return
		}
		comb, exists := dat.combCache[msg.Author.ID]
		if !exists {
			return
		}
		b.combine(comb.elem1, comb.elem2, msg, rsp)
		return
	}

	for _, comb := range combs {
		if strings.Contains(m.Content, comb) {
			b.checkServer(msg, rsp)
			parts := strings.Split(m.Content, comb)
			if len(parts) < 2 {
				return
			}
			b.combine(parts[0], parts[1], msg, rsp)
			return
		}
	}
}
