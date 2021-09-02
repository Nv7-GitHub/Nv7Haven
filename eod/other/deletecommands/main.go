package main

import (
	"fmt"
	"io/ioutil"

	"github.com/bwmarrin/discordgo"
)

const (
	clientID = "819076922867712031"
	guild    = "" // 819077688371314718
)

type empty struct{}

var toDelete = map[string]empty{
	"setvotes":         {},
	"setpolls":         {},
	"setplaychannel":   {},
	"setvotingchannel": {},
	"setnewschannel":   {},
	"setmodrole":       {},
}

func main() {
	tok, err := ioutil.ReadFile("../../token.txt")
	if err != nil {
		panic(err)
	}
	token := string(tok)

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
		_, exists := toDelete[cmd.Name]
		if exists {
			fmt.Println(cmd.ID, cmd.Name)
			dg.ApplicationCommandDelete(clientID, guild, cmd.ID)
		}
	}

	dg.Close()
}
