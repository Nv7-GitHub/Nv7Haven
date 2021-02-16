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

const yearsBack = 78

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
	_, err = db.Exec("DELETE FROM names WHERE 1")
	handle(err)
	endTimer("Connected to SQL database and cleared data")

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
	names := make(map[string]dat)
	for i := 0; i <= yearsBack; i++ {
		fl, err := files[0].file.Open()
		handle(err)
		defer fl.Close()
		csv := csv.NewReader(fl)
		for true {
			vals, err := csv.Read()
			if vals == nil {
				break
			}
			handle(err)

			count, err := strconv.Atoi(vals[2])
			handle(err)
			_, exists := names[vals[0]]
			isMale := vals[1] == "M"
			if exists {
				val := names[vals[0]]
				if val.isMale == isMale {
					val.count += count
				} else {
					if val.count < count {
						val.isMale = isMale
						val.count = count
					}
				}
				names[vals[0]] = val
			} else {
				names[vals[0]] = dat{
					count:  count,
					isMale: isMale,
				}
			}
		}
	}
	endTimer("Read and processed all data")

	// Read CSV
	times := 0
	query := "INSERT INTO names VALUES "
	args := make([]interface{}, 0)
	for k, v := range names {
		query += "(?,?,?),"
		args = append(args, k, v.isMale, v.count)
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

type dat struct {
	count  int
	isMale bool
}
