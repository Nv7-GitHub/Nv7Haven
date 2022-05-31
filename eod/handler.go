package eod

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/dgutil"
	"github.com/Nv7-Github/Nv7Haven/eod/logs"
	"github.com/bwmarrin/discordgo"
)

const guild = "" // 819077688371314718 for testing

func (b *EoD) initHandlers() {
	// Debugging
	discordgo.Logger = func(msgL, caller int, format string, a ...any) {
		// This code is a slightly modified version of https://github.com/bwmarrin/discordgo/blob/577e7dd4f6ccf1beb10acdb1871300c7638b84c4/logging.go#L46
		pc, file, line, _ := runtime.Caller(caller)

		files := strings.Split(file, "/")
		file = files[len(files)-1]

		name := runtime.FuncForPC(pc).Name()
		fns := strings.Split(name, ".")
		name = fns[len(fns)-1]

		msg := fmt.Sprintf(format, a...)

		log.SetOutput(logs.DiscordLogs)
		log.Printf("[DG%d] %s:%d:%s() %s\n", msgL, file, line, name, msg)
	}

	dgutil.UpdateBotCommands(b.dg, clientID, guild, commands)

	b.dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Member == nil {
			return
		}

		switch i.Type {
		// Command
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}

		// Button
		case discordgo.InteractionMessageComponent:
			data, res := b.GetData(i.GuildID)
			if !res.Exists {
				return
			}

			// Check if page switch handler or component handler
			_, exists := data.PageSwitchers[i.Message.ID]
			if exists {
				b.base.PageSwitchHandler(s, i)
				return
			}

			compMsg, exists := data.ComponentMsgs[i.Message.ID]
			if exists {
				compMsg.Handler(s, i)
				return
			}

		// Autocomplete
		case discordgo.InteractionApplicationCommandAutocomplete:
			if h, ok := autocompleteHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}

		// Modal
		case discordgo.InteractionModalSubmit:
			data, res := b.GetData(i.GuildID)
			if !res.Exists {
				return
			}

			handler, exists := data.Modals[i.ModalSubmitData().CustomID]
			if exists {
				handler.Handler(s, i, b.newRespSlash(i))
				return
			}
		}
	})
	b.dg.AddHandler(b.cmdHandler)
	b.dg.AddHandler(b.polls.ReactionHandler)
	b.dg.AddHandler(b.polls.UnReactionHandler)
}
