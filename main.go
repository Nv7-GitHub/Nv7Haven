package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/r3labs/sse/v2"
)

var events *sse.Server

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

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

	// HTTP Server
	m := mux.NewRouter()

	// Service list
	events = sse.New()
	events.CreateStream("services")
	m.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		events.ServeHTTP(w, r)
	})
	m.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		enableCors(&w)
		w.Write(marshalServices())
	})

	// Stream output
	m.HandleFunc("/logs/{id}", func(w http.ResponseWriter, r *http.Request) {
		// Get service
		vars := mux.Vars(r)
		id := vars["id"]
		var serv *Service
		lock.Lock()
		for _, servic := range services {
			if servic.ID == id {
				serv = servic
				break
			}
		}
		lock.Unlock()
		if serv == nil {
			http.NotFound(w, r)
			return
		}

		// Stream
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		// Send current logs
		err = c.WriteMessage(websocket.TextMessage, []byte(serv.Output.Content.String()))
		if err != nil {
			return
		}

		// Start listening for new logs
		for {
			serv.Output.Cond.Wait()
			err = c.WriteMessage(websocket.TextMessage, serv.Output.Data)
			if err != nil {
				break
			}
		}
	})

	// Start services
	for _, serv := range services {
		err = Run(serv)
		if err != nil {
			panic(err)
		}
	}

	// Run
	fmt.Println("Listening on port", os.Getenv("MAIN_PORT"))
	err = http.ListenAndServe(":"+os.Getenv("MAIN_PORT"), m)
	if err != nil {
		panic(err)
	}
}
