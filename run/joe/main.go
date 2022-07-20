package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	joe "github.com/Nv7-Github/average-joe"
)

//go:embed token.txt
var token string

func main() {
	j, err := joe.NewJoe(token)
	if err != nil {
		panic(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	fmt.Println("Gracefully shutting down...")
	j.Close()
}
