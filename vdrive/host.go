package vdrive

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type HostParams struct {
	RoomName string            `json:"room_name"`
	Binds    map[string]string `json:"binds"`
}

func (r *VDrive) Host(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	_, createParams, err := conn.ReadMessage()
	if err != nil {
		return
	}
	var params HostParams
	err = json.Unmarshal(createParams, &params)
	if err != nil {
		return
	}

	room := &Room{
		Events:   make(chan Event),
		Keybinds: params.Binds,
	}

	r.Lock.Lock()
	r.Rooms[params.RoomName] = room
	r.Lock.Unlock()

	for elem := range room.Events {
		v, err := json.Marshal(elem)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = conn.WriteMessage(websocket.TextMessage, v)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
