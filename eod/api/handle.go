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

	// Get ID
	res, state := a.GenURL()
	err = conn.WriteMessage(websocket.TextMessage, res.JSON())
	if err != nil {
		return
	}
	if res.Error != nil {
		return
	}

	// Wait for ID
	a.loginLock.RLock()
	ch := a.loginLinks[state]
	a.loginLock.RUnlock()
	id := <-ch
	a.loginLock.Lock()
	delete(a.loginLinks, state)
	a.loginLock.Unlock()
	close(ch)
	err = conn.WriteMessage(websocket.TextMessage, []byte(id))
	if err != nil {
		return
	}

	gld := ""

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return
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

		case data.MethodCategory:
			out = a.MethodCategory(msg.Params, id, gld)
		}

		// Respond
		err = conn.WriteMessage(websocket.TextMessage, out.JSON())
		if err != nil {
			return
		}
	}
}
