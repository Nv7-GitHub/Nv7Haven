package remodrive

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func (r *RemoDrive) Drive(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	room := ""

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if room == "" {
			room = string(message)
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
			return
		}
	}
}

func (r *RemoDrive) Host(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	fmt.Println(err)
	if err != nil {
		return
	}
	defer conn.Close()

	_, roomName, err := conn.ReadMessage()
	fmt.Println(err)
	if err != nil {
		return
	}

	lock.RLock()
	room, exists := r.Rooms[string(roomName)]
	lock.RUnlock()
	fmt.Println(exists)
	if !exists {
		return
	}

	for msg := range room.Msgs {
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
		fmt.Println(err)
		if err != nil {
			return
		}
	}
}
