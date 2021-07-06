package gdo

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sasha-s/go-deadlock"
)

type GDO struct {
	Streams map[string]*EventStream
	Lock    *deadlock.RWMutex
}

func (g *GDO) routing(app *fiber.App) {
	app.Get("/gdo/events/:uid", g.clientConnect)
	app.Post("/gdo/finish", g.clientResponse)
	app.Post("/gdo/send", g.sendEvent)
	app.Post("/gdo/disconnect", g.clientDisconnect)
}

func InitGDO(app *fiber.App) {
	gdo := GDO{
		Streams: make(map[string]*EventStream),
		Lock:    &deadlock.RWMutex{},
	}
	gdo.routing(app)
}
