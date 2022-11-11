package main

import (
	_ "embed"
	"os"

	"github.com/Nv7-Github/Nv7Haven/eod"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq" // postgres
)

//go:embed token.txt
var token string

const (
	dbUser = "postgres"
	dbName = "eod"
	dbPort = "5432"
)

func main() {
	db, err := sqlx.Connect("postgres", "user="+dbUser+" dbname="+dbName+" sslmode=disable port="+dbPort+"host="+os.Getenv("MYSQL_HOST")+" password="+os.Getenv("PASSWORD"))
	if err != nil {
		panic(err)
	}

	eod.InitEod(db, token)
}
