package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

type Element struct {
	ID        int       `db:"id"`
	Guild     string    `db:"guild"`
	Name      string    `db:"name"`
	Image     string    `db:"image"`
	Color     int       `db:"color"`
	Comment   string    `db:"comment"`
	Creator   string    `db:"creator"`
	CreatedOn time.Time `db:"createdon"`
	Commenter string    `db:"commenter"`
	Colorer   string    `db:"colorer"`
	Imager    string    `db:"imager"`
}

type Combo struct {
	Guild    string      `db:"guild"`
	Elements interface{} `db:"els"` // pq array
	Result   int         `db:"result"`
}

type Inventory struct {
	Guild string      `db:"guild"`
	User  string      `db:"user"`
	Inv   interface{} `db:"inv"` // pq array
}

type Category struct {
	Guild   string `db:"guild"`
	Name    string `db:"name"`
	Comment string `db:"comment"`
	Image   string `db:"image"`
	Color   int    `db:"color"`

	Commenter string `db:"commenter"`
	Imager    string `db:"imager"`
	Colorer   string `db:"colorer"`

	Elements interface{} `db:"elements"` // pq array
}

// TODO: VCats (integrate into eod_categories?), polls

func main() {
	// Eodb
	home, err := os.UserHomeDir()
	handle(err)
	dbPath := filepath.Join(home, "go/src/github.com/Nv7-Github/Nv7haven/data/eod")
	fmt.Println("Loading...")
	start := time.Now()
	eodb, err := eodb.NewData(dbPath)
	handle(err)
	fmt.Println("Loaded in", time.Since(start))

	// DB
	start = time.Now()
	db, err := sqlx.Connect("postgres", "user=postgres dbname=nv7haven sslmode=disable port = 5432 host="+os.Getenv("MYSQL_HOST")+" password="+os.Getenv("PASSWORD"))
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected in", time.Since(start))

	// Add elements
	/*start = time.Now()
	els := make([]Element, 0)
	for _, db := range eodb.DB {
		for _, el := range db.Elements {
			var timeV time.Time
			if el.CreatedOn == nil {
				timeV = time.Now()
			} else {
				timeV = el.CreatedOn.Time
			}
			if strings.Contains(el.Comment, string(rune(0))) {
				el.Comment = strings.ReplaceAll(el.Comment, string(rune(0)), "")
			}
			els = append(els, Element{
				ID:        el.ID,
				Guild:     db.Guild,
				Name:      el.Name,
				Image:     el.Image,
				Color:     el.Color,
				Comment:   el.Comment,
				Creator:   el.Creator,
				CreatedOn: timeV,
				Commenter: el.Commenter,
				Colorer:   el.Colorer,
				Imager:    el.Imager,
			})
		}
	}
	fmt.Println("Got elements in", time.Since(start))

	BulkInsert("INSERT INTO eod_elements (id, guild, name, image, color, comment, creator, createdon, commenter, colorer, imager) VALUES (:id, :guild, :name, :image, :color, :comment, :creator, :createdon, :commenter, :colorer, :imager)", els, db)*/

	// Add combos
	/*start = time.Now()
	combs := make([]Combo, 0)
	for _, db := range eodb.DB {
	skip:
		for k, com := range db.Combos() {
			rawV := strings.Split(k, "+")
			els := make([]int, len(rawV))
			for i, el := range rawV {
				els[i], err = strconv.Atoi(el)
				if err != nil {
					continue skip
				}
			}
			sort.Ints(els)
			combs = append(combs, Combo{
				Guild:    db.Guild,
				Elements: pq.Array(els),
				Result:   com,
			})
		}
	}
	fmt.Println("Got combos in", time.Since(start))

	BulkInsert("INSERT INTO eod_combos (guild, result, els) VALUES (:guild, :result, :els)", combs, db)*/

	// Add invs
	/*start = time.Now()
	invs := make([]Inventory, 0)
	for _, db := range eodb.DB {
		for user, inv := range db.Invs() {
			items := make([]int, 0, len(inv.Elements))
			for el := range inv.Elements {
				items = append(items, el)
			}
			invs = append(invs, Inventory{
				Guild: db.Guild,
				User:  user,
				Inv:   pq.Array(items),
			})
		}
	}
	fmt.Println("Got invs in", time.Since(start))

	BulkInsert("INSERT INTO eod_invs (guild, \"user\", inv) VALUES (:guild, :user, :inv)", invs, db)*/

	// Add categories
	/*start = time.Now()
	cats := make([]Category, 0)
	for _, db := range eodb.DB {
		for _, cat := range db.Cats() {
			items := make([]int, 0, len(cat.Elements))
			for el := range cat.Elements {
				items = append(items, el)
			}
			cats = append(cats, Category{
				Guild:     db.Guild,
				Name:      cat.Name,
				Image:     cat.Image,
				Color:     cat.Color,
				Comment:   cat.Comment,
				Imager:    cat.Imager,
				Colorer:   cat.Colorer,
				Commenter: cat.Commenter,
				Elements:  pq.Array(items),
			})
		}
	}
	fmt.Println("Got cats in", time.Since(start))

	BulkInsert("INSERT INTO eod_categories (guild, name, image, color, comment, imager, colorer, commenter, elements) VALUES (:guild, :name, :image, :color, :comment, :imager, :colorer, :commenter, :elements)", cats, db)*/
}

func BulkInsert[T any](insertQuery string, myStructs []T, db *sqlx.DB) {
	tx, err := db.Beginx()
	handle(err)

	v := reflect.ValueOf(myStructs[0])
	maxBulkInsert := ((1 << 16) / v.NumField()) - 1

	// send batch requests
	for i := 0; i < len(myStructs); i += maxBulkInsert {
		start := time.Now()
		batch := myStructs[i:Min(i+maxBulkInsert, len(myStructs))]
		_, err := tx.NamedExec(insertQuery, batch)
		if err != nil {
			log.Print(err)
			e := tx.Rollback()
			handle(e)
		}
		fmt.Println("Put in batch in", time.Since(start))
	}

	start := time.Now()
	err = tx.Commit()
	if err != nil {
		e := tx.Rollback()
		if e != nil {
			log.Printf("Couldn't rollback %+v", err)
		}
		log.Printf("Couldn't commit %+v", err)
	}
	fmt.Println("Committed in", time.Since(start))
}

// Min takes 2 ints and returns the lesser of them.
func Min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
