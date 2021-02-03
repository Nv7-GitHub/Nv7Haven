package discord

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/bwmarrin/discordgo"
)

var varput = regexp.MustCompile(`var (.+)=([0-9.,]+)`)

func newFunc(fu func(float64) float64) func(...interface{}) (interface{}, error) {
	return func(args ...interface{}) (interface{}, error) {
		if len(args) < 1 {
			return nil, errors.New("Not enough arguments")
		}
		val, ok := args[0].(float64)
		if !ok {
			intval, ok := args[0].(int)
			if !ok {
				return nil, errors.New("Argument is not number")
			}
			val = float64(intval)
		}
		return fu(val), nil
	}
}

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
		functions := map[string]govaluate.ExpressionFunction{
			"sqrt": newFunc(math.Sqrt),
			"cos":  newFunc(math.Cos),
			"tan":  newFunc(math.Tan),
			"sin":  newFunc(math.Sin),
			"asin": newFunc(math.Asin),
			"acos": newFunc(math.Acos),
			"atan": newFunc(math.Atan),
			"exp": func(args ...interface{}) (interface{}, error) {
				if len(args) < 2 {
					return nil, errors.New("Not enough arguments")
				}
				val1, ok := args[0].(float64)
				if !ok {
					intval, ok := args[0].(int)
					if !ok {
						return nil, errors.New("Argument is not number")
					}
					val1 = float64(intval)
				}
				val2, ok := args[1].(float64)
				if !ok {
					intval, ok := args[1].(int)
					if !ok {
						return nil, errors.New("Argument is not number")
					}
					val2 = float64(intval)
				}
				return math.Pow(val1, val2), nil
			},
		}

		gexp, err := govaluate.NewEvaluableExpressionWithFunctions(m.Content[1:], functions)
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
		_, ok := result.(float64)
		if ok {
			b.mathvars[m.GuildID]["ans"] = result
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%v", result))
		return
	}
}
