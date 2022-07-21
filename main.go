package main

import (
	_ "embed"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql" // mysql
)

func main() {
	http.ListenAndServe(":"+os.Getenv("MAIN_PORT"), nil)
}
