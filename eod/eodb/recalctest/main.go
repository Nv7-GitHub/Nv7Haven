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

func main() {
	home, err := os.UserHomeDir()
	handle(err)
	dbPath := filepath.Join(home, "go/src/github.com/Nv7-Github/Nv7haven/data/eod")
	fmt.Println("Loading...")
	start := time.Now()
	db, err := eodb.NewData(dbPath)
	handle(err)
	fmt.Println("Loaded in", time.Since(start))

	f, _ := os.Create("prof.pprof")
	pprof.StartCPUProfile(f)

	d, _ := db.GetDB("705084182673621033")
	fmt.Println("Recalcing...")
	start = time.Now()
	err = d.Recalc()
	fmt.Println(err)
	fmt.Println(time.Since(start))

	f.Close()
	pprof.StopCPUProfile()

	f2, _ := os.Create("heap.pprof")
	pprof.WriteHeapProfile(f2)
	f2.Close()

	for _, db := range db.DB {
		db.Close()
	}
}
