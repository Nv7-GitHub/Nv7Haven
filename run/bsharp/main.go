package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	bsharp "github.com/Nv7-Github/bsharp/bot"
)

//go:embed token.txt
var token string

func main() {
	bsharp, err := bsharp.NewBot("data/bsharp", token, "947593278147162113", "")
	if err != nil {
		panic(err)
	}
	fmt.Println("Running!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	fmt.Println("Gracefully shutting down...")
	bsharp.Close()
}
