package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	joe "github.com/Nv7-Github/average-joe"
)

//go:embed token.txt
var token string

func main() {
	fmt.Println("Loading...")
	start := time.Now()
	j, err := joe.NewJoe(token)
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded in", time.Since(start))

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	fmt.Println("Gracefully shutting down...")
	j.Close()
}
