package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" // mysql
	"golang.org/x/crypto/bcrypt"
)

const (
	dbUser = "root"
	dbName = "nv7haven"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	dbPassword := os.Getenv("PASSWORD")

	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	handle(err)
	defer db.Close()

	fmt.Println("Connected")

	var name string
	fmt.Print("Username: ")
	fmt.Scanln(&name)

	var pwd string
	fmt.Print("New Password: ")
	fmt.Scanln(&pwd)

	hashed, err := bcrypt.GenerateFromPassword([]byte(pwd), 8)
	handle(err)

	_, err = db.Exec("UPDATE users SET password=? WHERE name=?", string(hashed), name)
	handle(err)
}
