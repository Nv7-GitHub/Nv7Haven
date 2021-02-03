package discord

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/bwmarrin/discordgo"
)

var varput = regexp.MustCompile(`var (.+)=([0-9.,]+)`)

func (b *Bot) math(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "var") {
		out := varput.FindAllStringSubmatch(m.Content, -1)
		if len(out) < 1 || len(out[0]) < 3 {
			s.ChannelMessageSend(m.ChannelID, "Invalid format. You need to use `var <name>=<value>`.")
			return
		}
		name := out[0][1]
		val, err := strconv.ParseFloat(out[0][2], 64)
		if b.handle(err, m) {
			return
		}

		_, exists := b.mathvars[m.GuildID]
		if !exists {
			b.mathvars[m.GuildID] = make(map[string]interface{})
		}
		b.mathvars[m.GuildID][name] = val
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfuly set variable %s to %f", name, val))
		return
	}

	if strings.HasPrefix(m.Content, "=") {
		gexp, err := govaluate.NewEvaluableExpression(m.Content[1:])
		if b.handle(err, m) {
			return
		}

		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()

		_, exists := b.mathvars[m.GuildID]
		if !exists {
			b.mathvars[m.GuildID] = make(map[string]interface{})
		}
		result, err := gexp.Evaluate(b.mathvars[m.GuildID])
		if b.handle(err, m) {
			return
		}
		b.mathvars[m.GuildID]["ans"] = result

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%v", result))
		return
	}
}
