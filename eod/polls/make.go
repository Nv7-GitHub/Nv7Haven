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

	// Add reactions
	err = dg.MessageReactionAdd(p.Channel, msg.ID, UpArrow)
	if err != nil {
		return types.Error(err)
	}
	err = dg.MessageReactionAdd(p.Channel, msg.ID, DownArrow)
	if err != nil {
		return types.Error(err)
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
		//DM user
		emb, _ := b.makePollEmbed(p)
		dm, err := dg.UserChannelCreate(p.Creator)
		if err != nil {

			return
		}
		guild, _ := dg.Guild(p.Guild)

		upvotetext := "upvotes"
		if p.Upvotes == 0 {
			upvotetext = "upvote"
		}
		downvotetext := "downvotes"
		if p.Downvotes == 0 {
			downvotetext = "downvote"
		}
		msg := sevcord.NewMessage(fmt.Sprintf("Your poll in **%s** was rejected with **%d %s** and **%d %s**.\n\n**Your Poll**", guild.Name, p.Upvotes+1, upvotetext, p.Downvotes+1, downvotetext)).AddEmbed(emb)
		_, err = dg.ChannelMessageSendComplex(dm.ID, msg.Dg())
		if err != nil {
			return
		}

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
