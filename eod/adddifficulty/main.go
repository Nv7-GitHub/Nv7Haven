package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql
)

type empty struct{}

type element struct {
	Name       string
	Categories map[string]empty
	Image      string
	Guild      string
	Comment    string
	Creator    string
	CreatedOn  time.Time
	Parents    []string
	Complexity int
	Difficulty int
}

const (
	dbUser     = "u57_fypTHIW9t8"
	dbPassword = "C7HgI6!GF0NaHCrdUi^tEMGy"
	dbName     = "s57_nv7haven"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

// Fixelems fixes the elements
func main() {
	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	handle(err)
	defer db.Close()
	db.SetMaxOpenConns(10)

	fmt.Println("Connected")

	res, err := db.Query("SELECT name, parent1, parent2, complexity, guild FROM eod_elements WHERE 1")
	handle(err)
	defer res.Close()
	for res.Next() {
		var elem element
		elem.Parents = make([]string, 2)
		err = res.Scan(&elem.Name, &elem.Parents[0], &elem.Parents[1], &elem.Complexity, &elem.Guild)
		handle(err)
		if (elem.Parents[0] == "") && (elem.Parents[1] == "") {
			elem.Parents = make([]string, 0)
		}
		pars := make(map[string]empty)
		for _, val := range elem.Parents {
			pars[strings.ToLower(val)] = empty{}
		}
		data, err := json.Marshal(pars)
		handle(err)
		_, err = db.Exec("UPDATE eod_elements SET parents=? WHERE name=? AND guild=?", string(data), elem.Name, elem.Guild)
		handle(err)
		fmt.Println(elem.Name, string(data))
	}
}
