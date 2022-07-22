package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/Nv7-Github/Nv7Haven/db"
	"github.com/Nv7-Github/Nv7Haven/discord"
	"github.com/Nv7-Github/Nv7Haven/elemental"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"

	_ "github.com/go-sql-driver/mysql" // mysql
)

// TODO: Split this up one day

const (
	dbUser = "root"
	dbName = "nv7haven"
)

func main() {
	lis, err := net.Listen("tcp", ":"+os.Getenv("ELEMENTAL_PORT"))
	if err != nil {
		panic(err)
	}
	grpcS := grpc.NewServer()

	app := fiber.New()
	app.Use(cors.New())

	mysqldb, err := sql.Open("mysql", dbUser+":"+os.Getenv("PASSWORD")+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	if err != nil {
		panic(err)
	}
	db := db.NewDB(mysqldb)

	e, err := elemental.InitElemental(app, db, grpcS)
	if err != nil {
		panic(err)
	}

	b := discord.InitDiscord(db, e)

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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Gracefully shutting down...")
		b.Close()
		httpS.Close()
		app.Shutdown()
	}()

	go func() {
		err = httpS.Serve(lis)
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	if err := app.Listen(":" + os.Getenv("LOGIN_PORT")); err != nil {
		panic(err)
	}
}
