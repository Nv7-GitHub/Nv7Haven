package single

import (
	"database/sql"
	"os"

	"github.com/gofiber/fiber/v2"
)

func (s *Single) routing(app *fiber.App) {
	app.Post("/single_upload", s.upload)
	app.Get("/single_like/:id/:uid", s.like)
	app.Get("/single_list/:kind/:query", s.list)
	app.Get("/single_list/:kind", s.list)
	app.Get("/single_download/:id/:uid", s.download)
}

// Single is the Nv7 Singleplayer server for elemental 4 (https://elemental4.net)
type Single struct {
	db *sql.DB
}

// InitSingle initializes all of Nv7 Single's handlers on the app.
func InitSingle(app *fiber.App, db *sql.DB) {
	if _, err := os.Stat("/home/container/packs"); os.IsNotExist(err) {
		err = os.Mkdir("/home/container/packs", 0777)
		if err != nil && os.Getenv("MYSQL_HOST") != "c.filipk.in" {
			panic(err)
		}
	}

	s := Single{
		db: db,
	}
	s.routing(app)
}
