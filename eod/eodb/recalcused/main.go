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
	dat, err := eodb.NewData(dbPath)
	handle(err)
	fmt.Println("Loaded in", time.Since(start))

	for _, db := range dat.DB {
		// Recalc user used cnt
		userCnts := make(map[string]int)
		for _, el := range db.Elements {
			userCnts[el.Creator]++
		}
		for user, cnt := range userCnts {
			inv := db.GetInv(user)
			inv.UsedCnt += cnt
			handle(db.SaveInv(inv))
		}
		db.Close()
	}
}
