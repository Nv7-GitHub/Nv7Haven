package eod

import (
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type serverDataType int
type pollType int
type pageSwitchType int
type pageSwitchGetter func(pageSwitcher) (string, int, int, error) // text, newPage, maxPages, err

const (
	playChannel   = 0
	votingChannel = 1
	newsChannel   = 2
	voteCount     = 3
	pollCount     = 4
	modRole       = 5

	pollCombo        = 0
	pollCategorize   = 1
	pollSign         = 2
	pollImage        = 3
	pollUnCategorize = 4
	pollCatImage     = 5

	pageSwitchLdb      = 0
	pageSwitchInv      = 1
	pageSwitchElemSort = 2
)

type empty struct{}

type componentMsg interface {
	handler(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type serverData struct {
	playChannels  map[string]empty // channelID
	votingChannel string
	newsChannel   string
	voteCount     int
	pollCount     int
	modRole       string                      // role ID
	combCache     map[string]comb             // map[userID]comb
	invCache      map[string]map[string]empty // map[userID]map[elementName]empty
	elemCache     map[string]element          //map[elementName]element
	catCache      map[string]category         // map[catName]category
	polls         map[string]poll             // map[messageid]poll
	pageSwitchers map[string]pageSwitcher     // map[messageid]pageswitcher
	componentMsgs map[string]componentMsg     // map[messageid]componentMsg
	lock          *sync.RWMutex
}

type pageSwitcher struct {
	Kind       pageSwitchType
	Title      string
	PageGetter pageSwitchGetter
	Thumbnail  string
	PageLength int

	// Inv
	Items []string

	// Ldb
	User string
	Sort string

	// Element sorting
	Query  string
	Length int

	// Don't need to set these
	Guild string
	Page  int
}

type comb struct {
	elems []string
	elem3 string
}

type element struct {
	ID         int
	Name       string
	Image      string
	Guild      string
	Comment    string
	Creator    string
	CreatedOn  time.Time
	Parents    []string
	Complexity int
	Difficulty int
	UsedIn     int
}

type poll struct {
	Channel string
	Message string
	Guild   string
	Kind    pollType
	Value1  string
	Value2  string
	Value3  string
	Value4  string
	Data    map[string]interface{}

	Upvotes   int
	Downvotes int
}

type category struct {
	Name     string
	Guild    string
	Elements map[string]empty
	Image    string
}

type msg struct {
	Author    *discordgo.User
	ChannelID string
	GuildID   string
}

type rsp interface {
	Error(err error) bool
	ErrorMessage(msg string)
	Message(msg string, components ...discordgo.MessageComponent) string
	Embed(emb *discordgo.MessageEmbed, components ...discordgo.MessageComponent) string
	RawEmbed(emb *discordgo.MessageEmbed) string
	Resp(msg string, components ...discordgo.MessageComponent)
	Acknowledge()
	DM(msg string)
}
