package eod

import (
	"database/sql"
	_ "embed"
	"strings"

	"github.com/bwmarrin/discordgo"
	deadlock "github.com/sasha-s/go-deadlock"
)

const (
	clientID = "819076922867712031"
	status   = "Use /help to view the bot's commands!"
)

//go:embed token.txt
var token string

var bot EoD
var lock deadlock.RWMutex

// EoD contains the data for an EoD bot
type EoD struct {
	dg  *discordgo.Session
	db  *sql.DB
	dat map[string]serverData // map[guild]data
}

// InitEoD initializes the EoD bot
func InitEoD(db *sql.DB) EoD {
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

	bot = EoD{
		dg:  dg,
		db:  db,
		dat: make((map[string]serverData)),
	}

	dg.UpdateGameStatus(0, status)
	bot.init()
	return bot
}

// Close cleans up
func (b *EoD) Close() {
	b.dg.Close()
}
