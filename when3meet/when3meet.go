package when3meet

import (
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
)

type When3meet struct {
	app *pocketbase.PocketBase
}

func (a *When3meet) Init() {}

func Run() {
	err := os.MkdirAll("data/when3meet", os.ModePerm)
	if err != nil {
		panic(err)
	}
	app := pocketbase.NewWithConfig(&pocketbase.Config{
		DefaultDataDir: "data/when3meet",
	})

	a := &When3meet{app: app}
	a.Init()

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
