package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Nv7-Github/Nv7Haven/discord"
	"github.com/Nv7-Github/Nv7Haven/elemental"
	"github.com/Nv7-Github/Nv7Haven/nv7haven"
	"github.com/Nv7-Github/Nv7Haven/single"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 1000000000,
	})
	app.Use(cors.New())

	app.Static("/", "./index.html")

	//mysqlsetup.Mysqlsetup()

	e, err := elemental.InitElemental(app)
	if err != nil {
		panic(err)
	}

	err = nv7haven.InitNv7Haven(app)
	if err != nil {
		panic(err)
	}

	single.InitSingle(app)

	b := discord.InitDiscord()

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

	e.Close()
	b.Close()
}
