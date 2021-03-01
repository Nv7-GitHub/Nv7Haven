package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const (
	token    = "Nzg4MTg1MzY1NTMzNTU2NzM2.X9f00g.krA6cjfFWYdzbqOPXq8NvRjxb3k"
	clientID = "788185365533556736"
	guild    = "806258286043070545"
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
		fmt.Println(cmd.ID)
		dg.ApplicationCommandDelete(clientID, guild, cmd.ID)
	}

	dg.Close()
}
