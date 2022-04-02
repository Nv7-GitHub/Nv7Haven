package polls

import (
	"errors"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *Polls) CreatePoll(p types.Poll) error {
	db, res := b.GetDB(p.Guild)
	if !res.Exists {
		return nil
	}
	if db.Config.VoteCount == 0 {
		b.handlePollSuccess(p)
		return nil
	}
	msg := ""
	// check poll limit
	if db.Config.PollCount > 0 {
		uPolls := 0
		db.RLock()
		for _, val := range db.Polls {
			if val.Suggestor == p.Suggestor {
				uPolls++
			}
		}
		db.RUnlock()
		msg = "Too many active polls!"
		if uPolls >= db.Config.PollCount {
			return errors.New(msg)
		}
	}
	// Get embed
	emb, err := b.GetPollEmbed(db, p)
	if err != nil {
		return err
	}
	m, err := b.dg.ChannelMessageSendEmbed(db.Config.VotingChannel, emb)
	if err != nil {
		return err
	}
	p.Message = m.ID

	// Add reactions
	err = b.dg.MessageReactionAdd(p.Channel, p.Message, types.UpArrow)
	if err != nil {
		return err
	}
	err = b.dg.MessageReactionAdd(p.Channel, p.Message, types.DownArrow)
	if err != nil {
		return err
	}
	p.CreatedOn = types.NewTimeStamp(time.Now())
	err = db.NewPoll(p)
	if err != nil {
		return err
	}
	return err
}
