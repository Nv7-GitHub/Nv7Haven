package main

//https://www.ssa.gov/oact/babynames/limits.html

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

const url = "https://www.ssa.gov/oact/babynames/names.zip"

var reg = regexp.MustCompile("yob([0-9][0-9][0-9][0-9]).txt")

const (
	dbUser     = "u51_iYXt7TBZ0e"
	dbPassword = "W!QnD2u896yo.J4fww9X.h+J"
	dbName     = "s51_nv7haven"
)

var tm time.Time = time.Now()

func endTimer(print string) {
	newtime := time.Now()
	fmt.Println(print, "in", newtime.Sub(tm))
	tm = newtime
}

func main() {
	// SQL
	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp(c.filipk.in:3306)/"+dbName)
	handle(err)
	defer db.Close()
	endTimer("Connected to SQL database")

	// Download
	resp, err := http.Get(url)
	handle(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handle(err)
	endTimer("Downloaded baby name statistics ZIP")

	// Read ZIP
	zipfile, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	handle(err)
	endTimer("Converted baby name statistics ZIP")

	// Process files
	files := make([]file, 0)
	for _, fl := range zipfile.File {
		match := reg.FindAllStringSubmatch(fl.Name, -1)
		if len(match) < 1 || len(match[0]) < 2 {
			continue
		}
		yr, err := strconv.Atoi(match[0][1])
		handle(err)
		files = append(files, file{
			year: yr,
			file: fl,
		})
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].year > files[j].year
	})
	endTimer("Processed file names")

	// Convert file to CSV
	fl, err := files[0].file.Open()
	handle(err)
	defer fl.Close()
	csv := csv.NewReader(fl)
	endTimer("Read ZIP")

	// Read CSV
	times := 0
	query := "INSERT INTO names VALUES "
	args := make([]interface{}, 0)
	for true {
		vals, err := csv.Read()
		if vals == nil {
			break
		}
		handle(err)

		// Write vals
		count, err := strconv.Atoi(vals[2])
		handle(err)
		query += "(?,?,?),"
		args = append(args, vals[0], vals[1] == "M", count)
		handle(err)

		if (times % 10000) == 0 {
			query = query[:len(query)-1]
			_, err = db.Exec(query, args...)
			handle(err)
			endTimer("Processed and wrote 10000 records to SQL database")

			query = "INSERT INTO names VALUES "
			args = make([]interface{}, 0)
		}
		times++
	}
	query = query[:len(query)-1]
	_, err = db.Exec(query, args...)
	handle(err)
	endTimer("Processed and wrote final records to SQL database")
}

type file struct {
	file *zip.File
	year int
}
