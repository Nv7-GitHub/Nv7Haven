package eod

import (
	"encoding/json"

	"github.com/bwmarrin/discordgo"
)

const upArrow = "⬆️"
const downArrow = "⬇️"

func (b *EoD) createPoll(p poll) error {
	lock.RLock()
	dat, exists := b.dat[p.Guild]
	lock.RUnlock()
	if !exists {
		return nil
	}
	if dat.voteCount == 0 {
		b.handlePollSuccess(p)
		return nil
	}
	if p.Kind == pollCombo {
		m, err := b.dg.ChannelMessageSendEmbed(dat.votingChannel, &discordgo.MessageEmbed{
			Title:       "Combination",
			Description: p.Value1 + " + " + p.Value2 + " = " + p.Value3,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
			Fields: []*discordgo.MessageEmbedField{{
				Value: "Suggested by <@" + p.Value4 + ">",
			}},
		})
		if err != nil {
			return err
		}
		p.Message = m.ID
		err = b.dg.MessageReactionAdd(p.Channel, p.Message, upArrow)
		if err != nil {
			return err
		}
		err = b.dg.MessageReactionAdd(p.Channel, p.Message, downArrow)
		if err != nil {
			return err
		}
	}
	cnt, err := json.Marshal(p.Data)
	if err != nil {
		return err
	}
	_, err = b.db.Exec("INSERT INTO polls VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ? )", p.Guild, p.Channel, p.Message, p.Kind, p.Value1, p.Value2, p.Value3, p.Value4, string(cnt))
	return err
}

func (b *EoD) reactionHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	lock.RLock()
	dat, exists := b.dat[r.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	p, exists := dat.polls[r.MessageID]
	if !exists {
		return
	}
	if r.Emoji.Name == upArrow {
		p.Upvotes++
		if (p.Upvotes - p.Downvotes) >= dat.voteCount {
			b.handlePollSuccess(p)
			delete(dat.polls, r.MessageID)
			b.db.Exec("DELETE FROM eod_polls WHERE guild=? AND channel=? AND message=?", p.Guild, p.Channel, p.Message)
			b.dg.ChannelMessageDelete(p.Channel, p.Message)

			lock.Lock()
			b.dat[r.GuildID] = dat
			lock.Unlock()
			return
		}
	} else if r.Emoji.Name == downArrow {
		p.Downvotes++
		if (p.Downvotes - p.Upvotes) >= dat.voteCount {
			delete(dat.polls, r.MessageID)
			b.db.Exec("DELETE FROM eod_polls WHERE guild=? AND channel=? AND message=?", p.Guild, p.Channel, p.Message)
			b.dg.ChannelMessageDelete(p.Channel, p.Message)

			lock.Lock()
			b.dat[r.GuildID] = dat
			lock.Unlock()
			return
		}
	}
	lock.Lock()
	b.dat[r.GuildID] = dat
	lock.Unlock()
}

func (b *EoD) handlePollSuccess(p poll) {
	switch p.Kind {
	case pollCombo:
		b.elemCreate(p.Value3, p.Value1, p.Value2, p.Value4, p.Guild)
		break
	}
}
