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
		return
	}

	for _, comb := range combs {
		if strings.Contains(m.Content, comb) {
			b.checkServer(msg, rsp)
			return
		}
	}
}
