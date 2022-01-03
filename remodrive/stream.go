package remodrive

import (
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

func (r *RemoDrive) Drive(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	room := ""
	name := ""

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			// Send leave message
			lock.RLock()
			room, exists := r.Rooms[room]
			lock.RUnlock()
			if exists {
				room.Msgs <- "host_event_leave:" + name
			}
			return
		}

		if room == "" {
			parts := strings.Split(string(message), ":")
			room = parts[0]
			name = parts[1]

			// Save name
			lock.RLock()
			room, exists := r.Rooms[room]
			lock.RUnlock()
			if !exists {
				return
			}

			room.Msgs <- "host_event_join:" + name
			continue
		}

		lock.RLock()
		room, exists := r.Rooms[room]
		lock.RUnlock()
		if !exists {
			return
		}

		room.Msgs <- string(message)

		err = conn.WriteMessage(websocket.TextMessage, []byte("recv"))
		if err != nil {
			room.Msgs <- "host_event_leave:" + name
			return
		}
	}
}

func (r *RemoDrive) Host(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	_, roomName, err := conn.ReadMessage()
	if err != nil {
		return
	}

	lock.RLock()
	room, exists := r.Rooms[string(roomName)]
	lock.RUnlock()
	if !exists {
		return
	}

	for msg := range room.Msgs {
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			return
		}
	}
}
