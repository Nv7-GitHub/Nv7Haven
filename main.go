package main

import (
	"context"
	_ "embed"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const uidLength = 16

var uids = make(map[string]struct{})

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func getServ(r *http.Request) (serv *Service) {
	vars := mux.Vars(r)
	id := vars["id"]
	lock.Lock()
	for _, servic := range services {
		if servic.ID == id {
			serv = servic
			break
		}
	}
	lock.Unlock()
	return
}

func checkUID(r *http.Request) bool {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false
	}
	_, res := uids[string(body)]
	return res
}

const printerr = false

func main() {
	// HTTP Server
	m := mux.NewRouter()

	// Service list
	m.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		// Stream
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		// Send current services
		err = c.WriteMessage(websocket.TextMessage, marshalServices())
		if err != nil {
			return
		}

		// Listen for new services
		for {
			servicesCond.L.Lock()
			servicesCond.Wait()
			servicesCond.L.Unlock()
			err = c.WriteMessage(websocket.TextMessage, marshalServices())
			if err != nil {
				return
			}
		}
	})

	// Stream output
	m.HandleFunc("/logs/{id}", func(w http.ResponseWriter, r *http.Request) {
		// Get service
		serv := getServ(r)
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
			serv.Output.Cond.L.Lock()
			serv.Output.Cond.Wait()
			serv.Output.Cond.L.Unlock()
			err = c.WriteMessage(websocket.TextMessage, serv.Output.Data)
			if err != nil {
				break
			}
		}
	})

	// Rebuild/Stop/Start
	m.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		// Get body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if matches secret
		if string(body) == os.Getenv("MOD_PASSWORD") {
			rand.Seed(time.Now().UnixNano())
			b := make([]byte, uidLength)
			rand.Read(b)
			uid := fmt.Sprintf("%x", b)[:uidLength]
			uids[uid] = struct{}{}
			w.Write([]byte(uid))
		} else {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
		}
	}).Methods("POST")

	m.HandleFunc("/rebuild/{id}", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		// Get service
		serv := getServ(r)
		if serv == nil {
			http.NotFound(w, r)
			return
		}
		if !checkUID(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Rebuild
		err := Build(serv)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("OK"))
	}).Methods("POST")

	m.HandleFunc("/stop/{id}", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		// Get service
		serv := getServ(r)
		if serv == nil {
			http.NotFound(w, r)
			return
		}
		if !checkUID(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Rebuild
		err := Stop(serv)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("OK"))
	}).Methods("POST")

	m.HandleFunc("/start/{id}", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		// Get service
		serv := getServ(r)
		if serv == nil {
			http.NotFound(w, r)
			return
		}
		if !checkUID(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Rebuild
		err := Run(serv)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("OK"))
	}).Methods("POST")

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

	// Start services
	for _, serv := range services {
		err = Run(serv)
		if err != nil {
			panic(err)
		}
	}

	// Handle shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	server := &http.Server{Addr: ":" + os.Getenv("MAIN_PORT"), Handler: m}
	go func() {
		<-c
		fmt.Println("Shutting down...")

		for _, serv := range services {
			err = Stop(serv)
			if err != nil {
				fmt.Println(err)
			}
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	// Run
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://nv7haven.com", http.StatusMovedPermanently)
	})
	fmt.Println("Listening on port", os.Getenv("MAIN_PORT"))
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
