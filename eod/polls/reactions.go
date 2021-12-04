package polls

import (
	"fmt"
	"log"
	"os"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

func (b *Polls) RejectPoll(dat types.ServerDat, p types.OldPoll, messageid, user string) types.ServerDat {
	dat.Lock.Lock()
	delete(dat.Polls, messageid)
	dat.Lock.Unlock()

	b.db.Exec("DELETE FROM eod_polls WHERE guild=? AND channel=? AND message=?", p.Guild, p.Channel, p.Message)
	b.dg.ChannelMessageDelete(p.Channel, p.Message)

	if user != p.Value4 {
		b.dg.ChannelMessageSend(dat.NewsChannel, fmt.Sprintf("%s **Poll Rejected** (By <@%s>)", types.X, p.Value4))

		chn, err := b.dg.UserChannelCreate(p.Value4)
		if err == nil {
			servname, err := b.dg.Guild(p.Guild)
			if err == nil {
				pollemb, err := b.GetPollEmbed(dat, p)
				if err == nil {
					upvotes := ""
					downvotes := ""
					if p.Upvotes > 1 {
						upvotes = "s"
					}
					if p.Downvotes > 1 {
						downvotes = "s"
					}

					b.dg.ChannelMessageSendComplex(chn.ID, &discordgo.MessageSend{
						Content: fmt.Sprintf("Your poll in **%s** was rejected with **%d upvote%s** and **%d downvote%s**.\n\n**Your Poll**:", servname.Name, p.Upvotes, upvotes, p.Downvotes, downvotes),
						Embed:   pollemb,
					})
				}
			}
		}
	}
	return dat
}

func (b *Polls) CheckReactions(dat types.ServerDat, p types.OldPoll, reactor string, downvote bool) (types.ServerDat, bool) {
	if (p.Upvotes - p.Downvotes) >= dat.VoteCount {
		b.dg.ChannelMessageDelete(p.Channel, p.Message)
		b.handlePollSuccess(p)

		dat.Lock.Lock()
		delete(dat.Polls, p.Message)
		dat.Lock.Unlock()

		b.db.Exec("DELETE FROM eod_polls WHERE guild=? AND channel=? AND message=?", p.Guild, p.Channel, p.Message)
		return dat, true
	}

	if ((p.Downvotes - p.Upvotes) >= dat.VoteCount) || (downvote && (reactor == p.Value4)) {
		dat = b.RejectPoll(dat, p, p.Message, reactor)

		return dat, true
	}

	return dat, false
}

func (b *Polls) UnReactionHandler(_ *discordgo.Session, r *discordgo.MessageReactionRemove) {
	if r.UserID == b.dg.State.User.ID {
		return
	}
	b.lock.RLock()
	dat, exists := b.dat[r.GuildID]
	b.lock.RUnlock()
	if !exists {
		return
	}
	p, res := dat.GetPoll(r.MessageID)
	if !res.Exists {
		return
	}
	if r.Emoji.Name == types.DownArrow {
		p.Downvotes--
		dat.SavePoll(r.MessageID, p)
		b.lock.Lock()
		b.dat[r.GuildID] = dat
		b.lock.Unlock()

		var change bool
		dat, change = b.CheckReactions(dat, p, r.UserID, false)
		if change {
			b.lock.Lock()
			b.dat[r.GuildID] = dat
			b.lock.Unlock()
			return
		}
	} else if r.Emoji.Name == types.UpArrow {
		p.Upvotes--
		dat.SavePoll(r.MessageID, p)
		b.lock.Lock()
		b.dat[r.GuildID] = dat
		b.lock.Unlock()

		var change bool
		dat, change = b.CheckReactions(dat, p, r.UserID, false)
		if change {
			b.lock.Lock()
			b.dat[r.GuildID] = dat
			b.lock.Unlock()
			return
		}
	}
	b.lock.Lock()
	b.dat[r.GuildID] = dat
	b.lock.Unlock()
}

func (b *Polls) ReactionHandler(_ *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == b.dg.State.User.ID {
		return
	}

	b.lock.RLock()
	dat, exists := b.dat[r.GuildID]
	b.lock.RUnlock()
	if !exists {
		return
	}

	if len(dat.Polls) == 0 {
		log.SetOutput(os.Stdout)
		log.Println("no polls", r.GuildID)
	}

	p, res := dat.GetPoll(r.MessageID)
	if !res.Exists {
		return
	}

	if r.Emoji.Name == types.UpArrow {
		p.Upvotes++
		dat.SavePoll(r.MessageID, p)
		b.lock.Lock()
		b.dat[r.GuildID] = dat
		b.lock.Unlock()

		var change bool
		dat, change = b.CheckReactions(dat, p, r.UserID, false)
		if change {
			b.lock.Lock()
			b.dat[r.GuildID] = dat
			b.lock.Unlock()
			return
		}
	} else if r.Emoji.Name == types.DownArrow {
		p.Downvotes++
		dat.SavePoll(r.MessageID, p)
		b.lock.Lock()
		b.dat[r.GuildID] = dat
		b.lock.Unlock()

		var change bool
		dat, change = b.CheckReactions(dat, p, r.UserID, true)
		if change {
			b.lock.Lock()
			b.dat[r.GuildID] = dat
			b.lock.Unlock()
			return
		}
	}
	b.lock.Lock()
	b.dat[r.GuildID] = dat
	b.lock.Unlock()
}
