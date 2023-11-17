package polls

import (
	"fmt"
	"log"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/bwmarrin/discordgo"
)

func (b *Polls) CreatePoll(c sevcord.Ctx, p *types.Poll) types.Resp {
	dg := c.Dg()
	p.Guild = c.Guild()
	p.Creator = c.Author().User.ID
	p.CreatedOn = time.Now()

	// Check poll count
	var pollcnt int
	err := b.db.QueryRow("SELECT pollcnt FROM config WHERE guild=$1", c.Guild()).Scan(&pollcnt)
	if err != nil {
		return types.Error(err)
	}
	if pollcnt > 0 {
		var cnt int
		err = b.db.QueryRow("SELECT COUNT(*) FROM polls WHERE guild=$1 AND creator=$2", c.Guild(), p.Creator).Scan(&cnt)
		if err != nil {
			return types.Error(err)
		}
		if cnt >= pollcnt {
			return types.Fail("Poll limit reached!")
		}
	}

	// Get embed
	emb, err := b.makePollEmbed(p)
	if err != nil {
		return types.Error(err)
	}
	v := emb.Dg()

	// Get voting channel
	err = b.db.QueryRow("SELECT voting FROM config WHERE guild=$1", c.Guild()).Scan(&p.Channel)
	if err != nil {
		return types.Error(err)
	}

	// Send
	msg, err := dg.ChannelMessageSendEmbed(p.Channel, v)
	if err != nil {
		return types.Error(err)
	}
	p.Message = msg.ID

	// Add reactions at the same time
	errs := make(chan error)
	go func() {
		errs <- dg.MessageReactionAdd(p.Channel, msg.ID, UpArrow)
	}()
	go func() {
		errs <- dg.MessageReactionAdd(p.Channel, msg.ID, DownArrow)
	}()
	err = <-errs
	if err != nil {
		return types.Error(err) // Up arrow error
	}
	err = <-errs
	if err != nil {
		return types.Error(err) // Down arrow error
	}

	// Add to database
	_, err = b.db.NamedExec("INSERT INTO polls (guild, channel, message, kind, creator, createdon, data, upvotes, downvotes) VALUES (:guild, :channel, :message, :kind, :creator, :createdon, :data, :upvotes, :downvotes)", p)
	if err != nil {
		return types.Error(err)
	}
	return types.Ok()
}

func (b *Polls) checkPoll(p *types.Poll, votecnt int, dg *discordgo.Session) {
	if p.Upvotes-p.Downvotes >= votecnt {
		b.pollSuccess(p, dg)
		return
	}

	if p.Downvotes-p.Upvotes >= votecnt {
		var news string
		err := b.db.QueryRow(`SELECT news FROM config WHERE guild=$1`, p.Guild).Scan(&news)
		if err != nil {
			log.Println("news err", err)
			return
		}
		_, err = dg.ChannelMessageSend(news, fmt.Sprintf("❌ **Poll Rejected** %s", b.pollContextMsg(p)))
		if err != nil {
			log.Println("news err", err)
		}
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
