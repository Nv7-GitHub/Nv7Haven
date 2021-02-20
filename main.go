package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Nv7-Github/Nv7Haven/elemental"
	"github.com/Nv7-Github/Nv7Haven/nv7haven"
	"github.com/Nv7-Github/Nv7Haven/single"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	_ "github.com/go-sql-driver/mysql" // mysql
)

const (
	dbUser     = "u51_iYXt7TBZ0e"
	dbPassword = "W!QnD2u896yo.J4fww9X.h+J"
	dbName     = "s51_nv7haven"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 1000000000,
	})
	app.Use(cors.New())

	/* Testing*/
	websockets(app)

	app.Static("/", "./index.html")

	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	if err != nil {
		panic(err)
	}

	//mysqlsetup.Mysqlsetup()

	e, err := elemental.InitElemental(app, db)
	if err != nil {
		panic(err)
	}

	err = nv7haven.InitNv7Haven(app, db)
	if err != nil {
		panic(err)
	}

	single.InitSingle(app, db)

	//b := discord.InitDiscord(db)

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
	//b.Close()
}
