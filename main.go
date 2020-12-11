package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Nv7-Github/Nv7Haven/elemental"
	"github.com/Nv7-Github/Nv7Haven/mysqlsetup"
	"github.com/Nv7-Github/Nv7Haven/nv7haven"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Static("/", "./index.html")

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	err := elemental.InitElemental(app)
	if err != nil {
		panic(err)
	}

	err = nv7haven.InitNv7Haven(app)
	if err != nil {
		panic(err)
	}

	mysqlsetup.Mysqlsetup()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	if err := app.Listen(":" + os.Getenv("PORT")); err != nil {
		log.Panic(err)
	}

	elemental.CloseElemental()
}
