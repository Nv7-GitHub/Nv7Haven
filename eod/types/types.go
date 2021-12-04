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
	PollColor        = 6
	PollCatColor     = 7

	PageSwitchLdb = 0
	PageSwitchInv = 1
)

type ComponentMsg interface {
	Handler(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type ServerConfig struct {
	UserColors    map[string]int
	VotingChannel string
	NewsChannel   string
	VoteCount     int
	PollCount     int
	ModRole       string    // role ID
	PlayChannels  Container // channelID
}

type ServerData struct {
	LastCombs     map[string]Comb         // map[userID]comb
	PageSwitchers map[string]PageSwitcher // map[messageid]pageswitcher
	ComponentMsgs map[string]ComponentMsg // map[messageid]componentMsg
	ElementMsgs   map[string]string       // map[messageid]elemname
}

type ServerDat struct {
	ServerData
	ServerConfig

	Inventories map[string]Inventory   // map[userID]map[elementName]types.Empty
	Elements    map[string]OldElement  //map[elementName]element
	Combos      map[string]string      // map[elems]elem3
	Categories  map[string]OldCategory // map[catName]category
	Polls       map[string]OldPoll     // map[messageid]poll
	Lock        *sync.RWMutex
}

type PageSwitcher struct {
	Kind       PageSwitchType
	Title      string
	PageGetter PageSwitchGetter
	Thumbnail  string
	Footer     string
	PageLength int
	Color      int

	// Inv
	Items []string

	// Ldb
	Users   []string
	Cnts    []int
	UserPos int
	User    string

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
	Color      int
	Guild      string
	Comment    string
	Creator    string
	CreatedOn  time.Time
	Parents    []int
	Complexity int
	Difficulty int
	UsedIn     int
	TreeSize   int
}

type OldElement struct {
	ID         int
	Name       string
	Image      string
	Color      int
	Guild      string
	Comment    string
	Creator    string
	CreatedOn  time.Time
	Parents    []string
	Complexity int
	Difficulty int
	UsedIn     int
	TreeSize   int
}

type Poll struct {
	Channel string
	Message string
	Guild   string
	Kind    PollType

	// Data, pointers to different types with omitempty so that you can selectively have some data

	Upvotes   int
	Downvotes int
}

type OldPoll struct {
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
	Lock *sync.RWMutex `json:"-"`

	Name     string
	Guild    string
	Elements map[int]Empty
	Image    string
	Color    int
}

type OldCategory struct {
	Name     string
	Guild    string
	Elements map[string]Empty
	Image    string
	Color    int
}

type Msg struct {
	Author    *discordgo.User
	ChannelID string
	GuildID   string
}

type Inventory struct {
	Elements Container
	MadeCnt  int
	User     string
}

type Rsp interface {
	Error(err error) bool
	ErrorMessage(msg string) string
	Message(msg string, components ...discordgo.MessageComponent) string
	Embed(emb *discordgo.MessageEmbed, components ...discordgo.MessageComponent) string
	RawEmbed(emb *discordgo.MessageEmbed) string
	Resp(msg string, components ...discordgo.MessageComponent)
	Acknowledge()
	DM(msg string)
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		UserColors:   make(map[string]int),
		PlayChannels: make(Container),
	}
}

func NewServerData() ServerDat {
	return ServerDat{
		ServerData: ServerData{
			LastCombs:     make(map[string]Comb),
			PageSwitchers: make(map[string]PageSwitcher),
			ComponentMsgs: make(map[string]ComponentMsg),
			ElementMsgs:   make(map[string]string),
		},
		ServerConfig: *NewServerConfig(),

		Lock:        &sync.RWMutex{},
		Polls:       make(map[string]OldPoll),
		Elements:    make(map[string]OldElement),
		Combos:      make(map[string]string),
		Categories:  make(map[string]OldCategory),
		Inventories: make(map[string]Inventory),
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

func NewInventory(user string) Inventory {
	return Inventory{
		Elements: make(map[string]Empty),
		User:     user,
	}
}

type ElemContainer struct {
	sync.RWMutex
	Data map[int]Empty

	Id string
}

func (e *ElemContainer) Contains(val int) bool {
	e.RLock()
	_, contains := e.Data[val]
	e.RUnlock()
	return contains
}

func (e *ElemContainer) Add(val int) {
	e.Lock()
	e.Data[val] = Empty{}
	e.Unlock()
}

func NewElemContainer(data map[int]Empty, id string) *ElemContainer {
	return &ElemContainer{
		Data: data,
		Id:   id,
	}
}
