package main

import (
	"log"

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

	log.Fatal(app.Listen(":8080"))
}
