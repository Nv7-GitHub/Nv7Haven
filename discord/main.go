package discord

import (
	"database/sql"
	"os"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql" // mysql
)

const (
	dbUser     = "u29_c99qmCcqZ3"
	dbPassword = "j8@tJ1vv5d@^xMixUqUl+NmA"
	dbName     = "s29_nv7haven"
	token      = "Nzg4MTg1MzY1NTMzNTU2NzM2.X9f00g.krA6cjfFWYdzbqOPXq8NvRjxb3k"
)

// Bot is a discord bot
type Bot struct {
	dg *discordgo.Session
	db *sql.DB
}

func (b *Bot) handlers() {
	b.dg.AddHandler(b.giveNum)
}

// InitDiscord creates a discord bot
func InitDiscord() Bot {
	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	if err != nil {
		panic(err)
	}
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	b := Bot{
		dg: dg,
		db: db,
	}
	b.handlers()
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	err = dg.Open()
	if err != nil {
		panic(err)
	}
	return b
}

func (b *Bot) handle(err error, m *discordgo.MessageCreate) bool {
	if err != nil {
		b.dg.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
		return true
	}
	return false
}

// Close cleans up
func (b *Bot) Close() {
	b.dg.Close()
}
