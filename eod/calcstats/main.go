package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql
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

type elem struct {
	name      string
	createdon time.Time
}

type row struct {
	time  time.Time
	count int
}

func main() {
	dbPassword := os.Getenv("PASSWORD")

	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	handle(err)
	defer db.Close()

	fmt.Println("Connected")

	// Get elements
	res, err := db.Query("SELECT name, (CASE WHEN createdon=1637536881 THEN 1605988759 ELSE createdon END) FROM `eod_elements` ORDER BY (CASE WHEN createdon=1637536881 THEN 1605988759 ELSE createdon END) ASC")
	handle(err)
	defer res.Close()

	elems := make([]elem, 0)
	start := time.Now()

	var name string
	var createdon int64
	for res.Next() {
		err = res.Scan(&name, &createdon)
		handle(err)

		elems = append(elems, elem{
			name:      name,
			createdon: time.Unix(createdon, 0),
		})
		if len(elems)%10000 == 0 {
			fmt.Println("Downloaded 10000 elements in", time.Since(start))
			start = time.Now()
		}
	}

	// Calculate stats for every 30 mins
	currTime := elems[0].createdon
	count := 0

	out := []row{
		{
			time:  currTime,
			count: count,
		},
	}

	for _, elem := range elems {
		count++
		if elem.createdon.Sub(currTime) > (time.Hour * 24) {
			currTime = currTime.Add(time.Hour * 24)
			out = append(out, row{
				time:  currTime,
				count: count,
			})
		}
	}

	// Save stats to DB
	times := 0
	query := "INSERT INTO eod_stats (time, elemcnt) VALUES "
	args := make([]interface{}, 0)

	start = time.Now()

	for _, row := range out {
		query += "(?,?),"
		args = append(args, row.time.Unix(), row.count)

		if (times % 10000) == 0 {
			fmt.Println("Wrote 10000 records in", time.Since(start))
			start = time.Now()

			query = query[:len(query)-1]
			_, err := db.Exec(query, args...)
			handle(err)
			query = "INSERT INTO eod_stats (time, elemcnt) VALUES "
			args = make([]interface{}, 0)
		}

		times++
	}
	start = time.Now()
	query = query[:len(query)-1]
	_, err = db.Exec(query, args...)
	handle(err)
	fmt.Println("Wrote", len(args)/2, " records in", time.Since(start))
}
