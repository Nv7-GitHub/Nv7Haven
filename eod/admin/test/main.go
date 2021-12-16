package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/admin"
	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	start := time.Now()
	fmt.Println("Loading DB...")
	db, err := eodb.NewData("../../../data/eod")
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded in", time.Since(start))

	app := fiber.New()
	app.Use(cors.New())

	admin.InitAdmin(db, app)

	err = app.Listen(":" + os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
}
