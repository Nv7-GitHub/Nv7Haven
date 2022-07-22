package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/Nv7-Github/Nv7Haven/db"
	"github.com/Nv7-Github/Nv7Haven/eod"

	_ "github.com/go-sql-driver/mysql" // mysql
)

const (
	dbUser = "root"
	dbName = "nv7haven"
)

func main() {
	mysqldb, err := sql.Open("mysql", dbUser+":"+os.Getenv("PASSWORD")+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	if err != nil {
		panic(err)
	}
	db := db.NewDB(mysqldb)

	e := eod.InitEoD(db)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Gracefully shutting down...")
		e.Close()
	}()

	err = http.ListenAndServe(":"+os.Getenv("EOD_PORT"), nil)
	if err != nil {
		panic(err)
	}
}
