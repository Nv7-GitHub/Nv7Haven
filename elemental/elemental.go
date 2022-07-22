package elemental

import (
	"fmt"
	"sync"
	"time"

	_ "embed"

	"github.com/Nv7-Github/Nv7Haven/db"
	"github.com/Nv7-Github/Nv7Haven/pb"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"google.golang.org/grpc"
)

// Elemental is the default server at https://elem7.tk
type Elemental struct {
	*pb.UnimplementedElementalServer

	db      *db.DB
	cache   map[string]*pb.Element
	recents *sync.Cond
}

func (e *Elemental) init() {
	var cnt int
	err := e.db.QueryRow(`SELECT COUNT(1) FROM elements`).Scan(&cnt)
	if err != nil {
		panic(err)
	}

	fmt.Println("Loading elements...")
	start := time.Now()

	res, err := e.db.Query("SELECT * FROM elements WHERE 1")
	if err != nil {
		panic(err)
	}
	defer res.Close()
	for res.Next() {
		elem := &pb.Element{}
		elem.Parents = make([]string, 2)
		err = res.Scan(&elem.Name, &elem.Color, &elem.Comment, &elem.Parents[0], &elem.Parents[1], &elem.Creator, &elem.Pioneer, &elem.CreatedOn, &elem.Complexity, &elem.Uses, &elem.FoundBy)
		if err != nil {
			panic(err)
		}
		if (elem.Parents[0] == "") && (elem.Parents[1] == "") {
			elem.Parents = make([]string, 0)
		}
		e.cache[elem.Name] = elem
	}
	fmt.Println("Loaded in", time.Since(start))
}

func (e *Elemental) routing(app *fiber.App) {
	app.Post("/create_user/:name", e.createUser)
	app.Post("/login_user/:name", e.loginUser)
	app.Get("/new_anonymous_user", e.newAnonymousUser)
	app.Get("/clear", func(c *fiber.Ctx) error {
		e.cache = make(map[string]*pb.Element)
		e.init()
		return nil
	})

	limit := limiter.New()
	app.Use("/create_user", limit)
	app.Use("/new_anonymous_user", limit)
}

// InitElemental initializes all of Elemental's handlers on the app.
func InitElemental(app *fiber.App, db *db.DB, grpc *grpc.Server) (*Elemental, error) {
	e := &Elemental{
		db:      db,
		cache:   make(map[string]*pb.Element),
		recents: sync.NewCond(&sync.Mutex{}),
	}
	e.init()
	e.routing(app)

	pb.RegisterElementalServer(grpc, e)

	return e, nil
}
