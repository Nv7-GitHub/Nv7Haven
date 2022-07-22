package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	fmt.Println("Running!")

	serv := &http.Server{Addr: ":" + os.Getenv("EOD_PORT"), Handler: http.DefaultServeMux}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Gracefully shutting down...")
		e.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		serv.Shutdown(ctx)
		defer cancel()
	}()

	err = serv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
