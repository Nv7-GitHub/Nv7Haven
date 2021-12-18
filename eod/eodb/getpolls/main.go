package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/bwmarrin/discordgo"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	home, err := os.UserHomeDir()
	handle(err)
	dbPath := filepath.Join(home, "go/src/github.com/Nv7-Github/Nv7haven/data/eod")
	fmt.Println("Loading...")
	start := time.Now()
	db, err := eodb.NewData(dbPath)
	handle(err)
	fmt.Println("Loaded in", time.Since(start))

	token, err := os.ReadFile("../../token.txt")
	handle(err)
	dg, err := discordgo.New("Bot " + strings.TrimSpace(string(token)))
	if err != nil {
		panic(err)
	}

	// Get polls
	for _, db := range db.DB {
		fmt.Println(db.Guild, len(db.Polls))
		if len(os.Args) > 1 && os.Args[1] == db.Guild {
			for _, poll := range db.Polls {
				fmt.Println(poll.Suggestor)

				// delete?
				if len(os.Args) > 2 {
					dg.ChannelMessageDelete(poll.Channel, poll.Message)
					db.DeletePoll(poll)
				}
			}
		}

		db.Close()
	}
}
