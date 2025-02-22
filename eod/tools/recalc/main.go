// Run like "go run . > output.txt"

package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"slices"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
) // postgres

const (
	dbUser = "postgres"
	dbName = "nv7haven"
	dbPort = "5432"
)

var lastTime = time.Now()

func TimingPrint(message string) {
	fmt.Printf("%s in %v\n", message, time.Since(lastTime))
	lastTime = time.Now()
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("===PATH RECALC===")
	db, err := sqlx.Connect("postgres", "user="+dbUser+" dbname="+dbName+" sslmode=disable port="+dbPort+" host="+os.Getenv("MYSQL_HOST")+" password="+os.Getenv("PASSWORD"))
	if err != nil {
		panic(err)
	}
	TimingPrint("Connected to DB")
	db.Exec(`DROP TABLE elements_update`)
	// Get guilds
	var guilds []string
	err = db.Select(&guilds, "SELECT DISTINCT(guild) FROM elements")
	handle(err)
	TimingPrint("Fetched guilds")
	// Create update table
	_, err = db.Exec(`CREATE TABLE elements_update (
		id integer,
		guild text,
		parents integer[],
		treesize integer,
		usedin integer,
		madewith integer,
foundby integer,
tier integer
 	)`)
	handle(err)

	// Recalc
	for _, guild := range guilds {
		RecalcGuild(guild, db)
	}

	// Update
	_, err = db.Exec(`UPDATE elements SET parents = elements_update.parents, treesize = elements_update.treesize FROM elements_update WHERE elements.id = elements_update.id AND elements.guild = elements_update.guild`)
	handle(err)
	_, err = db.Exec(`UPDATE elements SET usedin = elements_update.usedin FROM elements_update WHERE elements.id = elements_update.id`)
	handle(err)
	_, err = db.Exec(`UPDATE elements SET madewith = elements_update.madewith FROM elements_update WHERE elements.id = elements_update.id`)
	TimingPrint("Updated element")

	// Drop update table
	_, err = db.Exec(`DROP TABLE elements_update`)
	handle(err)
	TimingPrint("Dropped update table")
}

type Combo struct {
	Els    pq.Int32Array `db:"els"`
	Result int32         `db:"result"`
	Done   bool
}

type Element struct {
	Guild    string        `db:"guild"`
	ID       int32         `db:"id"`
	TreeSize int           `db:"treesize"`
	Parents  pq.Int32Array `db:"parents"`
	MadeWith int           `db:"madewith"`
	FoundBy  int           `db:"foundby"`
	UsedIn   int           `db:"usedin"`
	Tier     int           `db:"tier"`
}

func AverageTreeSize(els []Element) {
	var total int
	for _, el := range els {
		total += el.TreeSize
	}
	fmt.Printf("Average tree size: %d\n", total/len(els))
	lastTime = time.Now()
}

func CalcTreeSize(el int32, els []Element, done map[int32]struct{}) {
	_, exists := done[el]
	if exists {
		return
	}
	elem := els[el-1]
	for _, par := range elem.Parents {
		CalcTreeSize(par, els, done)
	}
	done[el] = struct{}{}
}

func RecalcGuild(guild string, db *sqlx.DB) {
	fmt.Printf("===GUILD RECALC: %s===\n", guild)
	done := make(map[int32]struct{})
	for i := int32(1); i <= 4; i++ { // Starters done
		done[i] = struct{}{}
	}

	// Fetch combos
	var combos []Combo
	err := db.Select(&combos, "SELECT els, result FROM combos WHERE guild=$1", guild)
	handle(err)
	TimingPrint("Fetched combos")

	// Fetch elements
	var elements []Element
	err = db.Select(&elements, "SELECT guild, id, treesize, parents, madewith, foundby, usedin,tier FROM elements WHERE guild=$1 ORDER BY id", guild)
	handle(err)
	TimingPrint("Fetched elements")

	mademap := make(map[int32]int32, len(elements))
	usedmap := make(map[int32]int32, len(elements))
	// Calc some stats
	AverageTreeSize(elements)

	// Recalc

	changed := -1
	for changed != 0 {
		changed = 0

		var elemdone []int32
		// Loop through combos
		for i, comb := range combos {
			// Check if done
			if comb.Done {
				continue
			}
			_, exists := done[comb.Result]
			if exists {
				mademap[comb.Result-1]++
				var combeldone []int32
				for _, combelem := range comb.Els {
					if !slices.Contains(combeldone, combelem) {
						usedmap[combelem-1]++
						combeldone = append(combeldone, combelem)
					}
				}
				combos[i].Done = true
				continue
			}

			// If not done, check if can be done
			valid := true
			for _, elem := range comb.Els {
				_, exists := done[elem] // Check if element has been done, if it hasnt then not ready
				if !exists {
					valid = false
					break
				}
			}
			if !valid {
				continue
			}

			// Update
			el := elements[comb.Result-1]

			if el.ID != comb.Result {
				fmt.Fprintf(os.Stderr, "FAILED UPDATE: %s\n", guild)
				return
			}
			el.Parents = comb.Els
			elements[comb.Result-1] = el

			// Add to donelist
			elemdone = append(elemdone, comb.Result)

		}
		// Update done
		for i := 0; i < len(elemdone); i++ {
			done[elemdone[i]] = struct{}{}
			changed++
		}
	}
	TimingPrint("Recalculated combos")

	// Recalc tree size
	for i, el := range elements {
		elements[i].UsedIn = int(usedmap[int32(i)])
		elements[i].MadeWith = int(mademap[int32(i)])
		if i < 4 {
			continue
		}
		done := make(map[int32]struct{}, el.TreeSize)
		CalcTreeSize(el.ID, elements, done)
		elements[i].TreeSize = len(done)

	}
	TimingPrint("Recalculated tree size")

	// Stats
	AverageTreeSize(elements)

	// Add elements to update table
	BulkInsert(`INSERT INTO elements_update (id, guild, parents, treesize, usedin, madewith,tier) VALUES (:id, :guild, :parents, :treesize, :usedin, :madewith,:tier)`, elements, db)
	handle(err)
	TimingPrint("Added elements to update table")
}

// From https://github.com/jmoiron/sqlx/issues/552
func BulkInsert(insertQuery string, myStructs []Element, db *sqlx.DB) {
	tx, err := db.Beginx()
	if err != nil {
		log.Printf("Couldn't begin transaction %+v", err)
		return
	}

	// The number of placeholders allowed in a query is capped at 2^16, therefore,
	// divide 2^16 by the number of fields in the struct, and that is the max
	// number of bulk inserts possible. Use that number to chunk the inserts.
	v := reflect.ValueOf(myStructs[0])
	maxBulkInsert := ((1 << 16) / v.NumField()) - 1

	// send batch requests
	for i := 0; i < len(myStructs); i += maxBulkInsert {
		batch := myStructs[i:Min(i+maxBulkInsert, len(myStructs))]
		_, err := tx.NamedExec(insertQuery, batch)
		if err != nil {
			e := tx.Rollback()
			if e != nil {
				log.Printf("Couldn't rollback %+v", e)
			}
			log.Printf("Couldn't insert batch %+v", err)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		e := tx.Rollback()
		if e != nil {
			log.Printf("Couldn't rollback %+v", err)
		}
		log.Printf("Couldn't commit %+v", err)
		return
	}
}
func Min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
