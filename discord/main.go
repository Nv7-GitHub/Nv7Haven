package discord

import (
	"database/sql"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	_ "embed"

	"github.com/Nv7-Github/Nv7Haven/elemental"
	"github.com/bwmarrin/discordgo"
)

const (
	clientID = "788185365533556736"
)

//go:embed token.txt
var token string

var helpText string
var currHelp string
var bot Bot

// Bot is a discord bot
type Bot struct {
	dg    *discordgo.Session
	db    *sql.DB
	props map[string]property

	memerefreshtime time.Time
	memedat         []meme
	memecache       map[string]map[int]empty
	cmemedat        []meme
	cmemecache      map[string]map[int]empty
	pmemedat        []meme
	pmemecache      map[string]map[int]empty

	mathvars map[string]map[string]interface{}

	prefixcache map[string]string

	pages map[string]reactionMsg

	e      elemental.Elemental
	combos map[string]comb
}

type comb struct {
	elem1 string
	elem2 string
	elem3 string
}

func (b *Bot) handlers() {
	b.dg.AddHandler(b.giveNum)
	b.dg.AddHandler(b.help)
	b.dg.AddHandler(b.memes)
	b.dg.AddHandler(b.currencyBasics)
	b.dg.AddHandler(b.properties)
	b.dg.AddHandler(b.specials)
	b.dg.AddHandler(b.mod)
	b.dg.AddHandler(b.other)
	b.dg.AddHandler(b.memeGen)
	b.dg.AddHandler(b.math)
	b.dg.AddHandler(b.elementalHandler)
	b.dg.AddHandler(b.pageSwitchHandler)
	for _, v := range commands {
		go func(val *discordgo.ApplicationCommand) {
			_, err := b.dg.ApplicationCommandCreate(clientID, "806258286043070545", val)
			if err != nil {
				panic(err)
			}
		}(v)
	}
	b.dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

// InitDiscord creates a discord bot
func InitDiscord(db *sql.DB, e elemental.Elemental) Bot {
	// Init
	rand.Seed(time.Now().UnixNano())

	// Discord bot
	dg, err := discordgo.New("Bot " + strings.TrimSpace(token))
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
	props := make(map[string]property)
	for _, prop := range upgrades {
		props[prop.ID] = prop
	}

	// Set up bot
	b := Bot{
		dg:    dg,
		db:    db,
		e:     e,
		props: props,

		mathvars:    make(map[string]map[string]interface{}),
		pages:       make(map[string]reactionMsg),
		prefixcache: make(map[string]string),
		combos:      make(map[string]comb),
	}
	b.handlers()
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions)
	err = dg.Open()
	if err != nil {
		panic(err)
	}
	dg.UpdateGameStatus(0, "Run 7help to get help on this bot's commands!")
	bot = b
	return b
}

// Close cleans up
func (b *Bot) Close() {
	b.dg.Close()
}

func (b *Bot) help(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if b.startsWith(m, "7help currency") {
		s.ChannelMessageSend(m.ChannelID, currHelp)
		return
	}

	if b.startsWith(m, "7help") {
		s.ChannelMessageSend(m.ChannelID, helpText)
		return
	}
}

type msg struct {
	Author    *discordgo.User
	ChannelID string
	GuildID   string
}

type rsp interface {
	Error(err error) bool
	ErrorMessage(msg string)
	Message(msg string) string
	Embed(emb *discordgo.MessageEmbed) string
	Resp(msg string)
}
