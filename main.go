package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Nv7-Github/Nv7Haven/elemental"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Static("/", "./index.html")

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	elemental.InitElemental(app)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	if err := app.Listen(":8080"); err != nil {
		log.Panic(err)
	}

	elemental.CloseElemental()
}
