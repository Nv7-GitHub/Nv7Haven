package main

import (
	"sort"
	"strings"
	"time"
)

func elems2txt(elems []string) string {
	sort.Strings(elems)
	return strings.Join(elems, "+")
}

func editDB() {
	_, err := db.Exec("DELETE FROM eod_elements WHERE 1")
	handle(err)
	end("Deleted existing data")
	times := 0

	args := make([]interface{}, 0)
	query := "INSERT INTO eod_elements VALUES "
	for _, gld := range glds {
		for _, elem := range gld.Elements {
			if elem.CreatedOn.Unix() < 0 {
				elem.CreatedOn = time.Now()
			}
			args = append(args, elem.Name, elem.Image, elem.Color, elem.Guild, elem.Comment, elem.Creator, int(elem.CreatedOn.Unix()), elems2txt(elem.Parents), elem.Complexity, elem.Difficulty, elem.UsedIn, elem.TreeSize)
			query += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),"

			if (times % 5000) == 0 {
				query = query[:len(query)-1]
				_, err = db.Exec(query, args...)
				handle(err)
				end("Processed and wrote 5000 records to DB")

				query = "INSERT INTO eod_elements VALUES "
				args = make([]interface{}, 0)
			}
			times++
		}
	}
	query = query[:len(query)-1]
	_, err = db.Exec(query, args...)
	handle(err)
	end("Processed and wrote final records to DB")
}
