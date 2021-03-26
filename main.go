package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Nv7-Github/Nv7Haven/discord"
	"github.com/Nv7-Github/Nv7Haven/elemental"
	"github.com/Nv7-Github/Nv7Haven/eod"
	"github.com/Nv7-Github/Nv7Haven/nv7haven"
	"github.com/Nv7-Github/Nv7Haven/single"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"

	_ "github.com/go-sql-driver/mysql" // mysql
)

const (
	dbUser     = "u57_fypTHIW9t8"
	dbPassword = "C7HgI6!GF0NaHCrdUi^tEMGy"
	dbName     = "s57_nv7haven"
)

func main() {
	// Error logging
	logFile, err := os.OpenFile("logs.txt", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	syscall.Dup2(int(logFile.Fd()), 2)

	app := fiber.New(fiber.Config{
		BodyLimit: 1000000000,
	})
	app.Use(cors.New())
	app.Use(pprof.New())

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
	b := discord.InitDiscord(db, e)
	eod := eod.InitEoD(db)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Gracefully shutting down...")
		app.Shutdown()
	}()

	if err := app.Listen(":" + os.Getenv("PORT")); err != nil {
		panic(err)
	}

	e.Close()
	b.Close()
	eod.Close()
	db.Close()
}
