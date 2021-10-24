package polls

import (
	"encoding/json"
	"errors"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *Polls) CreatePoll(p types.Poll) error {
	b.lock.RLock()
	dat, exists := b.dat[p.Guild]
	b.lock.RUnlock()
	if !exists {
		return nil
	}
	if dat.VoteCount == 0 {
		b.handlePollSuccess(p)
		return nil
	}
	msg := ""
	if dat.PollCount > 0 {
		uPolls := 0
		for _, val := range dat.Polls {
			if val.Value4 == p.Value4 {
				uPolls++
			}
		}
		msg = "Too many active polls!"
		if uPolls >= dat.PollCount {
			return errors.New(msg)
		}
	}
	emb, err := b.GetPollEmbed(dat, p)
	if err != nil {
		return err
	}
	m, err := b.dg.ChannelMessageSendEmbed(dat.VotingChannel, emb)
	if err != nil {
		return err
	}
	p.Message = m.ID

	if !base.IsFoolsMode {
		err := b.dg.MessageReactionAdd(p.Channel, p.Message, types.UpArrow)
		if err != nil {
			return err
		}
	}
	err = b.dg.MessageReactionAdd(p.Channel, p.Message, types.DownArrow)
	if err != nil {
		return err
	}
	if base.IsFoolsMode {
		err := b.dg.MessageReactionAdd(p.Channel, p.Message, types.UpArrow)
		if err != nil {
			return err
		}
	}

	cnt, err := json.Marshal(p.Data)
	if err != nil {
		return err
	}
	_, err = b.db.Exec("INSERT INTO eod_polls VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ? )", p.Guild, p.Channel, p.Message, p.Kind, p.Value1, p.Value2, p.Value3, p.Value4, string(cnt))

	dat.SavePoll(p.Message, p)

	b.lock.Lock()
	b.dat[p.Guild] = dat
	b.lock.Unlock()
	return err
}
