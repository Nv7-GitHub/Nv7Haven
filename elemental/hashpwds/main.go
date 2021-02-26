package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // mysql
	"golang.org/x/crypto/bcrypt"
)

const (
	dbUser     = "u51_iYXt7TBZ0e"
	dbPassword = "W!QnD2u896yo.J4fww9X.h+J"
	dbName     = "s51_nv7haven"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp(c.filipk.in:3306)/"+dbName)
	handle(err)
	defer db.Close()

	fmt.Println("Connected")

	res, err := db.Query("SELECT uid, password FROM users WHERE 1")
	handle(err)
	defer res.Close()
	var uid string
	var password string
	var hashedPassword []byte
	for res.Next() {
		err = res.Scan(&uid, &password)
		handle(err)

		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), 8)
		handle(err)
		fmt.Println(string(hashedPassword))

		_, err = db.Exec("UPDATE users SET password=? WHERE uid=?", string(hashedPassword), uid)
		handle(err)
	}
}
