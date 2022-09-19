package elemcraft

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type ElemCraft struct {
	app *pocketbase.PocketBase
}

func (e *ElemCraft) Handlers() {
	e.app.OnBeforeServe().Add(func(ev *core.ServeEvent) error {
		ev.Router.AddRoute(echo.Route{
			Method: http.MethodPost,
			Path:   "/api/combine",
			Handler: func(c echo.Context) error {
				var comb [][]int
				dec := json.NewDecoder(c.Request().Body)
				err := dec.Decode(&comb)
				if err != nil {
					return err
				}

				u := c.Get(apis.ContextUserKey).(*models.User)
				inv := u.Profile.Data()["inv"]
				fmt.Println(inv)
				return c.String(200, "Hello world!")
			},
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireUserAuth(),
			},
		})

		ev.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/element",
			Handler: e.GetElement,
		})
		return nil
	})
}

func StartElemCraft() {
	err := os.MkdirAll("data/elemcraft", os.ModePerm)
	if err != nil {
		panic(err)
	}
	p := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: "data/elemcraft",
	})
	e := &ElemCraft{app: p}
	e.Handlers()
	err = e.app.Start()
	if err != nil {
		panic(err)
	}
}
