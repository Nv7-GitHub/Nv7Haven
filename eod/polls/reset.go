package polls

import (
	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *Polls) ResetPolls(m types.Msg, rsp types.Rsp) {
	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		return
	}
	rsp.Acknowledge()

	for id, poll := range dat.Polls {
		poll.Upvotes = 0
		poll.Downvotes = 0

		// Delete
		_, err := b.db.Exec("DELETE FROM eod_polls WHERE guild=? AND channel=? AND message=?", poll.Guild, poll.Channel, poll.Message)
		if err != nil {
			rsp.Error(err)
			return
		}

		// Delete from cache
		dat.Lock.Lock()
		delete(dat.Polls, id)
		dat.Lock.Unlock()

		// Delete msg
		b.dg.ChannelMessageDelete(poll.Channel, poll.Message)
		err = b.CreatePoll(poll)
		if err != nil {
			rsp.Error(err)
			return
		}
	}

	rsp.Message("Done!")
}
