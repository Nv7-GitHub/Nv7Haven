package eod

import (
	"database/sql"
	_ "embed"
	"io/ioutil"
	"sync"

	"github.com/bwmarrin/discordgo"
)

const (
	token    = "ODE5MDc2OTIyODY3NzEyMDMx.YEhW1A.iCZTYR_8YH59k7vlYtUM5LZ8Kn8"
	clientID = "819076922867712031"
)

var bot EoD
var lock sync.RWMutex

var about string

// EoD contains the data for an EoD bot
type EoD struct {
	dg  *discordgo.Session
	db  *sql.DB
	dat map[string]serverData // map[guild]data
}

// InitEoD initializes the EoD bot
func InitEoD(db *sql.DB) EoD {
	abt, err := ioutil.ReadFile("eod/about.txt")
	if err != nil {
		panic(err)
	}
	about = string(abt)

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

func (b *EoD) aboutCmd(rsp rsp) {
	rsp.Resp(about)
}

// Close cleans up
func (b *EoD) Close() {
	b.dg.Close()
}
