package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

type empty struct{}

func main() {
	home, err := os.UserHomeDir()
	handle(err)
	dbPath := filepath.Join(home, "go/src/github.com/Nv7-Github/Nv7haven/data/eod")
	fmt.Println("Loading...")
	start := time.Now()
	db, err := eodb.NewData(dbPath)
	handle(err)
	fmt.Println("Loaded in", time.Since(start))

	f, err := os.Create("prof.pprof")
	if err != nil {
		panic(err)
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		panic(err)
	}

	for _, db := range db.DB {
		fmt.Printf("Recalcing %s...\n", db.Guild)
		start = time.Now()
		err = db.Recalc()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(time.Since(start))

		db.Close()
	}

	pprof.StopCPUProfile()
	f.Close()

	f2, _ := os.Create("heap.pprof")
	pprof.WriteHeapProfile(f2)
	f2.Close()
}
