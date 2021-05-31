package gdo

import (
	"sync"

	"github.com/gofiber/fiber/v2"
)

type GDO struct {
	Streams map[string]*EventStream
	Lock    *sync.RWMutex
}

func (g *GDO) routing(app *fiber.App) {
	app.Get("/gdo/events/:uid", g.clientConnect)
	app.Post("/gdo/finish", g.clientResponse)
	app.Post("/gdo/send", g.sendEvent)
}

func InitGDO(app *fiber.App) {
	gdo := GDO{
		Streams: make(map[string]*EventStream),
		Lock:    &sync.RWMutex{},
	}
	gdo.routing(app)
}
