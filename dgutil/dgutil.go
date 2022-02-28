package dgutil

import (
	"fmt"
	"sort"

	"github.com/bwmarrin/discordgo"
)

func UpdateBotCommands(dg *discordgo.Session, clientID string, guild string, commands []*discordgo.ApplicationCommand) {
	cmds, err := dg.ApplicationCommands(clientID, guild)
	if err != nil {
		panic(err)
	}
	cms := make(map[string]*discordgo.ApplicationCommand)
	for _, cmd := range cmds {
		cms[cmd.Name] = cmd
	}
	idealCmds := make(map[string]*discordgo.ApplicationCommand)
	for _, val := range commands {
		idealCmds[val.Name] = val
		cmd, exists := cms[val.Name]
		if !exists || !commandsAreEqual(cmd, val) {
			_, err := dg.ApplicationCommandCreate(clientID, guild, val)
			if err != nil {
				fmt.Printf("Failed to update command %s\n", val.Name)
			} else {
				fmt.Printf("Updated command %s\n", val.Name)
			}
		}
	}
	for _, cmd := range cmds {
		_, exists := idealCmds[cmd.Name]
		if !exists {
			err = dg.ApplicationCommandDelete(clientID, guild, cmd.ID)
			if err != nil {
				fmt.Printf("Failed to delete command %s\n", cmd.Name)
			} else {
				fmt.Printf("Deleted command %s\n", cmd.Name)
			}
		}
	}
}

func commandsAreEqual(a *discordgo.ApplicationCommand, b *discordgo.ApplicationCommand) bool {
	if a.Name != b.Name || a.Description != b.Description || len(a.Options) != len(b.Options) {
		return false
	}

	return optionsArrEqual(a.Options, b.Options)
}

func optionsArrEqual(a []*discordgo.ApplicationCommandOption, b []*discordgo.ApplicationCommandOption) bool {
	for i, o1 := range a {
		o2 := b[i]
		if o1.Type != o2.Type || o1.Name != o2.Name || o1.Description != o2.Description || len(o1.Choices) != len(o2.Choices) || o1.Autocomplete != o2.Autocomplete || o1.MaxValue != o2.MaxValue || o1.MinValue != nil && o2.MinValue == nil || o1.MinValue == nil && o2.MinValue != nil || (o1.MinValue != nil && o2.MinValue != nil && *o1.MinValue != *o2.MinValue) {
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
		if len(o1.Options) != len(o2.Options) {
			return false
		}
		if len(o1.Options) > 0 {
			if !optionsArrEqual(o1.Options, o2.Options) {
				return false
			}
		}
	}
	return true
}
