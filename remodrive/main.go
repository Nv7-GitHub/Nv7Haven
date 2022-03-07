package remodrive

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/websocket"
	"github.com/sasha-s/go-deadlock"
)

var lock = &deadlock.RWMutex{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type RemoDrive struct {
	Rooms map[string]Room
}

type Room struct {
	Msgs chan string
}

func InitRemoDrive(app *fiber.App) {
	rd := &RemoDrive{}
	rd.Rooms = make(map[string]Room)
	rd.Handlers(app)
}

func (r *RemoDrive) Handlers(app *fiber.App) {
	app.Post("/close_room", func(ctx *fiber.Ctx) error {
		return r.CloseRoomByName(string(ctx.Body()))
	})
	app.Post("/new_room", func(ctx *fiber.Ctx) error {
		r.NewRoom(string(ctx.Body()))
		return nil
	})
	http.HandleFunc("/drive", r.Drive)
	http.HandleFunc("/host", r.Host)
}
