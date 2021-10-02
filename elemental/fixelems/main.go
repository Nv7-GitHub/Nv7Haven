package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql" // mysql
	"github.com/sasha-s/go-deadlock"
)

// Element has the data for a created element
type Element struct {
	Color      string   `json:"color"`
	Comment    string   `json:"comment"`
	CreatedOn  int      `json:"createdOn"`
	Creator    string   `json:"creator"`
	Name       string   `json:"name"`
	Parents    []string `json:"parents"`
	Pioneer    string   `json:"pioneer"`
	Uses       int      `json:"uses"`
	FoundBy    int      `json:"foundby"`
	Complexity int      `json:"complexity"`
}

const (
	dbUser = "root"
	dbName = "nv7haven"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

var lock = &deadlock.RWMutex{}
var wg = &sync.WaitGroup{}

var complcache = make(map[string]int)

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

	res, err := db.Query("SELECT * FROM elements WHERE 1")
	handle(err)
	elems := make(map[string]Element)
	defer res.Close()
	for res.Next() {
		var elem Element
		elem.Parents = make([]string, 2)
		err = res.Scan(&elem.Name, &elem.Color, &elem.Comment, &elem.Parents[0], &elem.Parents[1], &elem.Creator, &elem.Pioneer, &elem.CreatedOn, &elem.Complexity, &elem.Uses, &elem.FoundBy)
		handle(err)
		if (elem.Parents[0] == "") && (elem.Parents[1] == "") {
			elem.Parents = make([]string, 0)
		}

		wg.Add(1)
		go func(elem Element) {
			uses := db.QueryRow("SELECT COUNT(1) FROM elem_combos WHERE elem1=? OR elem2=?", elem.Name, elem.Name)
			err = uses.Scan(&elem.Uses)
			handle(err)

			foundby := db.QueryRow("SELECT COUNT(1) FROM users WHERE found LIKE ?", `%`+elem.Name+`%`)
			err = foundby.Scan(&elem.FoundBy)
			handle(err)

			lock.Lock()
			elems[elem.Name] = elem
			lock.Unlock()
			fmt.Println(elem.Name, elem.FoundBy, elem.Uses)

			wg.Done()
		}(elem)
	}
	wg.Wait()
	for k, v := range elems {
		v.Complexity = calcComplexity(v, elems)
		elems[k] = v
		wg.Add(1)
		go func(v Element) {
			_, err = db.Exec("UPDATE elements SET complexity=?, foundby=?, uses=? WHERE name=?", v.Complexity, v.FoundBy, v.Uses, v.Name)
			handle(err)
			wg.Done()
			fmt.Println(v.Name, v.Complexity, v.FoundBy, v.Uses, v.Name)
		}(v)
	}
	wg.Wait()
}

func calcComplexity(elem Element, elems map[string]Element) int {
	scr, exists := complcache[elem.Name]
	if exists {
		return scr
	}
	if len(elem.Parents) == 0 {
		return 0
	}
	parent1 := elems[elem.Parents[0]]
	parent2 := elems[elem.Parents[1]]
	comp1 := calcComplexity(parent1, elems)
	comp2 := calcComplexity(parent2, elems)

	if comp1 > comp2 {
		scr = comp1 + 1
	} else {
		scr = comp2 + 1
	}
	complcache[elem.Name] = scr
	return scr
}
