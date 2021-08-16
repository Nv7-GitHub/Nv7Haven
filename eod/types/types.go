package types

import (
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type ServerDataType int
type PollType int
type PageSwitchType int
type PageSwitchGetter func(PageSwitcher) (string, int, int, error) // text, newPage, maxPages, err
type Empty struct{}

const (
	PlayChannel   = 0
	VotingChannel = 1
	NewsChannel   = 2
	VoteCount     = 3
	PollCount     = 4
	ModRole       = 5
	UserColor     = 6

	PollCombo        = 0
	PollCategorize   = 1
	PollSign         = 2
	PollImage        = 3
	PollUnCategorize = 4
	PollCatImage     = 5

	PageSwitchLdb      = 0
	PageSwitchInv      = 1
	PageSwitchElemSort = 2
)

type ComponentMsg interface {
	Handler(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type ServerData struct {
	PlayChannels  Container // channelID
	UserColors    map[string]int
	VotingChannel string
	NewsChannel   string
	VoteCount     int
	PollCount     int
	ModRole       string                  // role ID
	LastCombs     map[string]Comb         // map[userID]comb
	Inventories   map[string]Container    // map[userID]map[elementName]types.Empty
	Elements      map[string]Element      //map[elementName]element
	Categories    map[string]Category     // map[catName]category
	Polls         map[string]Poll         // map[messageid]poll
	PageSwitchers map[string]PageSwitcher // map[messageid]pageswitcher
	ComponentMsgs map[string]ComponentMsg // map[messageid]componentMsg
	Lock          *sync.RWMutex
}

type PageSwitcher struct {
	Kind       PageSwitchType
	Title      string
	PageGetter PageSwitchGetter
	Thumbnail  string
	Footer     string
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

type Comb struct {
	Elems []string
	Elem3 string
}

type Element struct {
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

type Poll struct {
	Channel string
	Message string
	Guild   string
	Kind    PollType
	Value1  string
	Value2  string
	Value3  string
	Value4  string
	Data    map[string]interface{}

	Upvotes   int
	Downvotes int
}

type Category struct {
	Name     string
	Guild    string
	Elements Container
	Image    string
}

type Msg struct {
	Author    *discordgo.User
	ChannelID string
	GuildID   string
}

type Rsp interface {
	Error(err error) bool
	ErrorMessage(msg string)
	Message(msg string, components ...discordgo.MessageComponent) string
	Embed(emb *discordgo.MessageEmbed, components ...discordgo.MessageComponent) string
	RawEmbed(emb *discordgo.MessageEmbed) string
	Resp(msg string, components ...discordgo.MessageComponent)
	Acknowledge()
	DM(msg string)
}

func NewServerData() ServerData {
	return ServerData{
		Lock:          &sync.RWMutex{},
		ComponentMsgs: make(map[string]ComponentMsg),
		UserColors:    make(map[string]int),
	}
}

type Container map[string]Empty

func (c Container) Contains(elem string) bool {
	_, exists := c[strings.ToLower(elem)]
	return exists
}

func (c Container) Add(elem string) {
	c[strings.ToLower(elem)] = Empty{}
}
