package api

import (
	"encoding/json"
	"net/http"

	"github.com/Nv7-Github/Nv7Haven/eod/api/data"
	"github.com/gorilla/websocket"
)

func (a *API) Handle(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	id := ""
	gld := ""

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if id == "" {
			id = string(message)
			continue
		}

		// Parse
		var msg data.Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, data.RSPError(err.Error()).JSON())
			continue
		}

		// Eval
		out := data.RSPBadRequest
		switch msg.Method {
		case data.MethodGuild:
			out = a.MethodGuild(msg.Params, id)
			if out.Error == nil {
				gld = msg.Params["gld"].(string)
			}

		case data.MethodElem:
			out = a.MethodElem(msg.Params, id, gld)

		case data.MethodCombo:
			out = a.MethodCombo(msg.Params, id, gld)

		case data.MethodElemInfo:
			out = a.MethodElemInfo(msg.Params, id, gld)

		case data.MethodInv:
			out = a.MethodInv(id, gld)
		}

		// Respond
		err = conn.WriteMessage(websocket.TextMessage, out.JSON())
		if err != nil {
			return
		}
	}
}
