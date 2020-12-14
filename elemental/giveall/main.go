package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // mysql
)

const (
	dbUser     = "u29_c99qmCcqZ3"
	dbPassword = "j8@tJ1vv5d@^xMixUqUl+NmA"
	dbName     = "s29_nv7haven"
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

	res, err := db.Query("SELECT name FROM elements WHERE 1")
	handle(err)
	defer res.Close()

	elemNames := make([]string, 0)
	for res.Next() {
		var data string
		err = res.Scan(&data)
		handle(err)
		elemNames = append(elemNames, data)
	}

	data, err := json.Marshal(elemNames)
	handle(err)

	var name string
	fmt.Print("Username: ")
	fmt.Scanln(&name)
	_, err = db.Exec("UPDATE users SET found=? WHERE name=?", data, name)
	handle(err)
}
