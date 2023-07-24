package base

import (
	"fmt"
	"log"
	"strings"

	"github.com/Nv7-Github/sevcord/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/lib/pq"
)

func (b *Base) configNewsMessage(c sevcord.Ctx, message string) {
	var news string
	err := b.db.QueryRow(`SELECT news FROM config WHERE guild=$1`, c.Guild()).Scan(&news)
	if err != nil {
		log.Println("news err", err)
		return
	}
	_, err = c.Dg().ChannelMessageSend(news, fmt.Sprintf("⚙ "+message))
	if err != nil {
		log.Println("news err", err)
	}
}

func (b *Base) configChannel(c sevcord.Ctx, field string, value any, message string) {
	c.Acknowledge()

	_, err := b.db.Exec(fmt.Sprintf("UPDATE config SET %s=$1 WHERE guild=$2", field), value, c.Guild())
	if err != nil {
		b.Error(c, err)
		return
	}

	c.Respond(sevcord.NewMessage(fmt.Sprintf("Successfully updated %s!", message)))
	b.configNewsMessage(c, fmt.Sprintf("Changed Config - **%s**", strings.Title(message)))
}

func (b *Base) ConfigVoting(c sevcord.Ctx, opts []any) {
	b.configChannel(c, "voting", opts[0], "voting channel")
}

func (b *Base) ConfigNews(c sevcord.Ctx, opts []any) {
	b.configChannel(c, "news", opts[0], "news channel")
}

func (b *Base) ConfigVoteCnt(c sevcord.Ctx, opts []any) {
	b.configChannel(c, "votecnt", opts[0], "vote count")
}

func (b *Base) ConfigPollCnt(c sevcord.Ctx, opts []any) {
	b.configChannel(c, "pollcnt", opts[0], "max poll count")
}

func (b *Base) ConfigPlayChannels(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	_, err := b.db.Exec("UPDATE config SET play=$1 WHERE guild=$2", pq.Array([]string{}), c.Guild())
	if err != nil {
		b.Error(c, err)
		return
	}

	c.Respond(sevcord.NewMessage("**PLAY CHANNELS HAVE BEEN RESET**\nUpdate them below!").AddComponentRow(
		sevcord.NewSelect("Play channels", "config_play", c.Author().User.ID).
			SetKind(sevcord.SelectKindChannel).
			ChannelMenuFilter(discordgo.ChannelTypeGuildText).
			SetRange(0, 25),
	))
	b.configNewsMessage(c, "Changed Config - **Play Channels**")
}

func (b *Base) ConfigPlayChannelsHandler(c sevcord.Ctx, params string, opts []string) {
	c.Acknowledge()

	if c.Author().User.ID != params {
		c.Respond(sevcord.NewMessage("You are not authorized to use this!"))
		return
	}

	_, err := b.db.Exec("UPDATE config SET play=$1 WHERE guild=$2", pq.Array(opts), c.Guild())
	if err != nil {
		b.Error(c, err)
		return
	}

	c.Respond(sevcord.NewMessage("Successfully updated play channels!"))
}
