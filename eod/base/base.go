package base

import (
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
)

const configCmdId = "1057375692431568994"

type Base struct {
	s  *sevcord.Sevcord
	db *sqlx.DB

	lock *sync.RWMutex
	mem  map[string]*types.ServerMem // map[guild]data
}

func (b *Base) Init() {
	b.s.AddMiddleware(b.CheckCtx)
	b.s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"stats",
		"View the statistics of this server!",
		b.Stats,
	))
	b.s.RegisterSlashCommand(sevcord.NewSlashCommandGroup("config", "Configure the server!",
		sevcord.NewSlashCommand("voting", "Configure the voting channel!", b.ConfigVoting, sevcord.NewOption("channel", "The new voting channel!", sevcord.OptionKindChannel, true).ChannelFilter(discordgo.ChannelTypeGuildText)),
		sevcord.NewSlashCommand("news", "Configure the news channel!", b.ConfigNews, sevcord.NewOption("channel", "The new news channel!", sevcord.OptionKindChannel, true).ChannelFilter(discordgo.ChannelTypeGuildText)),
		sevcord.NewSlashCommand("votecnt", "Configure the minimum number of votes!", b.ConfigVoteCnt, sevcord.NewOption("count", "The new vote count!", sevcord.OptionKindInt, true)),
		sevcord.NewSlashCommand("pollcnt", "Configure the maximum number of polls!", b.ConfigPollCnt, sevcord.NewOption("count", "The new poll count!", sevcord.OptionKindInt, true)),
		sevcord.NewSlashCommand("playchannels", "Configure the channels in which users can combine elements!", b.ConfigPlayChannels),
		sevcord.NewSlashCommand("voteicons", "Configure the emojis used for voting!", b.ConfigVoteIcons, sevcord.NewOption("upvote", "The new upvote emoji!", sevcord.OptionKindString, true), sevcord.NewOption("downvote", "The new downvote emoji!", sevcord.OptionKindString, true)),
		sevcord.NewSlashCommand("progressicons", "Configure the emojis for progress!", b.ConfigProgIcons, sevcord.NewOption("positive", "The positive progress emoji!", sevcord.OptionKindString, true), sevcord.NewOption("negative", "The negative progress emoji!", sevcord.OptionKindString, true)),
	).RequirePermissions(discordgo.PermissionManageChannels))
	b.s.AddSelectHandler("config_play", b.ConfigPlayChannelsHandler)
}

func NewBase(s *sevcord.Sevcord, db *sqlx.DB) *Base {
	b := &Base{
		lock: &sync.RWMutex{},
		mem:  make(map[string]*types.ServerMem),
		s:    s,
		db:   db,
	}
	b.Init()
	return b
}
