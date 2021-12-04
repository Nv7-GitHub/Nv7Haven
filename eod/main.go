package eod

import (
	_ "embed"
	"strings"
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/basecmds"
	"github.com/Nv7-Github/Nv7Haven/eod/categories"
	"github.com/Nv7-Github/Nv7Haven/eod/elements"
	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/polls"
	"github.com/Nv7-Github/Nv7Haven/eod/treecmds"
	"github.com/Nv7-Github/Nv7Haven/eod/types"

	"github.com/bwmarrin/discordgo"
)

const (
	clientID = "819076922867712031"
	status   = "Use /help to view the bot's commands!"
)

//go:embed token.txt
var token string

var bot EoD
var lock = &sync.RWMutex{}

// EoD contains the data for an EoD bot
type EoD struct {
	*eodb.Data

	dg *discordgo.Session

	// Subsystems
	base       *base.Base
	treecmds   *treecmds.TreeCmds
	polls      *polls.Polls
	basecmds   *basecmds.BaseCmds
	categories *categories.Categories
	elements   *elements.Elements
}

// InitEoD initializes the EoD bot
func InitEoD() EoD {
	// Discord bot
	dg, err := discordgo.New("Bot " + strings.TrimSpace(token))
	if err != nil {
		panic(err)
	}

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions | discordgo.IntentsGuildMembers | discordgo.IntentsGuilds)
	err = dg.Open()
	if err != nil {
		panic(err)
	}

	db, err := eodb.NewData("data/eod")
	bot = EoD{
		Data: db,

		dg: dg,
	}

	dg.UpdateGameStatus(0, status)
	bot.init()

	// FOOLS
	bot.base.InitFools(foolsRaw)
	if base.IsFoolsMode {
		types.MaxComboLength = 2
	}

	return bot
}

// Close cleans up
func (b *EoD) Close() {
	b.Data.RLock()
	for _, db := range b.Data.DB {
		db.Close()
	}
	b.Data.RUnlock()
	b.dg.Close()
}
