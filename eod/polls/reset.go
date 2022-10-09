package polls

import (
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

const MaxPollTime = 24 * time.Hour

func (b *Polls) CheckPollTime() {
	b.Data.RLock()
	for _, db := range b.Data.DB {
		b.Data.RUnlock()
		db.RLock()
		for _, poll := range db.Polls {
			if poll.CreatedOn == nil || time.Since(poll.CreatedOn.Time) > MaxPollTime { // Too long, delete
				db.RUnlock()
				poll.Upvotes = 0
				poll.Downvotes = 0
				err := db.DeletePoll(poll)
				if err == nil { // Successfully deleted, delete message
					b.dg.ChannelMessageDelete(poll.Channel, poll.Message)
				}
				db.RLock()
			}
		}
		db.RUnlock()
		b.Data.RLock()
	}
	b.Data.RUnlock()
}

func (b *Polls) ResetPolls(m types.Msg, rsp types.Rsp) {
	rsp.Acknowledge()
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	db.RLock()
	for _, poll := range db.Polls {
		db.RUnlock()
		poll.Upvotes = 0
		poll.Downvotes = 0

		// Delete
		err := db.DeletePoll(poll)
		if rsp.Error(err) {
			return
		}

		// Delete msg
		b.dg.ChannelMessageDelete(poll.Channel, poll.Message)

		// Create new one
		err = b.CreatePoll(poll)
		if err != nil {
			rsp.Error(err)
			return
		}
		db.RLock()
	}
	db.RUnlock()

	rsp.Message("Done!")
}
