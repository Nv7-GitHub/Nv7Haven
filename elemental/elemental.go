package elemental

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/iterator"

	firebase "firebase.google.com/go"
	fire "github.com/Nv7-Github/firebase"
	authentication "github.com/Nv7-Github/firebase/auth"
	database "github.com/Nv7-Github/firebase/db"
	"google.golang.org/api/option"
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
	db      *database.Db
	store   *firestore.Client
	auth    *authentication.Auth
	fireapp *firebase.App
	cache   map[string]Element
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
	app.Get("/create_user/:email/:password", e.createUser)
	app.Get("/login_user/:email/:password", e.loginUser)
	app.Get("/clear", func(c *fiber.Ctx) error {
		e.cache = make(map[string]Element, 0)
		return nil
	})
}

func (e *Elemental) init() {
	iter := e.store.Collection("elements").Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err)
		}
		var data Element
		doc.DataTo(&data)
		e.cache[data.Name] = data
	}
}

// InitElemental initializes all of Elemental's handlers on the app.
func InitElemental(app *fiber.App) error {
	opt := option.WithCredentialsJSON([]byte(serviceAccount))
	config := &firebase.Config{
		DatabaseURL:   "https://elementalserver-8c6d0.firebaseio.com",
		ProjectID:     "elementalserver-8c6d0",
		StorageBucket: "elementalserver-8c6d0.appspot.com",
	}
	var err error
	fireapp, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return err
	}

	firebaseapp, err := fire.CreateAppWithServiceAccount("https://elementalserver-8c6d0.firebaseio.com", "AIzaSyCsqvV3clnwDTTgPHDVO2Yatv5JImSUJvU", []byte(serviceAccount))
	if err != nil {
		return err
	}
	auth := authentication.CreateAuth(firebaseapp)

	db := database.CreateDatabase(firebaseapp)

	store, err := fireapp.Firestore(context.Background())
	if err != nil {
		return err
	}

	e := Elemental{
		db:      db,
		auth:    auth,
		store:   store,
		fireapp: fireapp,
		cache:   make(map[string]Element, 0),
	}

	e.routing(app)
	e.init()

	CloseElemental = func() { e.store.Close() }

	return nil
}
