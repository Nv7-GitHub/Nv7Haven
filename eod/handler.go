package eod

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const guild = "" // 819077688371314718 for testing

func (b *EoD) initHandlers() {
	// Debugging
	var err error
	datafile, err = os.OpenFile("createlogs.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	discordlogs, err = os.OpenFile("discordlogs.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	discordgo.Logger = func(msgL, caller int, format string, a ...interface{}) {
		// This code is a slightly modified version of https://github.com/bwmarrin/discordgo/blob/577e7dd4f6ccf1beb10acdb1871300c7638b84c4/logging.go#L46
		pc, file, line, _ := runtime.Caller(caller)

		files := strings.Split(file, "/")
		file = files[len(files)-1]

		name := runtime.FuncForPC(pc).Name()
		fns := strings.Split(name, ".")
		name = fns[len(fns)-1]

		msg := fmt.Sprintf(format, a...)

		log.SetOutput(discordlogs)
		log.Printf("[DG%d] %s:%d:%s() %s\n", msgL, file, line, name, msg)
	}

	// Handlers
	b.initInfoChoices()

	cmds, err := b.dg.ApplicationCommands(clientID, guild)
	if err != nil {
		panic(err)
	}
	cms := make(map[string]*discordgo.ApplicationCommand)
	for _, cmd := range cmds {
		cms[cmd.Name] = cmd
	}
	for _, val := range commands {
		if val.Name == "elemsort" {
			val.Options[0].Choices = infoChoices
		}
		cmd, exists := cms[val.Name]
		if !exists || !commandsAreEqual(cmd, val) {
			_, err := b.dg.ApplicationCommandCreate(clientID, guild, val)
			if err != nil {
				fmt.Printf("Failed to update command %s\n", val.Name)
			} else {
				fmt.Printf("Updated command %s\n", val.Name)
			}
		}
	}

	b.dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Member == nil {
			return
		}

		// Command
		if i.Type == discordgo.InteractionApplicationCommand {
			rsp := b.newRespSlash(i)
			canRun, msg := b.canRunCmd(i)
			if !canRun {
				rsp.ErrorMessage(msg)
				return
			}

			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
			return
		}

		// Button
		if i.Type == discordgo.InteractionMessageComponent {
			lock.Lock()
			dat, exists := b.dat[i.GuildID]
			if !exists {
				return
			}
			lock.Unlock()

			// Check if page switch handler or component handler
			_, exists = dat.PageSwitchers[i.Message.ID]
			if exists {
				b.pageSwitchHandler(s, i)
				return
			}

			compMsg, exists := dat.ComponentMsgs[i.Message.ID]
			if exists {
				compMsg.Handler(s, i)
				return
			}
			return
		}
	})
	b.dg.AddHandler(b.cmdHandler)
	b.dg.AddHandler(b.reactionHandler)
	b.dg.AddHandler(b.unReactionHandler)
}

func commandsAreEqual(a *discordgo.ApplicationCommand, b *discordgo.ApplicationCommand) bool {
	if a.Name != b.Name || a.Description != b.Description || len(a.Options) != len(b.Options) {
		return false
	}
	for i, o1 := range a.Options {
		o2 := b.Options[i]
		if o1.Type != o2.Type || o1.Name != o2.Name || o1.Description != o2.Description || len(o1.Choices) != len(o2.Choices) {
			return false
		}
		sort.Slice(o1.Choices, func(i, j int) bool {
			return o1.Choices[i].Name < o1.Choices[j].Name
		})
		sort.Slice(o2.Choices, func(i, j int) bool {
			return o2.Choices[i].Name < o2.Choices[j].Name
		})
		for i, c1 := range o1.Choices {
			c2 := o2.Choices[i]
			if c1.Name != c2.Name || fmt.Sprintf("%v", c1.Value) != fmt.Sprintf("%v", c2.Value) {
				return false
			}
		}
	}
	return true
}
