package main

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql" // mysql
	"github.com/r3labs/sse/v2"
)

var events *sse.Server

func main() {
	// SSE
	events = sse.New()
	events.CreateStream("services")
	http.HandleFunc("/events", events.ServeHTTP)
	http.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		v, err := json.Marshal(services)
		if err != nil {
			panic(err) // Should never happen
		}
		w.Write(v)
	})

	// Run
	err := http.ListenAndServe(":"+os.Getenv("MAIN_PORT"), nil)
	if err != nil {
		panic(err)
	}
}
