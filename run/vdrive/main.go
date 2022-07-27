package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Nv7-Github/Nv7Haven/vdrive"
)

func main() {
	vdrive.InitVDrive()

	fmt.Println("Listening on port " + os.Getenv("VDRIVE_PORT"))
	err := http.ListenAndServe(":"+os.Getenv("VDRIVE_PORT"), nil)
	if err != nil {
		panic(err)
	}
}
