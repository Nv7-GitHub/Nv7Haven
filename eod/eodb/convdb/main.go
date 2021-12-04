package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

const (
	dbUser = "root"
	dbName = "nv7haven"
)

var start = time.Now()

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func end(reason string) {
	fmt.Println(reason, "in", time.Since(start))
	start = time.Now()
}

func main() {
	var err error
	db, err = sql.Open("mysql", dbUser+":"+os.Getenv("PASSWORD")+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	if err != nil {
		panic(err)
	}
	end("Connected to DB")

	loadDB(false)
	end("Loaded DB")

	convDB()
	end("Converted DB")
}
