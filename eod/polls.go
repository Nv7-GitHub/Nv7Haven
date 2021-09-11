package eod

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

const upArrow = "⬆️"
const downArrow = "⬇️"

func (b *EoD) createPoll(p types.Poll) error {
	lock.RLock()
	dat, exists := b.dat[p.Guild]
	lock.RUnlock()
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
	switch p.Kind {
	case types.PollCombo:
		txt := ""
		elems, ok := p.Data["elems"].([]string)
		if !ok {
			elemDat := p.Data["elems"].([]interface{})
			elems = make([]string, len(elemDat))
			for i, val := range elemDat {
				elems[i] = val.(string)
			}
		}
		if len(elems) < 1 {
			return errors.New("error: combo must have at least one element")
		}
		for _, val := range elems {
			el, _ := dat.GetElement(val)
			txt += el.Name + " + "
		}
		txt = txt[:len(txt)-2]
		if len(elems) == 1 {
			el, _ := dat.GetElement(elems[0])
			txt += " + " + el.Name
		}
		txt += " = " + p.Value3
		m, err := b.dg.ChannelMessageSendEmbed(dat.VotingChannel, &discordgo.MessageEmbed{
			Title:       "Combination",
			Description: txt + "\n\n" + "Suggested by <@" + p.Value4 + ">",
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
		})
		if err != nil {
			return err
		}
		p.Message = m.ID

	case types.PollSign:
		m, err := b.dg.ChannelMessageSendEmbed(dat.VotingChannel, &discordgo.MessageEmbed{
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

	case types.PollImage:
		description := fmt.Sprintf("**%s**\n[New Image](%s)\n[Old Image](%s)\n\nSuggested by <@%s>", p.Value1, p.Value2, p.Value3, p.Value4)
		if p.Value3 == "" {
			description = fmt.Sprintf("**%s**\n[New Image](%s)\n\nSuggested by <@%s>", p.Value1, p.Value2, p.Value4)
		}
		m, err := b.dg.ChannelMessageSendEmbed(dat.VotingChannel, &discordgo.MessageEmbed{
			Title:       "Add Image",
			Description: description,
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

	case types.PollCategorize:
		data, ok := p.Data["elems"].([]string)
		if !ok {
			dat := p.Data["elems"].([]interface{})
			data = make([]string, len(dat))
			for i, val := range dat {
				data[i] = val.(string)
			}
		}
		p.Data["elems"] = data
		m, err := b.dg.ChannelMessageSendEmbed(dat.VotingChannel, &discordgo.MessageEmbed{
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
	case types.PollUnCategorize:
		data, ok := p.Data["elems"].([]string)
		if !ok {
			dat := p.Data["elems"].([]interface{})
			data = make([]string, len(dat))
			for i, val := range dat {
				data[i] = val.(string)
			}
		}
		p.Data["elems"] = data
		m, err := b.dg.ChannelMessageSendEmbed(dat.VotingChannel, &discordgo.MessageEmbed{
			Title:       "Un-Categorize",
			Description: fmt.Sprintf("Elements:\n**%s**\n\nCategory: **%s**\n\nSuggested By <@%s>", strings.Join(p.Data["elems"].([]string), "\n"), p.Value1, p.Value4),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
		})
		if err != nil {
			return err
		}
		p.Message = m.ID

	case types.PollCatImage:
		description := fmt.Sprintf("**%s**\n[New Image](%s)\n[Old Image](%s)\n\nSuggested by <@%s>", p.Value1, p.Value2, p.Value3, p.Value4)
		if p.Value3 == "" {
			description = fmt.Sprintf("**%s**\n[New Image](%s)\n\nSuggested by <@%s>", p.Value1, p.Value2, p.Value4)
		}
		m, err := b.dg.ChannelMessageSendEmbed(dat.VotingChannel, &discordgo.MessageEmbed{
			Title:       "Add Category Image",
			Description: description,
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
	}

	if !isFoolsMode {
		err := b.dg.MessageReactionAdd(p.Channel, p.Message, upArrow)
		if err != nil {
			return err
		}
	}
	err := b.dg.MessageReactionAdd(p.Channel, p.Message, downArrow)
	if err != nil {
		return err
	}
	if isFoolsMode {
		err := b.dg.MessageReactionAdd(p.Channel, p.Message, upArrow)
		if err != nil {
			return err
		}
	}

	cnt, err := json.Marshal(p.Data)
	if err != nil {
		return err
	}
	_, err = b.db.Exec("INSERT INTO eod_polls VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ? )", p.Guild, p.Channel, p.Message, p.Kind, p.Value1, p.Value2, p.Value3, p.Value4, string(cnt))
	if dat.Polls == nil {
		dat.Polls = make(map[string]types.Poll)
	}

	dat.SavePoll(p.Message, p)

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
	p, res := dat.GetPoll(r.MessageID)
	if !res.Exists {
		return
	}
	if r.Emoji.Name == upArrow {
		p.Upvotes++
		dat.SavePoll(r.MessageID, p)
		lock.Lock()
		b.dat[r.GuildID] = dat
		lock.Unlock()
		if (p.Upvotes - p.Downvotes) >= dat.VoteCount {
			b.dg.ChannelMessageDelete(p.Channel, p.Message)
			b.handlePollSuccess(p)
			delete(dat.Polls, r.MessageID)
			b.db.Exec("DELETE FROM eod_polls WHERE guild=? AND channel=? AND message=?", p.Guild, p.Channel, p.Message)
			lock.Lock()
			b.dat[r.GuildID] = dat
			lock.Unlock()
			return
		}
	} else if r.Emoji.Name == downArrow {
		p.Downvotes++
		dat.SavePoll(r.MessageID, p)
		lock.Lock()
		b.dat[r.GuildID] = dat
		lock.Unlock()
		if ((p.Downvotes - p.Upvotes) >= dat.VoteCount) || (r.UserID == p.Value4) {
			delete(dat.Polls, r.MessageID)
			b.db.Exec("DELETE FROM eod_polls WHERE guild=? AND channel=? AND message=?", p.Guild, p.Channel, p.Message)
			b.dg.ChannelMessageDelete(p.Channel, p.Message)
			b.dg.ChannelMessageSend(dat.NewsChannel, fmt.Sprintf("%s **Poll Rejected** (By <@%s>)", x, p.Value4))

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

func (b *EoD) handlePollSuccess(p types.Poll) {
	lock.RLock()
	dat, exists := b.dat[p.Guild]
	lock.RUnlock()
	if !exists {
		return
	}

	controversial := dat.VoteCount != 0 && float32(p.Downvotes)/float32(dat.VoteCount) >= 0.3
	controversialTxt := ""
	if controversial {
		controversialTxt = "🌩️"
	}

	switch p.Kind {
	case types.PollCombo:
		els, ok := p.Data["elems"].([]string)
		if !ok {
			dat := p.Data["elems"].([]interface{})
			els = make([]string, len(dat))
			for i, val := range dat {
				els[i] = val.(string)
			}
		}
		b.elemCreate(p.Value3, els, p.Value4, controversialTxt, p.Guild)
	case types.PollSign:
		b.mark(p.Guild, p.Value1, p.Value2, p.Value4, controversialTxt)
	case types.PollImage:
		b.image(p.Guild, p.Value1, p.Value2, p.Value4, controversialTxt)
	case types.PollCategorize:
		els, ok := p.Data["elems"].([]string)
		if !ok {
			dat := p.Data["elems"].([]interface{})
			els := make([]string, len(dat))
			for i, val := range dat {
				els[i] = val.(string)
			}
		}
		for _, val := range els {
			b.categorize(val, p.Value1, p.Guild)
		}
		if len(els) == 1 {
			b.dg.ChannelMessageSend(dat.NewsChannel, fmt.Sprintf("🗃️ Added **%s** to **%s** (By <@%s>)%s", els[0], p.Value1, p.Value4, controversialTxt))
		} else {
			b.dg.ChannelMessageSend(dat.NewsChannel, fmt.Sprintf("🗃️ Added **%d elements** to **%s** (By <@%s>)%s", len(els), p.Value1, p.Value4, controversialTxt))
		}
	case types.PollUnCategorize:
		els := p.Data["elems"].([]string)
		for _, val := range els {
			b.unCategorize(val, p.Value1, p.Guild)
		}
		if len(els) == 1 {
			b.dg.ChannelMessageSend(dat.NewsChannel, fmt.Sprintf("🗃️ Removed **%s** from **%s** (By <@%s>)%s", els[0], p.Value1, p.Value4, controversialTxt))
		} else {
			b.dg.ChannelMessageSend(dat.NewsChannel, fmt.Sprintf("🗃️ Removed **%d elements** from **%s** (By <@%s>)%s", len(els), p.Value1, p.Value4, controversialTxt))
		}
	case types.PollCatImage:
		b.catImage(p.Guild, p.Value1, p.Value2, p.Value4, controversialTxt)
	}
}

func (b *EoD) unReactionHandler(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	if r.UserID == b.dg.State.User.ID {
		return
	}
	lock.RLock()
	dat, exists := b.dat[r.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	p, exists := dat.Polls[r.MessageID]
	if !exists {
		return
	}
	if r.Emoji.Name == downArrow {
		p.Downvotes--
		dat.Polls[r.MessageID] = p
		lock.Lock()
		b.dat[r.GuildID] = dat
		lock.Unlock()
		if (p.Upvotes - p.Downvotes) >= dat.VoteCount {
			b.dg.ChannelMessageDelete(p.Channel, p.Message)
			b.handlePollSuccess(p)
			delete(dat.Polls, r.MessageID)
			b.db.Exec("DELETE FROM eod_polls WHERE guild=? AND channel=? AND message=?", p.Guild, p.Channel, p.Message)
			lock.Lock()
			b.dat[r.GuildID] = dat
			lock.Unlock()
			return
		}
	} else if r.Emoji.Name == upArrow {
		p.Upvotes--
		dat.Polls[r.MessageID] = p
		lock.Lock()
		b.dat[r.GuildID] = dat
		lock.Unlock()
		if (p.Downvotes - p.Upvotes) >= dat.VoteCount {
			delete(dat.Polls, r.MessageID)
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
