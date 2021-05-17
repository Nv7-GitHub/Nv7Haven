package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql" // mysql
)

func elems2txt(elems []string) string {
	sort.Strings(elems)
	return strings.Join(elems, "+")
}

type empty struct{}

const (
	dbUser = "u57_fypTHIW9t8"
	dbName = "s57_nv7haven"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

// Fixelems fixes the elements
func main() {
	dbPass, err := ioutil.ReadFile("../../password.txt")
	handle(err)
	dbPassword := string(dbPass)

	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	handle(err)
	defer db.Close()
	db.SetMaxOpenConns(10)

	fmt.Println("Connected")

	res, err := db.Query("SELECT elems, guild FROM eod_combos WHERE 1")
	handle(err)
	defer res.Close()

	var combs map[string]empty
	var guild string
	for res.Next() {
		var dat string
		err = res.Scan(&dat, &guild)
		handle(err)
		combs = make(map[string]empty)
		err = json.Unmarshal([]byte(dat), &combs)
		handle(err)

		elems := make([]string, len(combs))
		i := 0
		for k := range combs {
			elems[i] = k
			i++
		}
		if len(elems) == 1 {
			elems = append(elems, elems[0])
		}

		dt := elems2txt(elems)
		_, err = db.Exec("UPDATE eod_combos SET elemsnew=? WHERE elems LIKE ? AND guild=?", dt, dat, guild)
		handle(err)
		fmt.Println(dt)
	}
}
