package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
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
	db, err := eodb.NewData(dbPath)
	handle(err)

	// Get polls
	for _, db := range db.DB {
		fmt.Println(db.Guild, db.Polls)
		if os.Args[1] == db.Guild {
			for _, poll := range db.Polls {
				fmt.Println(poll.Suggestor)
			}
		}
	}

	// Cleanup
	for _, db := range db.DB {
		db.Close()
	}
}
