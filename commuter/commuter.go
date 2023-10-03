package commuter

import (
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
)

type Commuter struct {
	app *pocketbase.PocketBase
}

func (a *Commuter) Init() {}

func Run() {
	err := os.MkdirAll("data/commuter", os.ModePerm)
	if err != nil {
		panic(err)
	}
	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: "data/commuter",
	})

	a := &Commuter{app: app}
	a.Init()

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
