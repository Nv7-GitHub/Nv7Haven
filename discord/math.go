package discord

import (
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) math(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "var") {
		var name string
		var val float64
		_, err := fmt.Sscanf(m.Content, "var %s=%f", &name, &val)
		if b.handle(err, m) {
			return
		}

		b.mathvars[m.GuildID][name] = val
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfuly set variable %s to %f", name, val))
		return
	}

	if strings.HasPrefix(m.Content, "=") {
		var expression string
		_, err := fmt.Sscanf(m.Content, "=%s", &expression)
		if b.handle(err, m) {
			return
		}

		gexp, err := govaluate.NewEvaluableExpression(expression)
		if b.handle(err, m) {
			return
		}

		result, err := gexp.Evaluate(b.mathvars[m.GuildID])
		if b.handle(err, m) {
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%v", result))
		return
	}
}
