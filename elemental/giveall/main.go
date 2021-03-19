package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // mysql
)

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

func main() {
	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp(c.filipk.in:3306)/"+dbName)
	handle(err)
	defer db.Close()

	fmt.Println("Connected")

	res, err := db.Query("SELECT name FROM elements ORDER BY createdOn ASC")
	handle(err)
	defer res.Close()

	elemNames := make([]string, 0)
	for res.Next() {
		var name string
		err = res.Scan(&name)
		handle(err)
		elemNames = append(elemNames, name)
	}
	data, err := json.Marshal(elemNames)
	handle(err)

	var name string
	fmt.Print("Username: ")
	fmt.Scanln(&name)
	_, err = db.Exec("UPDATE users SET found=? WHERE name=?", data, name)
	handle(err)
}
