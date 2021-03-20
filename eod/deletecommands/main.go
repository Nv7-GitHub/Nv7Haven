package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const (
	token    = "ODE5MDc2OTIyODY3NzEyMDMx.YEhW1A.iCZTYR_8YH59k7vlYtUM5LZ8Kn8"
	clientID = "819076922867712031"
	guild    = "819077688371314718"
)

func main() {
	// Discord bot
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	err = dg.Open()
	if err != nil {
		panic(err)
	}

	cmds, err := dg.ApplicationCommands(clientID, guild)
	if err != nil {
		panic(err)
	}

	for _, cmd := range cmds {
		fmt.Println(cmd.ID, cmd.Name)
		if cmd.Name == "elemsort" {
			dg.ApplicationCommandDelete(clientID, guild, cmd.ID)
		}
	}

	dg.Close()
}
