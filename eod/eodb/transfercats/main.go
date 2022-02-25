package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	home, err := os.UserHomeDir()
	handle(err)
	dbPath := filepath.Join(home, "go/src/github.com/Nv7-Github/Nv7haven/data/eod")
	fmt.Println("Loading...")
	start := time.Now()
	db, err := eodb.NewData(dbPath)
	handle(err)
	fmt.Println("Loaded in", time.Since(start))

	// Get polls
	for _, db := range db.DB {
		// Cache VCats
		regs := make(map[string]*regexp.Regexp)
		vcats := make(map[string]map[int]types.Empty)
		for _, vcat := range db.VCats() {
			if vcat.Rule == types.VirtualCategoryRuleRegex {
				vcat.Cache = make(map[int]types.Empty)
				vcats[vcat.Name] = vcat.Cache
				regs[vcat.Name] = regexp.MustCompile(vcat.Data["regex"].(string))
			}
		}
		for _, el := range db.Elements {
			for k, reg := range regs {
				match := reg.Match([]byte(el.Name))
				if match {
					vcats[k][el.ID] = types.Empty{}
				}
			}
		}

		// Cache cats
		for _, cat := range db.Cats() {
			err := db.SaveCatCache(cat.Name, cat.Elements)
			handle(err)
		}

		// Save VCat caches
		for _, vcat := range db.VCats() {
			if vcat.Rule == types.VirtualCategoryRuleRegex {
				err := db.SaveCatCache(vcat.Name, vcat.Cache)
				handle(err)
			}
		}

		db.Close()
	}
}
