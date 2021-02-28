package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"

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

	res, err := db.Query("SELECT name, createdOn FROM elements WHERE 1")
	handle(err)
	defer res.Close()

	elemNames := make([]struct {
		Name      string
		CreatedOn int
	}, 0)
	for res.Next() {
		var name string
		var createdOn int
		err = res.Scan(&name, &createdOn)
		handle(err)
		elemNames = append(elemNames, struct {
			Name      string
			CreatedOn int
		}{name, createdOn})
	}

	sort.Slice(elemNames, func(i, j int) bool { return elemNames[i].CreatedOn < elemNames[j].CreatedOn })

	elemDat := make([]string, len(elemNames))
	for i, val := range elemNames {
		elemDat[i] = val.Name
	}
	data, err := json.Marshal(elemDat)
	handle(err)

	var name string
	fmt.Print("Username: ")
	fmt.Scanln(&name)
	_, err = db.Exec("UPDATE users SET found=? WHERE name=?", data, name)
	handle(err)
}
