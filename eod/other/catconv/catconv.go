package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql
)

type empty struct{}

type category struct {
	Name     string
	Guild    string
	Elements map[string]empty
	Image    string
}

var tm = time.Now()

func endTimer(print string) {
	newtime := time.Now()
	fmt.Println(print, "in", newtime.Sub(tm))
	tm = newtime
}

const (
	dbUser    = "root"
	dbName    = "nv7haven"
	batchSize = 10000
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	dbPass, err := ioutil.ReadFile("../../password.txt")
	handle(err)
	dbPassword := string(dbPass)

	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	handle(err)
	defer db.Close()
	db.SetMaxOpenConns(10)

	endTimer("Connected to database")

	res, err := db.Query("SELECT DISTINCT guild FROM eod_elements")
	handle(err)
	defer res.Close()

	endTimer("Queried data")

	var guild string
	for res.Next() {
		err = res.Scan(&guild)
		handle(err)
		doGuild(guild, db)
	}
}

func doGuild(guild string, db *sql.DB) {
	cats := make(map[string]category)

	elems, err := db.Query("SELECT name, categories FROM eod_elements WHERE guild=?", guild)
	handle(err)
	defer elems.Close()

	var elName string
	var catDat string
	var elCats map[string]empty
	for elems.Next() {
		err = elems.Scan(&elName, &catDat)
		handle(err)
		elCats = make(map[string]empty)
		err = json.Unmarshal([]byte(catDat), &elCats)
		handle(err)

		for cat := range elCats {
			_, exists := cats[strings.ToLower(cat)]
			if !exists {
				cats[strings.ToLower(cat)] = category{
					Name:     cat,
					Guild:    guild,
					Elements: make(map[string]empty),
					Image:    "",
				}
			}

			categ := cats[strings.ToLower(cat)]
			categ.Elements[elName] = empty{}
			cats[strings.ToLower(cat)] = categ
		}
	}

	endTimer("Read Data for Guild " + guild)

	query := "INSERT INTO eod_categories VALUES "
	times := 0
	args := make([]interface{}, 0)
	for _, cat := range cats {
		query += "(?,?,?,?),"

		els, err := json.Marshal(cat.Elements)
		handle(err)

		args = append(args, cat.Guild, cat.Name, string(els), cat.Image)

		if (times % batchSize) == 0 {
			query = query[:len(query)-1]
			_, err = db.Exec(query, args...)
			handle(err)
			endTimer(fmt.Sprintf("Processed and wrote %d records to SQL database", batchSize))

			query = "INSERT INTO eod_categories VALUES "
			args = make([]interface{}, 0)
			times = 0
		}

		times++
	}

	if len(args) > 0 {
		query = query[:len(query)-1]
		_, err = db.Exec(query, args...)
		handle(err)
		endTimer("Processed and wrote final records to SQL database")
	}
}
