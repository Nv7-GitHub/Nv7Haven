package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func websockets(app *fiber.App) {
	app.Use("/ws", func(c *fiber.Ctx) error {
		fmt.Println(c.Get("Connection"))
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		fmt.Println("upgrade required")
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		fmt.Println(c.Locals("allowed"))  // true
		fmt.Println(c.Params("id"))       // 123
		fmt.Println(c.Query("v"))         // 1.0
		fmt.Println(c.Cookies("session")) // ""

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				fmt.Println("read:", err)
				break
			}
			fmt.Printf("recv: %s", msg)

			if err = c.WriteMessage(mt, msg); err != nil {
				fmt.Println("write:", err)
				break
			}
		}

	}))

}
