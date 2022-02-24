package eod

import (
	"fmt"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/admin"
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/basecmds"
	"github.com/Nv7-Github/Nv7Haven/eod/categories"
	"github.com/Nv7-Github/Nv7Haven/eod/elements"
	"github.com/Nv7-Github/Nv7Haven/eod/logs"
	"github.com/Nv7-Github/Nv7Haven/eod/polls"
	"github.com/Nv7-Github/Nv7Haven/eod/treecmds"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/gofiber/fiber/v2"
	"github.com/schollz/progressbar/v3"
)

func (b *EoD) init(app *fiber.App) {
	// Initialize subsystems
	logs.InitEoDLogs()
	b.base = base.NewBase(b.Data, b.dg)
	b.basecmds = basecmds.NewBaseCmds(b.base, b.db, b.dg, b.Data)
	b.treecmds = treecmds.NewTreeCmds(b.Data, b.dg, b.base)
	b.polls = polls.NewPolls(b.Data, b.dg, b.base)
	b.categories = categories.NewCategories(b.Data, b.base, b.dg, b.polls)
	b.elements = elements.NewElements(b.Data, b.polls, b.db, b.base, b.dg)
	admin.InitAdmin(b.Data, app)

	// Polls
	cnt := 0
	for _, db := range b.DB {
		cnt += len(db.Polls)
	}
	bar := progressbar.New(cnt)

	for _, db := range b.DB {
		for _, po := range db.Polls {
			msg, err := b.dg.ChannelMessage(po.Channel, po.Message)
			if err != nil {
				err := db.DeletePoll(po)
				if err != nil {
					panic(err)
				}
				continue
			}
			for _, r := range msg.Reactions {
				if r.Emoji.Name == types.UpArrow {
					po.Upvotes = r.Count - 1
				}

				if r.Emoji.Name == types.DownArrow {
					po.Downvotes = r.Count - 1
				}
			}

			// Get downs to see who last reacted
			downs, err := b.dg.MessageReactions(po.Channel, po.Message, types.DownArrow, 100, "", "")
			if err != nil {
				err := db.DeletePoll(po)
				if err != nil {
					panic(err)
				}
				continue
			}

			lastDown := downs[len(downs)-1].ID
			b.polls.CheckReactions(db, po, lastDown, false)

			db.SavePoll(po)
			bar.Add(1)
		}
	}
	bar.Finish()

	// Cache vcats
	fmt.Println("Caching Virtual Categories...")
	start := time.Now()
	b.categories.CacheVCats()
	fmt.Println("Cached Virtual Categories in", time.Since(start))

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

	// Change starters
	/*for _, db := range b.DB {
		for _, el := range base.StarterElements {
			e, res := db.GetElement(el.ID)
			if !res.Exists {
				continue
			}
			e.Creator = el.Creator
			err := db.SaveElement(e)
			if err != nil {
				fmt.Println(err)
			}
		}
	}*/
}
