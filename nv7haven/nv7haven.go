package nv7haven

import (
	"github.com/Nv7-Github/firebase"
	"github.com/gofiber/fiber/v2"

	database "github.com/Nv7-Github/firebase/db"
)

// Nv7Haven is the backend for https://nv7haven.tk
type Nv7Haven struct {
	db *database.Db
}

func (c *Nv7Haven) routing(app *fiber.App) {
	app.Get("/hella/:input", c.calcHella)
	app.Get("/bestever_new_suggest/:suggestion", c.newSuggestion)
	app.Get("/bestever_get_suggest", c.getSuggestion)
	app.Get("/bestever_vote/:item", c.vote)
	app.Get("/bestever_get_ldb/:len", c.getLdb)
	app.Get("/bestever_refresh", c.refresh)
	app.Get("/bestever_mod", c.deleteBad)
	app.Get("/getmyip", c.getIP)
}

// InitNv7Haven initializes the handlers for Nv7Haven
func InitNv7Haven(app *fiber.App) error {
	fireapp, err := firebase.CreateAppWithServiceAccount("https://nv7haven.firebaseio.com", "AIzaSyA8ySJ5bATo7OADU75TMfbtnvKmx_g5rSs", []byte(serviceAccount))
	if err != nil {
		return err
	}
	db := database.CreateDatabase(fireapp)
	nv7haven := Nv7Haven{
		db: db,
	}
	err = nv7haven.initBestEver()
	if err != nil {
		return err
	}
	nv7haven.routing(app)
	return nil
}
