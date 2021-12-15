package eod

import (
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
	"github.com/schollz/progressbar/v3"
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
	admin.InitAdmin(b.Data)

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

	// Recalc autocats?
	if types.RecalcAutocats {
		for _, db := range b.DB {
			for _, elem := range db.Elements {
				b.polls.Autocategorize(elem.Name, db.Guild)
			}
		}
	}
}
