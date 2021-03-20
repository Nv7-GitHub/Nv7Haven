package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql
)

type empty struct{}

type element struct {
	Name       string
	Categories map[string]empty
	Image      string
	Guild      string
	Comment    string
	Creator    string
	CreatedOn  time.Time
	Parents    []string
	Complexity int
	Difficulty int
}

const (
	dbUser     = "u57_fypTHIW9t8"
	dbPassword = "C7HgI6!GF0NaHCrdUi^tEMGy"
	dbName     = "s57_nv7haven"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

var wg = &sync.WaitGroup{}

var complcache = make(map[string]int)

// Fixelems fixes the elements
func main() {
	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	handle(err)
	defer db.Close()
	db.SetMaxOpenConns(10)

	fmt.Println("Connected")

	res, err := db.Query("SELECT name, parent1, parent2, complexity FROM eod_elements WHERE 1")
	handle(err)
	elems := make(map[string]element)
	defer res.Close()
	for res.Next() {
		var elem element
		elem.Parents = make([]string, 2)
		err = res.Scan(&elem.Name, &elem.Parents[0], &elem.Parents[1], &elem.Complexity)
		handle(err)
		if (elem.Parents[0] == "") && (elem.Parents[1] == "") {
			elem.Parents = make([]string, 0)
		}
		elems[strings.ToLower(elem.Name)] = elem
	}
	for k, v := range elems {
		v.Difficulty = calcComplexity(v, elems)
		elems[k] = v
		wg.Add(1)
		go func(v element) {
			_, err = db.Exec("UPDATE eod_elements SET difficulty=? WHERE name=?", v.Difficulty, v.Name)
			handle(err)
			wg.Done()
			fmt.Println(v.Name, v.Complexity, v.Difficulty)
		}(v)
	}
	wg.Wait()
}

func calcComplexity(elem element, elems map[string]element) int {
	scr, exists := complcache[elem.Name]
	if exists {
		return scr
	}
	if len(elem.Parents) == 0 {
		return 0
	}
	parent1 := elems[strings.ToLower(elem.Parents[0])]
	parent2 := elems[strings.ToLower(elem.Parents[1])]
	comp1 := calcComplexity(parent1, elems)
	comp2 := calcComplexity(parent2, elems)

	if comp1 > comp2 {
		scr = comp1
	} else {
		scr = comp2
	}
	if strings.EqualFold(elem.Parents[0], elem.Parents[1]) {
		scr++
	}
	complcache[elem.Name] = scr
	return scr
}
