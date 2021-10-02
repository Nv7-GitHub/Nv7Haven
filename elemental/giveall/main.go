package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" // mysql
)

const (
	dbUser = "root"
	dbName = "nv7haven"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	dbPassword := os.Getenv("PASSWORD")

	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	handle(err)
	defer db.Close()

	fmt.Println("Connected")

	res, err := db.Query("SELECT name FROM elements ORDER BY createdOn ")
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
