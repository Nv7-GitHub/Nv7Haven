package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/Nv7-Github/Nv7Haven/anarchy"
	"github.com/Nv7-Github/Nv7Haven/db"
	"github.com/Nv7-Github/Nv7Haven/discord"
	"github.com/Nv7-Github/Nv7Haven/elemental"
	"github.com/Nv7-Github/Nv7Haven/eod"
	"github.com/Nv7-Github/Nv7Haven/gdo"
	"github.com/Nv7-Github/Nv7Haven/nv7haven"
	"github.com/Nv7-Github/Nv7Haven/remodrive"
	"github.com/Nv7-Github/Nv7Haven/single"
	joe "github.com/Nv7-Github/average-joe"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"

	_ "embed"

	_ "github.com/go-sql-driver/mysql" // mysql
)

const (
	dbUser = "root"
	dbName = "nv7haven"
)

//go:embed joe_token.txt
var joe_token string

func main() {
	logFile, err := os.OpenFile("logs.txt", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	dupfn(int(logFile.Fd()))

	// Error logging
	//defer recoverer()

	// Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: 1000000000,
	})
	app.Use(cors.New())
	app.Use(pprof.New())
	systemHandlers(app)

	// gRPC
	lis, err := net.Listen("tcp", ":"+os.Getenv("RPC_PORT"))
	if err != nil {
		panic(err)
	}
	grpcS := grpc.NewServer()

	app.Static("/", "./index.html")

	mysqldb, err := sql.Open("mysql", dbUser+":"+os.Getenv("PASSWORD")+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	if err != nil {
		panic(err)
	}
	db := db.NewDB(mysqldb)

	//mysqlsetup.Mysqlsetup()

	e, err := elemental.InitElemental(app, db, grpcS)
	if err != nil {
		panic(err)
	}

	err = nv7haven.InitNv7Haven(app, db)
	if err != nil {
		panic(err)
	}

	j, err := joe.NewJoe(joe_token)
	if err != nil {
		panic(err)
	}

	single.InitSingle(app, db)
	b := discord.InitDiscord(db, e)
	eodB := eod.InitEoD(db)
	anarchy.InitAnarchy(db, grpcS)
	gdo.InitGDO(app)
	remodrive.InitRemoDrive(app)

	go func() {
		err := http.ListenAndServe(":"+os.Getenv("HTTP_PORT"), nil)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		wrapped := grpcweb.WrapServer(grpcS)
		httpS := &http.Server{
			Handler: http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				// CORS
				resp.Header().Set("Access-Control-Allow-Origin", "*")
				resp.Header().Set("Access-Control-Allow-Methods", "*")
				resp.Header().Set("Access-Control-Allow-Headers", "*")
				wrapped.ServeHTTP(resp, req)
			}),
		}
		defer httpS.Close()

		err = httpS.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", ":"+os.Getenv("GRPC_PORT"))
		if err != nil {
			panic(err)
		}
		err = grpcS.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

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

	b.Close()
	eodB.Close()
	db.Close()
	j.Close()
}
