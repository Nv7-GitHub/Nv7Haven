package studier

import (
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
)

type Altruity struct {
	app *pocketbase.PocketBase
}

func (a *Altruity) Init() {}

func Run() {
	err := os.MkdirAll("data/studier", os.ModePerm)
	if err != nil {
		panic(err)
	}
	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: "data/studier",
	})

	a := &Altruity{app: app}
	a.Init()

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
