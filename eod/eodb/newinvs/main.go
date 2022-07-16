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
	// Load
	home, err := os.UserHomeDir()
	handle(err)
	dbPath := filepath.Join(home, "go/src/github.com/Nv7-Github/Nv7haven/data/eod")

	// Create dirs
	dirs, err := os.ReadDir(dbPath)
	handle(err)
	for _, dir := range dirs {
		if dir.IsDir() {
			path := filepath.Join(dbPath, dir.Name())
			err = os.MkdirAll(filepath.Join(path, "invdata"), os.ModePerm)
			handle(err)
			invs, err := os.ReadDir(filepath.Join(path, "inventories"))
			handle(err)
			for _, inv := range invs {
				if !inv.IsDir() {
					f, err := os.Create(filepath.Join(path, "invdata", inv.Name()))
					handle(err)
					f.Close()
				}
			}
		}
	}

	fmt.Println("Loading DB...")
	start := time.Now()
	dat, err := eodb.NewData(dbPath)
	handle(err)
	fmt.Println("Loaded in", time.Since(start))

	for _, db := range dat.DB {
		for _, inv := range db.Invs() {
			fmt.Println(len(inv.Elements))
			db.SaveInv(inv)
		}
	}
}
