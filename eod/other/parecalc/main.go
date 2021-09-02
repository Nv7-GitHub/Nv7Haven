package main

import (
	"database/sql"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	_ "github.com/go-sql-driver/mysql"
) // mysql

const (
	dbUser = "root"
	dbName = "nv7haven"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

type Element struct {
	ID         int
	Name       string
	Image      string
	Guild      string
	Comment    string
	Creator    string
	CreatedOn  time.Time
	Parents    []string
	Complexity int
	Difficulty int
	UsedIn     int
}

type empty struct{}

type Guild struct {
	Combos   map[int]Combo
	Elements map[string]Element
	Finished map[string]empty
}

func NewGuild() Guild {
	return Guild{
		Combos:   make(map[int]Combo),
		Elements: make(map[string]Element),
		Finished: make(map[string]empty),
	}
}

type Combo struct {
	Elems []string
	Elem3 string
}

var glds = make(map[string]Guild)
var starters = []string{"Air", "Earth", "Fire", "Water"}

var db *sql.DB

var start = time.Now()

func end(reason string) {
	fmt.Println(reason, "in", time.Since(start))
	start = time.Now()
}

func main() {
	prof, err := os.Create("prof.pprof")
	handle(err)
	defer prof.Close()

	err = pprof.StartCPUProfile(prof)
	handle(err)
	defer pprof.StopCPUProfile()
	end("Started profiling")

	db, err = sql.Open("mysql", dbUser+":"+os.Getenv("PASSWORD")+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	if err != nil {
		panic(err)
	}
	end("Connected to DB")
	loadData(true) // false to load from cache
	end("Loaded data")
	recalcPars()
	end("Recalculated parents")
	recalcStats()
	end("Recalculated stats")
	editDB()
}
