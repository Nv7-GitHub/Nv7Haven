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

	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if strings.HasPrefix(m.Content, "*2") {
		b.checkServer(msg, rsp)
		lock.RLock()
		dat, exists := b.dat[msg.GuildID]
		lock.RUnlock()
		if !exists {
			return
		}
		if dat.combCache == nil {
			dat.combCache = make(map[string]comb)
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
			b.combine(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), msg, rsp)
			return
		}
	}

	if strings.HasPrefix(m.Content, "?") {
		b.infoCmd(strings.TrimSpace(m.Content[1:]), msg, rsp)
	}
}
