package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql
)

const (
	dbUser = "root"
	dbName = "nv7haven"
	guild  = "705084182673621033"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

var tm = time.Now()

func end(print string) {
	newtime := time.Now()
	fmt.Println(print, "in", newtime.Sub(tm))
	tm = newtime
}

type Element struct {
	Name       string   `json:"name"`
	Image      string   `json:"image"`
	Color      int      `json:"color"`
	Guild      string   `json:"guild"`
	Comment    string   `json:"comment"`
	Creator    string   `json:"creator"`
	Createdon  int      `json:"createdon"`
	Parents    []string `json:"parents"`
	Complexity int      `json:"complexity"`
	Difficulty int      `json:"difficulty"`
	UsedIn     int      `json:"usedin"`
	TreeSize   int      `json:"treesize"`
}

type Combo struct {
	Guild    string   `json:"guild"`
	Elements []string `json:"elements"`
	Result   string   `json:"result"`
}

type Output struct {
	Elements []Element `json:"elements"`
	Combos   []Combo   `json:"combos"`
}

func main() {
	db, err := sql.Open("mysql", dbUser+":"+os.Getenv("PASSWORD")+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	handle(err)
	defer db.Close()
	db.SetMaxOpenConns(10)
	end("Connected to DB")

	// Download elements
	elems := make([]Element, 0)

	res, err := db.Query("SELECT * FROM eod_elements WHERE guild=?", guild)
	handle(err)
	defer res.Close()
	var elem Element
	var parents string
	for res.Next() {
		err = res.Scan(&elem.Name, &elem.Image, &elem.Color, &elem.Guild, &elem.Comment, &elem.Creator, &elem.Createdon, &parents, &elem.Complexity, &elem.Difficulty, &elem.UsedIn, &elem.TreeSize)
		handle(err)

		if len(parents) > 0 {
			elem.Parents = strings.Split(parents, "+")
		}

		elems = append(elems, elem)
	}
	end("Downloaded elements")

	// Download combos
	combos := make([]Combo, 0)

	combs, err := db.Query("SELECT * FROM eod_combos WHERE guild=?", guild)
	handle(err)
	defer combs.Close()

	var els string
	var comb Combo
	for combs.Next() {
		err = combs.Scan(&comb.Guild, &els, &comb.Result)
		handle(err)

		comb.Elements = strings.Split(els, "+")
		combos = append(combos, comb)
	}
	end("Downloaded combos")

	// Save
	out := Output{
		Elements: elems,
		Combos:   combos,
	}
	outFile, err := os.Create("out.json")
	handle(err)

	enc := json.NewEncoder(outFile)
	enc.Encode(out)
	end("Saved data")
}
