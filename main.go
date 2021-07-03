package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"

	"github.com/Nv7-Github/Nv7Haven/discord"
	"github.com/Nv7-Github/Nv7Haven/elemental"
	"github.com/Nv7-Github/Nv7Haven/eod"
	"github.com/Nv7-Github/Nv7Haven/gdo"
	"github.com/Nv7-Github/Nv7Haven/nv7haven"
	"github.com/Nv7-Github/Nv7Haven/single"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"

	_ "embed"

	_ "github.com/go-sql-driver/mysql" // mysql
)

const (
	dbUser = "root"
	dbName = "nv7haven"
)

func main() {
	logFile, err := os.OpenFile("logs.txt", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	dupfn(int(logFile.Fd()))

	// Error logging
	//defer recoverer()

	app := fiber.New(fiber.Config{
		BodyLimit: 1000000000,
	})
	app.Use(cors.New())
	app.Use(pprof.New())
	app.Use(recover.New(recover.Config{
		Next:              nil,
		EnableStackTrace:  true,
		StackTraceHandler: traceHandler,
	}))
	app.Get("/freememory", func(c *fiber.Ctx) error {
		debug.FreeOSMemory()
		return nil
	})

	/* Testing*/
	websockets(app)

	app.Static("/", "./index.html")

	db, err := sql.Open("mysql", dbUser+":"+os.Getenv("PASSWORD")+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
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
	gdo.InitGDO(app)

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
