package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

const guild = "705084182673621033"
const user = "663139761128603651"

func main() {
	home, err := os.UserHomeDir()
	handle(err)
	dbPath := filepath.Join(home, "go/src/github.com/Nv7-Github/Nv7haven/data/eod")
	fmt.Println("Loading...")
	start := time.Now()
	dat, err := eodb.NewData(dbPath)
	handle(err)
	fmt.Println("Loaded in", time.Since(start))

	db, _ := dat.GetDB(guild)
	elemsRaw, err := os.ReadFile("inv.txt")
	handle(err)
	inv := db.GetInv(user)

	for _, elem := range strings.Split(string(elemsRaw), "\n") {
		elem, res := db.GetElementByName(elem)
		if !res.Exists {
			fmt.Println(res.Message)
		} else {
			inv.Add(elem.ID)
		}
	}

	db.SaveInv(inv)
	fmt.Println("Done!")
}
