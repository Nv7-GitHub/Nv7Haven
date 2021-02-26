package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"

	_ "github.com/go-sql-driver/mysql" // mysql
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

	res, err := db.Query("SELECT * FROM element_combos WHERE 1")
	handle(err)
	defer res.Close()
	var name string
	var combodat string
	var combos map[string]string
	var combs map[string]map[string]string = make(map[string]map[string]string)
	for res.Next() {
		combos = make(map[string]string)
		err = res.Scan(&name, &combodat)
		handle(err)
		err = json.Unmarshal([]byte(combodat), &combos)
		handle(err)
		for k, v := range combos {
			thing := []string{name, k}
			sort.Strings(thing)
			_, exists := combs[thing[0]]
			if exists {
				_, exists = combs[thing[0]][thing[1]]
				if !exists {
					combs[thing[0]][thing[1]] = v
				}
			} else {
				combs[thing[0]] = make(map[string]string)
				combs[thing[0]][thing[1]] = v
			}
		}
	}
	fmt.Println("Downloaded & processed element combos")

	times := 0
	query := "INSERT INTO elem_combos VALUES "
	args := make([]interface{}, 0)
	for k, v := range combs {
		for key, val := range v {
			query += "(?,?,?),"
			args = append(args, k, key, val)
			handle(err)

			if (times % 200) == 0 {
				query = query[:len(query)-1]
				_, err = db.Exec(query, args...)
				handle(err)
				fmt.Println("Processed and wrote 200 records to SQL database")

				query = "INSERT INTO elem_combos VALUES "
				args = make([]interface{}, 0)
			}
			times++
		}
	}
	query = query[:len(query)-1]
	_, err = db.Exec(query, args...)
	handle(err)
	fmt.Println("Processed and wrote final records to SQL database")

	res, err = db.Query("SELECT * FROM suggestion_combos WHERE 1")
	handle(err)
	defer res.Close()
	var suggcombos map[string][]string
	var suggcombs map[string]map[string][]string = make(map[string]map[string][]string)
	for res.Next() {
		suggcombos = make(map[string][]string)
		err = res.Scan(&name, &combodat)
		handle(err)
		err = json.Unmarshal([]byte(combodat), &suggcombos)
		handle(err)
		for k, v := range suggcombos {
			thing := []string{name, k}
			sort.Strings(thing)
			_, exists := suggcombs[thing[0]]
			if exists {
				_, exists = suggcombs[thing[0]][thing[1]]
				if !exists {
					suggcombs[thing[0]][thing[1]] = v
				}
			} else {
				suggcombs[thing[0]] = make(map[string][]string)
				suggcombs[thing[0]][thing[1]] = v
			}
		}
	}
	fmt.Println("Downloaded & processed element combos")

	times = 0
	query = "INSERT INTO sugg_combos VALUES "
	args = make([]interface{}, 0)
	for k, v := range suggcombs {
		for key, val := range v {
			for _, val2 := range val {
				query += "(?,?,?),"
				args = append(args, k, key, val2)
				handle(err)

				if (times % 500) == 0 {
					query = query[:len(query)-1]
					_, err = db.Exec(query, args...)
					handle(err)
					fmt.Println("Processed and wrote 500 records to SQL database")

					query = "INSERT INTO sugg_combos VALUES "
					args = make([]interface{}, 0)
				}
				times++
			}
		}
	}
	query = query[:len(query)-1]
	_, err = db.Exec(query, args...)
	handle(err)
	fmt.Println("Processed and wrote final records to SQL database")
}
