package eod

import (
	"encoding/json"
	"fmt"
	"strings"

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
	switch p.Kind {
	case pollCombo:
		m, err := b.dg.ChannelMessageSendEmbed(dat.votingChannel, &discordgo.MessageEmbed{
			Title:       "Combination",
			Description: p.Value1 + " + " + p.Value2 + " = " + p.Value3 + "\n\n" + "Suggested by <@" + p.Value4 + ">",
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
		})
		if err != nil {
			return err
		}
		p.Message = m.ID

	case pollSign:
		m, err := b.dg.ChannelMessageSendEmbed(dat.votingChannel, &discordgo.MessageEmbed{
			Title:       "Sign Note",
			Description: fmt.Sprintf("**%s**\nNew Note: %s\n\nOld Note: %s\n\nSuggested by <@%s>", p.Value1, p.Value2, p.Value3, p.Value4),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
		})
		if err != nil {
			return err
		}
		p.Message = m.ID

	case pollImage:
		m, err := b.dg.ChannelMessageSendEmbed(dat.votingChannel, &discordgo.MessageEmbed{
			Title:       "Add Image",
			Description: fmt.Sprintf("**%s**\n[New Image](%s)\n[Old Image](%s)\n\nSuggested by <@%s>", p.Value1, p.Value2, p.Value3, p.Value4),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: p.Value2,
			},
		})
		if err != nil {
			return err
		}
		p.Message = m.ID

	case pollCategorize:
		m, err := b.dg.ChannelMessageSendEmbed(dat.votingChannel, &discordgo.MessageEmbed{
			Title:       "Categorize",
			Description: fmt.Sprintf("Elements:\n**%s**\n\nCategory: **%s**\n\nSuggested By <@%s>", strings.Join(p.Data["elems"].([]string), "\n"), p.Value1, p.Value4),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
		})
		if err != nil {
			return err
		}
		p.Message = m.ID
	}
	err := b.dg.MessageReactionAdd(p.Channel, p.Message, upArrow)
	if err != nil {
		return err
	}
	err = b.dg.MessageReactionAdd(p.Channel, p.Message, downArrow)
	if err != nil {
		return err
	}
	cnt, err := json.Marshal(p.Data)
	if err != nil {
		return err
	}
	_, err = b.db.Exec("INSERT INTO eod_polls VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ? )", p.Guild, p.Channel, p.Message, p.Kind, p.Value1, p.Value2, p.Value3, p.Value4, string(cnt))
	if dat.polls == nil {
		dat.polls = make(map[string]poll)
	}
	dat.polls[p.Message] = p

	lock.Lock()
	b.dat[p.Guild] = dat
	lock.Unlock()
	return err
}

func (b *EoD) reactionHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == b.dg.State.User.ID {
		return
	}
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
		dat.polls[r.MessageID] = p
		lock.Lock()
		b.dat[r.GuildID] = dat
		lock.Unlock()
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
		dat.polls[r.MessageID] = p
		lock.Lock()
		b.dat[r.GuildID] = dat
		lock.Unlock()
		if ((p.Downvotes - p.Upvotes) >= dat.voteCount) || (r.UserID == p.Value4) {
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
	lock.RLock()
	dat, exists := b.dat[p.Guild]
	lock.RUnlock()
	if !exists {
		return
	}
	switch p.Kind {
	case pollCombo:
		b.elemCreate(p.Value3, p.Value1, p.Value2, p.Value4, p.Guild)
	case pollSign:
		b.mark(p.Guild, p.Value1, p.Value2, p.Value4)
	case pollImage:
		b.image(p.Guild, p.Value1, p.Value2, p.Value4)
	case pollCategorize:
		els := p.Data["elems"].([]string)
		for _, val := range els {
			b.categorize(val, p.Value1, p.Guild)
		}
		if len(els) == 1 {
			b.dg.ChannelMessageSend(dat.newsChannel, fmt.Sprintf("🗃️ Added **%s** to **%s** (By <@%s>)", els[0], p.Value1, p.Value4))
		} else {
			b.dg.ChannelMessageSend(dat.newsChannel, fmt.Sprintf("🗃️ Added **%d elements** to **%s** (By <@%s>)", len(els), p.Value1, p.Value4))
		}
	}
}
