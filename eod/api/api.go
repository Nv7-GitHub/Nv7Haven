package api

import (
	"net/http"
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
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
	base       *base.Base
	loginLinks map[string]chan string
	loginLock  *sync.RWMutex
}

func NewAPI(data *eodb.Data, base *base.Base) *API {
	return &API{
		Data:       data,
		base:       base,
		loginLinks: make(map[string]chan string),
		loginLock:  &sync.RWMutex{},
	}
}

func (a *API) Run() {
	http.HandleFunc("/eode", a.Handle)
	http.HandleFunc("/eode/oauth", a.HandleOAuth)
}
