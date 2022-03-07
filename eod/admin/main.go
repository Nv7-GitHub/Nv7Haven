package admin

import (
	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/gofiber/fiber/v2"
	"github.com/sasha-s/go-deadlock"
)

type empty struct{}

type Admin struct {
	data *eodb.Data

	lock *deadlock.RWMutex
	uids map[string]empty
}

func (a *Admin) Route(app *fiber.App) {
	app.Post("/admin/login", a.Login)
	app.Post("/admin/config", a.Config)
	app.Post("/admin/setconfig", a.SetConfig)
	app.Post("/admin/element", a.Element)
	app.Post("/admin/setelement", a.SetElement)
}

func InitAdmin(data *eodb.Data, app *fiber.App) {
	a := &Admin{
		data: data,
		lock: &deadlock.RWMutex{},
		uids: make(map[string]empty),
	}
	a.Route(app)
}
