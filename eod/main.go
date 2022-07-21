package eod

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/Nv7-Github/Nv7Haven/db"
	"github.com/Nv7-Github/Nv7Haven/eod/api"
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/basecmds"
	"github.com/Nv7-Github/Nv7Haven/eod/categories"
	"github.com/Nv7-Github/Nv7Haven/eod/elements"
	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/polls"
	"github.com/Nv7-Github/Nv7Haven/eod/treecmds"

	"github.com/bwmarrin/discordgo"
)

const (
	clientID = "819076922867712031"
	status   = "Use /help to view the bot's commands!"
)

//go:embed token.txt
var token string

var bot EoD

// EoD contains the data for an EoD bot
type EoD struct {
	*eodb.Data

	db *db.DB
	dg *discordgo.Session

	// Subsystems
	base       *base.Base
	treecmds   *treecmds.TreeCmds
	polls      *polls.Polls
	basecmds   *basecmds.BaseCmds
	categories *categories.Categories
	elements   *elements.Elements
	api        *api.API
}

// InitEoD initializes the EoD bot
func InitEoD(sqldb *db.DB) EoD {
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

	start := time.Now()
	fmt.Println("Loading DB...")
	db, err := eodb.NewData("data/eod")
	if err != nil {
		panic(err)
	}
	fmt.Println("started in", time.Since(start))
	bot = EoD{
		Data: db,

		dg: dg,
		db: sqldb,
	}

	dg.UpdateGameStatus(0, status)
	bot.init()

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
