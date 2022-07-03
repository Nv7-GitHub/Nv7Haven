package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

const element = 179745 // Intelligent Unit
const guild = "705084182673621033"

func main() {
	home, err := os.UserHomeDir()
	handle(err)
	dbPath := filepath.Join(home, "go/src/github.com/Nv7-Github/Nv7haven/data/eod")
	fmt.Println("Loading...")
	start := time.Now()
	db, err := eodb.NewData(dbPath)
	handle(err)
	fmt.Println("Loaded in", time.Since(start))
	gld, _ := db.GetDB(guild)

	// Get elem
	els := make(map[int]types.Element)
	addEl(els, element, gld)

	// Save
	out, err := os.Create("export.json")
	handle(err)
	defer out.Close()
	enc := json.NewEncoder(out)
	err = enc.Encode(els)
	handle(err)
}

func addEl(tree map[int]types.Element, el int, db *eodb.DB) {
	_, exists := tree[el]
	if exists {
		return
	}

	elem, _ := db.GetElement(el)
	tree[elem.ID] = elem
	for _, par := range elem.Parents {
		addEl(tree, par, db)
	}
}
