package main

import (
	"database/sql"
	"os"

	"github.com/Nv7-Github/Nv7Haven/db"
	"github.com/Nv7-Github/Nv7Haven/nv7haven"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

const (
	dbUser = "root"
	dbName = "nv7haven"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 1000000000,
	})
	app.Use(cors.New())

	mysqldb, err := sql.Open("mysql", dbUser+":"+os.Getenv("PASSWORD")+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	if err != nil {
		panic(err)
	}
	db := db.NewDB(mysqldb)

	err = nv7haven.InitNv7Haven(app, db)
	if err != nil {
		panic(err)
	}

	if err := app.Listen(":" + os.Getenv("NV7HAVEN_PORT")); err != nil {
		panic(err)
	}
}
