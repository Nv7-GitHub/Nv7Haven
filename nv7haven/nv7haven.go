package nv7haven

import (
	"github.com/Nv7-Github/firebase"
	"github.com/gofiber/fiber/v2"

	database "github.com/Nv7-Github/firebase/db"
)

var db *database.Db

// InitNv7Haven initializes the handlers for Nv7Haven
func InitNv7Haven(app *fiber.App) error {
	app.Get("/hella/:input", calcHella)
	fireapp, err := firebase.CreateAppWithServiceAccount("https://nv7haven.firebaseio.com", "AIzaSyA8ySJ5bATo7OADU75TMfbtnvKmx_g5rSs", []byte(serviceAccount))
	if err != nil {
		return err
	}
	db = database.CreateDatabase(fireapp)
	return nil
}
