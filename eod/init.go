package eod

import (
	"fmt"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/api"
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/basecmds"
	"github.com/Nv7-Github/Nv7Haven/eod/categories"
	"github.com/Nv7-Github/Nv7Haven/eod/elements"
	"github.com/Nv7-Github/Nv7Haven/eod/logs"
	"github.com/Nv7-Github/Nv7Haven/eod/polls"
	"github.com/Nv7-Github/Nv7Haven/eod/treecmds"
)

func (b *EoD) init() {
	// Initialize subsystems
	logs.InitEoDLogs()
	b.base = base.NewBase(b.Data, b.dg)
	b.basecmds = basecmds.NewBaseCmds(b.base, b.db, b.dg, b.Data)
	b.treecmds = treecmds.NewTreeCmds(b.Data, b.dg, b.base)
	b.polls = polls.NewPolls(b.Data, b.dg, b.base)
	b.categories = categories.NewCategories(b.Data, b.base, b.dg, b.polls)
	b.elements = elements.NewElements(b.Data, b.polls, b.db, b.base, b.dg)
	b.api = api.NewAPI(b.Data, b.base)

	// Run API
	b.api.Run()

	// Calc VCats
	start := time.Now()
	fmt.Println("Calculating VCats...")
	b.categories.CacheVCats()
	fmt.Println("Calculated in", time.Since(start))

	b.initHandlers()
	b.start()

	// Start stats saving
	go func() {
		b.basecmds.SaveStats()
		for {
			time.Sleep(time.Minute * 30)
			b.basecmds.SaveStats()
		}
	}()

	// Remove #0
	for _, db := range b.DB {
		for _, cat := range db.Cats() {
			_, exists := cat.Elements[0]
			if exists {
				delete(cat.Elements, 0)
				db.SaveCat(cat)
			}
		}
	}

	// Change elements created by devi
	/*dbv, _ := b.GetDB("705084182673621033")
	for _, el := range dbv.Elements {
		if el.Creator == "278931380191100929" {
			el.Creator = "812106732045205566"
			dbv.SaveElement(el)
			fmt.Println(el.Name)
		}
	}*/
}
