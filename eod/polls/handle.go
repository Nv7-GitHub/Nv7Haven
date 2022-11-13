package polls

import (
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/bwmarrin/discordgo"
)

const UpArrow = "⬆️"
const DownArrow = "⬇️"

func (b *Polls) CreatePoll(c sevcord.Ctx, p *types.Poll) error {
	dg := c.Dg()
	p.Guild = c.Guild()

	// Get embed
	emb, err := b.makePollEmbed(p)
	if err != nil {
		return err
	}
	v := emb.Dg()

	// Get voting channel
	err = b.db.QueryRow("SELECT voting FROM config WHERE guild=$1", c.Guild()).Scan(&p.Channel)
	if err != nil {
		return err
	}

	// Send
	msg, err := dg.ChannelMessageSendEmbed(p.Channel, v)
	if err != nil {
		return err
	}
	p.Message = msg.ID

	// Add reactions
	err = dg.MessageReactionAdd(p.Channel, msg.ID, UpArrow)
	if err != nil {
		return err
	}
	err = dg.MessageReactionAdd(p.Channel, msg.ID, DownArrow)
	if err != nil {
		return err
	}

	// Add to database
	_, err = b.db.NamedExec("INSERT INTO polls (guild, channel, message, kind, creator, createdon, data, upvotes, downvotes) VALUES (:guild, :channel, :message, :kind, :creator, :createdon, :data, :upvotes, :downvotes)", p)
	return err
}

func (b *Polls) reactionHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}

	// Get poll & vote cnt
	var p types.Poll
	err := b.db.Get(&p, "SELECT * FROM polls WHERE guild=$1 AND message=$3", r.GuildID, r.ChannelID, r.MessageID)
	if err != nil {
		return
	}
	var votecnt int
	err = b.db.QueryRow("SELECT votecnt FROM config WHERE guild=$1", r.GuildID).Scan(&votecnt)
	if err != nil {
		return
	}

	// User trying to delete?
	if r.UserID == p.Creator && r.Emoji.Name == DownArrow {
		b.deletePoll(&p, s)
		return
	}

	// Handle
	if r.Emoji.Name == UpArrow {
		p.Upvotes++
	} else {
		p.Downvotes++
	}

	// Update
	_, err = b.db.NamedExec("UPDATE polls SET upvotes=:upvotes, downvotes=:downvotes WHERE guild=:guild AND message=:message", p)
	if err != nil {
		return
	}

	// Check
	b.checkPoll(&p, votecnt, s)
}

func (b *Polls) unReactionHandler(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	if r.UserID == s.State.User.ID {
		return
	}

	// Get poll & vote cnt
	var p types.Poll
	err := b.db.Get(&p, "SELECT * FROM polls WHERE guild=$1 AND message=$3", r.GuildID, r.ChannelID, r.MessageID)
	if err != nil {
		return
	}
	var votecnt int
	err = b.db.QueryRow("SELECT votecnt FROM config WHERE guild=$1", r.GuildID).Scan(&votecnt)
	if err != nil {
		return
	}

	// Handle
	if r.Emoji.Name == UpArrow {
		p.Upvotes--
	} else {
		p.Downvotes--
	}

	// Update
	_, err = b.db.NamedExec("UPDATE polls SET upvotes=:upvotes, downvotes=:downvotes WHERE guild=:guild AND message=:message", p)
	if err != nil {
		return
	}

	// Check
	b.checkPoll(&p, votecnt, s)
}

func (b *Polls) checkPoll(p *types.Poll, votecnt int, dg *discordgo.Session) {
	if p.Upvotes-p.Downvotes >= votecnt {
		b.pollSuccess(p)
		return
	}

	if p.Downvotes-p.Upvotes >= votecnt {
		b.deletePoll(p, dg)
		return
	}
}

func (b *Polls) deletePoll(p *types.Poll, dg *discordgo.Session) {
	// Delete from DB
	_, err := b.db.Exec("DELETE FROM polls WHERE guild=$1 AND message=$2", p.Guild, p.Message)
	if err != nil {
		return
	}

	// Delete from channel
	err = dg.ChannelMessageDelete(p.Channel, p.Message)
	if err != nil {
		return
	}
}
