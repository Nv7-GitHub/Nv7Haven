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
	"github.com/Nv7-Github/Nv7Haven/eod/types"
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

	// Check polls
	start = time.Now()
	fmt.Println("Checking polls...")
	for _, db := range b.Data.DB {
		for _, poll := range db.Polls {
			msg, err := b.dg.ChannelMessage(poll.Channel, poll.Message)
			if err != nil || msg == nil {
				db.DeletePoll(poll)
				continue
			}
			for _, r := range msg.Reactions {
				if r.Emoji.Name == types.UpArrow {
					poll.Upvotes = r.Count - 1
				} else if r.Emoji.Name == types.DownArrow {
					poll.Downvotes = r.Count - 1
				}
			}

			// Check if being deleted by author
			reactor := ""
			downvote := false
			if poll.Downvotes > 0 {
				r, err := b.dg.MessageReactions(poll.Channel, poll.Message, types.DownArrow, 100, "", "")
				if err == nil {
					for _, u := range r {
						if u.ID == poll.Suggestor {
							downvote = true
							reactor = u.Username
							break
						}
					}
				}
			}

			// Handle poll votes
			b.polls.CheckReactions(db, poll, reactor, downvote)
		}
	}
	fmt.Println("Checked in", time.Since(start))

	b.initHandlers()

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
