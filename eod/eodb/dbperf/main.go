package main

import (
	"os"
	"path/filepath"
	"runtime/pprof"

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

	// Profiling
	cpu, err := os.Create("cpu.pprof")
	handle(err)
	defer cpu.Close()
	err = pprof.StartCPUProfile(cpu)
	handle(err)

	db, err := eodb.NewData(dbPath)
	handle(err)

	// Profiling
	heap, err := os.Create("heap.pprof")
	handle(err)
	defer heap.Close()
	err = pprof.WriteHeapProfile(heap)
	handle(err)

	// Cleanup
	for _, db := range db.DB {
		db.Close()
	}
}
