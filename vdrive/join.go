package vdrive

import (
	"encoding/json"
	"net/http"
)

type JoinParams struct {
	RoomName string `json:"room_name"`
	Name     string `json:"name"`
}

func (r *VDrive) Join(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	_, createParams, err := conn.ReadMessage()
	if err != nil {
		return
	}
	var params JoinParams
	err = json.Unmarshal(createParams, &params)
	if err != nil {
		return
	}

	// Get room
	r.Lock.Lock()
	room, exists := r.Rooms[params.RoomName]
	r.Lock.Unlock()
	if !exists {
		return
	}

	// Join event
	room.Events <- Event{Kind: EventKindJoin, Value: params.Name}
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			// Leave event
			room.Events <- Event{Kind: EventKindLeave, Value: params.Name}
			return
		}

		// Parse event
		room.Events <- Event{Kind: EventKindKey, Value: string(message)}
	}
}
