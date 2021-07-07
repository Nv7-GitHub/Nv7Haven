package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "embed"

	_ "github.com/go-sql-driver/mysql" // mysql
)

//go:embed inv.txt
var invDat string

const (
	dbUser = "root"
	dbName = "nv7haven"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

type empty struct{}

func main() {
	dbPassword := os.Getenv("PASSWORD")

	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	handle(err)
	defer db.Close()

	fmt.Println("Connected")

	res, err := db.Query("SELECT name FROM eod_elements")
	handle(err)
	defer res.Close()

	names := make(map[string]empty)
	var nm string
	for res.Next() {
		err = res.Scan(&nm)
		handle(err)
		names[strings.ToLower(nm)] = empty{}
	}

	var author string
	var guild string
	fmt.Print("User ID: ")
	fmt.Scanln(&author)
	fmt.Print("Guild: ")
	fmt.Scanln(&guild)

	row := db.QueryRow("SELECT inv FROM eod_inv WHERE guild=? AND user=?", guild, author)
	var invRaw string
	var inv map[string]empty
	err = row.Scan(&invRaw)
	handle(err)
	err = json.Unmarshal([]byte(invRaw), &inv)
	handle(err)

	elems := strings.Split(invDat, "\n")
	for _, elem := range elems {
		_, exists := names[strings.ToLower(elem)]
		if exists {
			inv[strings.ToLower(elem)] = empty{}
		}
	}

	newInv, err := json.Marshal(inv)
	handle(err)
	_, err = db.Exec("UPDATE eod_inv SET inv=? WHERE guild=? AND user=?", string(newInv), guild, author)
	handle(err)
}
