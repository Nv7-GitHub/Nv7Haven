package api

import (
	"net/http"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type API struct {
	*eodb.Data
}

func NewAPI(data *eodb.Data) *API {
	return &API{
		Data: data,
	}
}

func (a *API) Run() {
	http.HandleFunc("/eode", a.Handle)
}
