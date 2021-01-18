package single

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql" // mysql
	"github.com/gofiber/fiber/v2"
)

const (
	dbUser     = "u51_iYXt7TBZ0e"
	dbPassword = "W!QnD2u896yo.J4fww9X.h+J"
	dbName     = "s51_nv7haven"
)

func (s *Single) routing(app *fiber.App) {
	app.Post("/single_upload", s.upload)
	app.Get("/single_like/:id/:uid", s.like)
	app.Get("/single_list", s.list)
}

// Single is the Nv7 Singleplayer server for elemental 4 (https://elemental4.net)
type Single struct {
	db *sql.DB
}

// InitSingle initializes all of Nv7 Single's handlers on the app.
func InitSingle(app *fiber.App) {
	if _, err := os.Stat("packs"); os.IsNotExist(err) {
		err = os.Mkdir("/home/container/packs", os.ModeDir)
		if err != nil {
			panic(err)
		}
	}

	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	if err != nil {
		panic(err)
	}

	s := Single{
		db: db,
	}
	s.routing(app)
}
