package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/admin"
	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
)

func main() {
	start := time.Now()
	fmt.Println("Loading DB...")
	db, err := eodb.NewData("../../../data/eod")
	if err != nil {
		panic(err)
	}
	fmt.Println("started in", time.Since(start))

	admin.InitAdmin(db)

	fmt.Println("Listening")
	err = http.ListenAndServe(":"+os.Getenv("HTTP_PORT"), nil)
	if err != nil {
		panic(err)
	}
}
