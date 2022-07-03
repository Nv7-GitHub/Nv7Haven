package discord

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/bwmarrin/discordgo"
)

var varput = regexp.MustCompile(`var (.+)=(.+)`)
var functions = map[string]govaluate.ExpressionFunction{
	"sqrt": newFunc(math.Sqrt),
	"cos":  newFunc(math.Cos),
	"tan":  newFunc(math.Tan),
	"sin":  newFunc(math.Sin),
	"asin": newFunc(math.Asin),
	"acos": newFunc(math.Acos),
	"atan": newFunc(math.Atan),
	"exp": func(args ...any) (any, error) {
		if len(args) < 2 {
			return nil, errors.New("not enough arguments")
		}
		val1, ok := args[0].(float64)
		if !ok {
			intval, ok := args[0].(int)
			if !ok {
				return nil, errors.New("argument is not number")
			}
			val1 = float64(intval)
		}
		val2, ok := args[1].(float64)
		if !ok {
			intval, ok := args[1].(int)
			if !ok {
				return nil, errors.New("argument is not number")
			}
			val2 = float64(intval)
		}
		return math.Pow(val1, val2), nil
	},
}

func newFunc(fu func(float64) float64) func(...any) (any, error) {
	return func(args ...any) (any, error) {
		if len(args) < 1 {
			return nil, errors.New("not enough arguments")
		}
		val, ok := args[0].(float64)
		if !ok {
			intval, ok := args[0].(int)
			if !ok {
				return nil, errors.New("argument is not number")
			}
			val = float64(intval)
		}
		return fu(val), nil
	}
}

func (b *Bot) math(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if b.startsWith(m, "var") {
		out := varput.FindAllStringSubmatch(m.Content, -1)
		if len(out) < 1 || len(out[0]) < 3 {
			s.ChannelMessageSend(m.ChannelID, "Invalid format. You need to use `var <name>=<value>`.")
			return
		}
		name := out[0][1]

		gexp, err := govaluate.NewEvaluableExpressionWithFunctions(out[0][2], functions)
		if b.handle(err, m) {
			return
		}
		_, exists := b.mathvars[m.GuildID]
		if !exists {
			b.mathvars[m.GuildID] = make(map[string]any)
		}
		result, err := gexp.Evaluate(b.mathvars[m.GuildID])
		if b.handle(err, m) {
			return
		}
		b.mathvars[m.GuildID][name] = result
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfuly set variable %s to %v", name, result))
		return
	}

	if strings.HasPrefix(m.Content, "=") {
		gexp, err := govaluate.NewEvaluableExpressionWithFunctions(m.Content[1:], functions)
		if b.handle(err, m) {
			return
		}

		_, exists := b.mathvars[m.GuildID]
		if !exists {
			b.mathvars[m.GuildID] = make(map[string]any)
		}
		result, err := gexp.Evaluate(b.mathvars[m.GuildID])
		if b.handle(err, m) {
			return
		}
		_, ok := result.(float64)
		if ok {
			b.mathvars[m.GuildID]["ans"] = result
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%v", result))
		return
	}
}
