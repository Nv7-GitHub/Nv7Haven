package nv7haven

import (
	"database/sql"

	"github.com/Nv7-Github/firebase"
	database "github.com/Nv7-Github/firebase/db"
	"github.com/gofiber/fiber/v2"
)

// Nv7Haven is the backend for https://nv7haven.tk
type Nv7Haven struct {
	db  *database.Db
	sql *sql.DB
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
	app.Get("/new_note/:name/:password", c.newNote)
	app.Get("/new_note/:name/", c.newNote)
	app.Post("/change_note/:name/:password", c.changeNote)
	app.Post("/change_note/:name", c.changeNote)
	app.Get("/get_note/:name/:password", c.getNote)
	app.Get("/delete_note/:name/:password", c.deleteNote)
	app.Get("/delete_note/:name", c.deleteNote)
	app.Get("/get_note/:name", c.getNote)
	app.Get("/has_password/:name", c.hasPassword)
	app.Get("/note_search/:query", c.searchNotes)
	app.Get("/search_elems/:query", c.searchElems)
	app.Get("/search_tf/:query/:order", c.searchTf)
	app.Post("/new_tf/:name", c.newTf)
	app.Get("/tf_like/:name", c.like)
	app.Post("/tf_comment/:name", c.comment)
	app.Get("/tf_get/:name", c.getPost)
	app.Post("/upload", c.upload)
	app.Get("/get_file/:id", c.getFile)
	app.Get("/get_ideas/:sort", c.getIdeas)
	app.Get("/new_idea/:title", c.newIdea)
	app.Get("/update_idea/:id/:vote", c.updateIdea)
	app.Get("/breakdown/:input", c.breakdown)
	app.Get("/search_names/:query", c.searchNames)
	app.Get("/get_name/:name", c.getName)
}

// InitNv7Haven initializes the handlers for Nv7Haven
func InitNv7Haven(app *fiber.App, sql *sql.DB) error {
	// Firebase DB
	fireapp, err := firebase.CreateAppWithServiceAccount("https://nv7haven.firebaseio.com", "AIzaSyA8ySJ5bATo7OADU75TMfbtnvKmx_g5rSs", []byte(serviceAccount))
	if err != nil {
		return err
	}
	db := database.CreateDatabase(fireapp)

	nv7haven := Nv7Haven{
		db:  db,
		sql: sql,
	}

	err = nv7haven.initBestEver()
	if err != nil {
		return err
	}
	nv7haven.routing(app)

	return nil
}
