package single

import (
	"os"

	"github.com/Nv7-Github/Nv7Haven/db"
	"github.com/gofiber/fiber/v2"
)

func (s *Single) routing(app *fiber.App) {
	app.Post("/single_upload", s.upload)
	app.Get("/single_like/:id/:uid", s.like)
	app.Get("/single_list/:kind/:query", s.list)
	app.Get("/single_list/:kind", s.list)
	app.Get("/single_download/:id/:uid", s.download)
}

// Single is the Nv7 Singleplayer server for elemental 7 (https://elem7.tk)
type Single struct {
	db *db.DB
}

// InitSingle initializes all of Nv7 Single's handlers on the app.
func InitSingle(app *fiber.App, db *db.DB) {
	err := os.MkdirAll("data/packs", os.ModePerm)
	if err != nil {
		panic(err)
	}

	s := Single{
		db: db,
	}
	s.routing(app)
}
