package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"

	"github.com/r3labs/sse/v2"
)

var events *sse.Server

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	// Fill out build cache
	err := os.MkdirAll("build", os.ModePerm)
	if err != nil {
		panic(err)
	}
	files, err := os.ReadDir("build")
	if err != nil {
		panic(err)
	}
	needed := make(map[string]struct{})
	for _, serv := range services {
		needed[serv.ID] = struct{}{}
	}
	for _, file := range files {
		delete(needed, file.Name())
	}
	for k := range needed {
		fmt.Printf("Building %s...\n", k)
		err = Build(services[k])
		if err != nil {
			panic(err)
		}
	}

	// SSE
	events = sse.New()
	events.CreateStream("services")
	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		events.ServeHTTP(w, r)
	})
	http.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		enableCors(&w)
		w.Write(marshalServices())
	})

	// Run
	fmt.Println("Listening on port", os.Getenv("MAIN_PORT"))
	err = http.ListenAndServe(":"+os.Getenv("MAIN_PORT"), nil)
	if err != nil {
		panic(err)
	}
}
