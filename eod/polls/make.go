package polls

import (
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/bwmarrin/discordgo"
)

func (b *Polls) CreatePoll(c sevcord.Ctx, p *types.Poll) error {
	dg := c.Dg()
	p.Guild = c.Guild()
	p.Creator = c.Author().User.ID
	p.CreatedOn = time.Now()

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

func (b *Polls) checkPoll(p *types.Poll, votecnt int, dg *discordgo.Session) {
	if p.Upvotes-p.Downvotes >= votecnt {
		b.pollSuccess(p, dg)
		return
	}

	if p.Downvotes-p.Upvotes >= votecnt {
		b.deletePoll(p, dg)
		return
	}
}

func (b *Polls) deletePoll(p *types.Poll, dg *discordgo.Session) {
	// Delete from channel
	err := dg.ChannelMessageDelete(p.Channel, p.Message)
	if err != nil {
		return
	}

	// Delete from DB
	_, err = b.db.Exec("DELETE FROM polls WHERE guild=$1 AND message=$2", p.Guild, p.Message)
	if err != nil {
		return
	}
}
