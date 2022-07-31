package vdrive

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type VDrive struct {
	Rooms map[string]*Room
	Lock  *sync.Mutex
}

type EventKind string

const (
	EventKindJoin  EventKind = "join"
	EventKindLeave EventKind = "leave"
	EventKindKey   EventKind = "key"
)

type Event struct {
	Kind  EventKind `json:"kind"`
	Value string    `json:"value"`
}

type Room struct {
	Events   chan Event
	Keybinds map[string]string
}

func InitVDrive() {
	r := &VDrive{Lock: &sync.Mutex{}, Rooms: make(map[string]*Room)}
	r.Handlers()
}

func (r *VDrive) Handlers() {
	http.HandleFunc("/host", r.Host)
	http.HandleFunc("/join", r.Join)
}
