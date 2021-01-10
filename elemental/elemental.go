package elemental

import (
	"database/sql"
	"os"

	"github.com/Nv7-Github/firebase"
	"github.com/Nv7-Github/firebase/db"
	database "github.com/Nv7-Github/firebase/db"
	_ "github.com/go-sql-driver/mysql" // mysql
	"github.com/gofiber/fiber/v2"
)

// CloseElemental cleans up elemental
var CloseElemental func()

// Suggestion has the data for a suggestion
type Suggestion struct {
	Creator string   `json:"creator"`
	Name    string   `json:"name"`
	Votes   int      `json:"votes"`
	Color   Color    `json:"color"`
	Voted   []string `json:"voted"`
}

// Recent has the data of a recent element
type Recent struct {
	Parents [2]string `json:"recipe"`
	Result  string    `json:"result"`
}

// Elemental is the "Nv7's Elemental" server at https://elemental4.net, the elemental.json is at https://nv7haven.tk/elemental
type Elemental struct {
	db    *sql.DB
	cache map[string]Element
	fdb   *db.Db
}

func (e *Elemental) init() {
	res, err := e.db.Query("SELECT * FROM elements WHERE 1")
	if err != nil {
		panic(err)
	}
	defer res.Close()
	for res.Next() {
		var elem Element
		elem.Parents = make([]string, 2)
		err = res.Scan(&elem.Name, &elem.Color, &elem.Comment, &elem.Parents[0], &elem.Parents[1], &elem.Creator, &elem.Pioneer, &elem.CreatedOn)
		if err != nil {
			panic(err)
		}
		if (elem.Parents[0] == "") && (elem.Parents[1] == "") {
			elem.Parents = make([]string, 0)
		}
		e.cache[elem.Name] = elem
	}
}

func (e *Elemental) routing(app *fiber.App) {
	app.Get("/get_combo/:elem1/:elem2", e.getCombo)
	app.Get("/get_elem/:elem", e.getElem)
	app.Get("/get_found/:uid", e.getFound)
	app.Get("/new_found/:uid/:elem", e.newFound)
	app.Get("/recents", e.getRecents)
	app.Get("/get_suggestion/:id", e.getSuggestion)
	app.Get("/suggestion_combos/:elem1/:elem2", e.getSuggestionCombos)
	app.Get("/down_suggestion/:id/:uid", e.downVoteSuggestion)
	app.Get("/up_suggestion/:id/:uid", e.upVoteSuggestion)
	app.Get("/create_suggestion/:elem1/:elem2/:id/:mark/:pioneer", e.createSuggestion)
	app.Get("/new_suggestion/:elem1/:elem2/:data", e.newSuggestion)
	app.Get("/create_user/:name/:password", e.createUser)
	app.Get("/login_user/:name/:password", e.loginUser)
	app.Get("/new_anonymous_user", e.newAnonymousUser)
	app.Get("/clear", func(c *fiber.Ctx) error {
		e.cache = make(map[string]Element, 0)
		return nil
	})
}

const (
	dbUser     = "u51_iYXt7TBZ0e"
	dbPassword = "W!QnD2u896yo.J4fww9X.h+J"
	dbName     = "s51_nv7haven"
)

// InitElemental initializes all of Elemental's handlers on the app.
func InitElemental(app *fiber.App) (Elemental, error) {
	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	if err != nil {
		panic(err)
	}

	firebaseapp, err := firebase.CreateAppWithServiceAccount("https://elementalserver-8c6d0.firebaseio.com", "AIzaSyCsqvV3clnwDTTgPHDVO2Yatv5JImSUJvU", []byte(serviceAccount))
	if err != nil {
		return Elemental{}, err
	}

	fdb := database.CreateDatabase(firebaseapp)

	e := Elemental{
		db:    db,
		cache: make(map[string]Element),
		fdb:   fdb,
	}
	e.init()
	e.routing(app)

	return e, nil
}

// Close cleans up elemental
func (e *Elemental) Close() {
	e.db.Close()
}
