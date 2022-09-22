package elemcraft

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type ElemCraft struct {
	app *pocketbase.PocketBase
}

func (e *ElemCraft) Handlers() {
	e.app.OnBeforeServe().Add(func(ev *core.ServeEvent) error {
		ev.Router.AddRoute(echo.Route{
			Method:  http.MethodPost,
			Path:    "/api/combine",
			Handler: e.Combo,
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireUserAuth(),
			},
		})

		ev.Router.AddRoute(echo.Route{
			Method:  http.MethodGet,
			Path:    "/api/element",
			Handler: e.GetElement,
		})

		ev.Router.AddRoute(echo.Route{
			Method:  http.MethodPost,
			Path:    "/api/suggest",
			Handler: e.Suggest,
			Middlewares: []echo.MiddlewareFunc{
				apis.RequireUserAuth(),
			},
		})

		ev.Router.AddRoute(echo.Route{
			Method:  http.MethodPost,
			Path:    "/api/existing",
			Handler: e.Existing,
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
