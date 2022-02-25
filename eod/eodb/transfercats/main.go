package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

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
	fmt.Println("Loading...")
	start := time.Now()
	db, err := eodb.NewData(dbPath)
	handle(err)
	fmt.Println("Loaded in", time.Since(start))

	// Get polls
	for _, db := range db.DB {
		// Delete inline cat data
		for _, cat := range db.Cats() {
			err := db.SaveCat(cat)
			handle(err)
		}

		db.Close()
	}
}
