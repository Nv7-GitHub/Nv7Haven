package polls

import (
	"fmt"
	"log"
	"os"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

func (b *Polls) RejectPoll(db *eodb.DB, p types.Poll, messageid, user string) {
	_ = db.DeletePoll(p)
	b.dg.ChannelMessageDelete(p.Channel, p.Message)

	if user != p.Suggestor {
		// Inform them
		lasted := fmt.Sprintf(db.Config.LangProperty("Lasted"), time.Since(p.CreatedOn.Time).Round(time.Second).String()) + " • "
		b.dg.ChannelMessageSend(db.Config.NewsChannel, fmt.Sprintf(db.Config.LangProperty("RejectedPollNews"), types.X, lasted, p.Suggestor))

		chn, err := b.dg.UserChannelCreate(p.Suggestor)
		if err == nil {
			servname, err := b.dg.Guild(p.Guild)
			if err == nil {
				pollemb, err := b.GetPollEmbed(db, p)
				if err == nil {

					b.dg.ChannelMessageSendComplex(chn.ID, &discordgo.MessageSend{
						Content: fmt.Sprintf(db.Config.LangProperty("RejectedPollDM"), servname.Name, p.Upvotes, p.Downvotes),
						Embed:   pollemb,
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
		b.RejectPoll(db, p, p.Message, reactor)

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
