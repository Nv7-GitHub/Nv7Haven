package discord

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql" // mysql
)

const (
	dbUser     = "u29_c99qmCcqZ3"
	dbPassword = "j8@tJ1vv5d@^xMixUqUl+NmA"
	dbName     = "s29_nv7haven"
	token      = "Nzg4MTg1MzY1NTMzNTU2NzM2.X9f00g.krA6cjfFWYdzbqOPXq8NvRjxb3k"
)

var helpText string
var currHelp string

// Bot is a discord bot
type Bot struct {
	dg              *discordgo.Session
	db              *sql.DB
	memedat         []meme
	memerefreshtime time.Time
	memecache       map[string][]int
	props           map[string]property
}

func (b *Bot) handlers() {
	b.dg.AddHandler(b.giveNum)
	b.dg.AddHandler(b.help)
	b.dg.AddHandler(b.memes)
	b.dg.AddHandler(b.currencyBasics)
	b.dg.AddHandler(b.properties)
}

// InitDiscord creates a discord bot
func InitDiscord() Bot {
	// MySQL DB
	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+os.Getenv("MYSQL_HOST")+":3306)/"+dbName)
	if err != nil {
		panic(err)
	}

	// Discord bot
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	// Help message
	data, err := ioutil.ReadFile("discord/help.txt")
	if err != nil {
		panic(err)
	}
	helpText = string(data)
	data, err = ioutil.ReadFile("discord/currency.txt")
	if err != nil {
		panic(err)
	}
	currHelp = string(data)

	// Init properties
	props := make(map[string]property, 0)
	for _, prop := range upgrades {
		props[prop.ID] = prop
	}

	// Convert users
	res, err := db.Query(`SELECT user, properties FROM currency WHERE properties != "[]"`)
	if err != nil {
		panic(err)
	}
	defer res.Close()
	var user string
	var properties string
	var propDat []struct {
		ID       string
		Upgrades int
	}
	for res.Next() {
		err = res.Scan(&user, &properties)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal([]byte(properties), &propDat)
		if err != nil {
			panic(err)
		}
		newData := make(map[string]int, 0)
		for _, val := range propDat {
			newData[val.ID] = val.Upgrades
		}
		dat, err := json.Marshal(newData)
		if err != nil {
			panic(err)
		}
		_, err = db.Exec("UPDATE currency SET properties=? WHERE user=?", string(dat), user)
		if err != nil {
			panic(err)
		}
	}

	// Set up bot
	b := Bot{
		dg:        dg,
		db:        db,
		memecache: make(map[string][]int, 0),
		props:     props,
	}
	b.handlers()
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	err = dg.Open()
	if err != nil {
		panic(err)
	}
	dg.UpdateStatus(0, "Run 7help to get help on this bot's commands!")
	return b
}

// Close cleans up
func (b *Bot) Close() {
	b.dg.Close()
}

func (b *Bot) help(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "7help") {
		if strings.HasPrefix(m.Content, "7help currency") {
			s.ChannelMessageSend(m.ChannelID, currHelp)
			return
		}
		s.ChannelMessageSend(m.ChannelID, helpText)
		return
	}
}
