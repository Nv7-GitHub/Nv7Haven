package polls

import (
	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

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

		// Cleate new one
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
