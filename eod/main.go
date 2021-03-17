package eod

import (
	"database/sql"
	"sync"

	"github.com/bwmarrin/discordgo"
)

const (
	token    = "ODE5MDc2OTIyODY3NzEyMDMx.YEhW1A.iCZTYR_8YH59k7vlYtUM5LZ8Kn8"
	clientID = "819076922867712031"
)

var bot EoD
var lock sync.RWMutex

// EoD contains the data for an EoD bot
type EoD struct {
	dg  *discordgo.Session
	db  *sql.DB
	dat map[string]serverData // map[guild]data
}

// InitEoD initializes the EoD bot
func InitEoD(db *sql.DB) EoD {
	// Discord bot
	dg, err := discordgo.New("Bot " + token)
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

	dg.UpdateGameStatus(0, "Type / to see the bot's commands!")
	bot.init()
	return bot
}

// Close cleans up
func (b *EoD) Close() {
	b.dg.Close()
}
