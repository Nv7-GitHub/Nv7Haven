package main

import (
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/Nv7-Github/Nv7Haven/db"
	"github.com/Nv7-Github/Nv7Haven/elemental"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
)

// TODO: Don't need fiber app, don't return elemental.Elemental object

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

	mysqldb, err := sql.Open("mysql", dbUser+":"+os.Getenv("PASSWORD")+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	if err != nil {
		panic(err)
	}
	db := db.NewDB(mysqldb)

	err = elemental.InitElemental(db, grpcS)
	if err != nil {
		panic(err)
	}

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
}
