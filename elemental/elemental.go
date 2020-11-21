package elemental

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"

	firebase "firebase.google.com/go"
	database "firebase.google.com/go/db"
	"google.golang.org/api/option"
)

var db *database.Client
var store *firestore.Client

// Element has the data for a created element
type Element struct {
	Color     string   `json:"color"`
	Comment   string   `json:"comment"`
	CreatedOn int      `json:"createdOn"`
	Creator   string   `json:"creator"`
	Name      string   `json:"name"`
	Parents   []string `json:"parents"`
	Pioneer   string   `json:"pioneer"`
}

// Color has the data for a suggestion's color
type Color struct {
	Base       string  `json:"base"`
	Lightness  float32 `json:"lightness"`
	Saturation float32 `json:"saturation"`
}

// Suggestion has the data for a suggestion
type Suggestion struct {
	Creator string   `json:"creator"`
	Name    string   `json:"name"`
	Votes   int      `json:"votes"`
	Color   Color    `json:"color"`
	Voted   []string `json:"voted"`
}

// ComboMap has the data that maps combos
type ComboMap map[string]map[string]string

// SuggMap has the data that maps suggestion combos
type SuggMap map[string]map[string][]string

// Recent has the data of a recent element
type Recent struct {
	Parents [2]string `json:"recipe"`
	Result  string    `json:"result"`
}

// InitElemental initializes all of Elemental's handlers on the app.
func InitElemental(app *fiber.App) error {
	opt := option.WithCredentialsJSON([]byte(json))
	config := &firebase.Config{
		DatabaseURL:   "https://elementalserver-8c6d0.firebaseio.com",
		ProjectID:     "elementalserver-8c6d0",
		StorageBucket: "elementalserver-8c6d0.appspot.com",
	}
	fireapp, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return err
	}

	db, err = fireapp.Database(context.Background())
	if err != nil {
		return err
	}

	store, err = fireapp.Firestore(context.Background())
	if err != nil {
		return err
	}

	app.Get("/get_combo/:elem1/:elem2", getCombo)
	app.Get("/get_elem/:elem", getElem)
	app.Get("/clear", func(c *fiber.Ctx) error {
		cache = make(map[string]Element, 0)
		elemMap = make(map[string]map[string]string, 0)
		return nil
	})
	return nil
}

// CloseElemental has the cleanup functions
func CloseElemental() {
	store.Close()
}
