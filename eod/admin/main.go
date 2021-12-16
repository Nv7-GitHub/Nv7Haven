package admin

import (
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/gofiber/fiber/v2"
)

type empty struct{}

type Admin struct {
	data *eodb.Data

	lock *sync.RWMutex
	uids map[string]empty
}

func (a *Admin) Route(app *fiber.App) {
	app.Post("/admin/config", a.Config)
	app.Post("/admin/login", a.Login)
	app.Post("/admin/setconfig", a.SetConfig)
}

func InitAdmin(data *eodb.Data, app *fiber.App) {
	a := &Admin{
		data: data,
		lock: &sync.RWMutex{},
		uids: make(map[string]empty),
	}
	a.Route(app)
}
