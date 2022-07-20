package main

import (
	"net/http"
	"os"

	"github.com/Nv7-Github/Nv7Haven/remodrive"
)

// TODO: Move handlers using fiber app over to HTTP

func main() {
	remodrive.InitRemoDrive()

	err := http.ListenAndServe(":"+os.Getenv("REMODRIVE_PORT"), nil)
	if err != nil {
		panic(err)
	}
}
