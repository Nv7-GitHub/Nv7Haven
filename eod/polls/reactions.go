package polls

import (
	"log"
	"os"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

func (b *Polls) RejectPoll(db *eodb.DB, p types.Poll, messageid, user, lasted string) {
	_ = db.DeletePoll(p)
	b.dg.ChannelMessageDelete(p.Channel, p.Message)

	if user != p.Suggestor {
		// Inform them
		b.dg.ChannelMessageSend(db.Config.NewsChannel, db.Config.LangProperty("RejectedPollNews", map[string]any{
			"LastedText": lasted,
			"Creator":    p.Suggestor,
		}))

		chn, err := b.dg.UserChannelCreate(p.Suggestor)
		if err == nil {
			servname, err := b.dg.Guild(p.Guild)
			if err == nil {
				pollemb, err := b.GetPollEmbed(db, p)
				if err == nil {

					b.dg.ChannelMessageSendComplex(chn.ID, &discordgo.MessageSend{
						Content: db.Config.LangProperty("RejectedPollDM", map[string]any{
							"Server":    servname.Name,
							"Upvotes":   p.Upvotes,
							"Downvotes": p.Downvotes,
						}),
						Embed: pollemb,
					})
				}
			}
		}
	}
}

func (b *Polls) CheckReactions(db *eodb.DB, p types.Poll, reactor string, downvote bool) {
	if (p.Upvotes - p.Downvotes) >= db.Config.VoteCount {
		b.dg.ChannelMessageDelete(p.Channel, p.Message)
		b.handlePollSuccess(p)

		db.DeletePoll(p)
		return
	}

	if ((p.Downvotes - p.Upvotes) >= db.Config.VoteCount) || (downvote && (reactor == p.Suggestor)) {
		b.RejectPoll(db, p, p.Message, reactor, b.getLasted(db, p))

		return
	}
}

func (b *Polls) UnReactionHandler(_ *discordgo.Session, r *discordgo.MessageReactionRemove) {
	if r.UserID == b.dg.State.User.ID {
		return
	}
	db, res := b.GetDB(r.GuildID)
	if !res.Exists {
		return
	}
	p, res := db.GetPoll(r.MessageID)
	if !res.Exists {
		return
	}
	if r.Emoji.Name == types.DownArrow {
		p.Downvotes--
		db.SavePoll(p)
		b.CheckReactions(db, p, r.UserID, false)
	} else if r.Emoji.Name == types.UpArrow {
		p.Upvotes--
		db.SavePoll(p)
		b.CheckReactions(db, p, r.UserID, false)
	}
}

func (b *Polls) ReactionHandler(_ *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == b.dg.State.User.ID {
		return
	}
	db, res := b.GetDB(r.GuildID)
	if !res.Exists {
		return
	}

	if len(db.Polls) == 0 {
		log.SetOutput(os.Stdout)
		log.Println("no polls", r.GuildID)
	}

	p, res := db.GetPoll(r.MessageID)
	if !res.Exists {
		return
	}
	if r.Emoji.Name == types.UpArrow {
		p.Upvotes++
		db.SavePoll(p)
		b.CheckReactions(db, p, r.UserID, false)
	} else if r.Emoji.Name == types.DownArrow {
		p.Downvotes++
		db.SavePoll(p)
		b.CheckReactions(db, p, r.UserID, true)
	}
}
