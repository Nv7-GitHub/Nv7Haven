package types

import (
	"strconv"
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
	*sync.RWMutex

	UserColors    map[string]int
	VotingChannel string
	NewsChannel   string
	VoteCount     int
	PollCount     int
	ModRole       string    // role ID
	PlayChannels  Container // channelID
}

type ServerData struct {
	*sync.RWMutex

	LastCombs     map[string]Comb         // map[userID]comb
	PageSwitchers map[string]PageSwitcher // map[messageid]pageswitcher
	ComponentMsgs map[string]ComponentMsg // map[messageid]componentMsg
	ElementMsgs   map[string]int          // map[messageid]elemname
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

	// Don't need to set these
	Guild string
	Page  int
}

type Comb struct {
	Elems []int
	Elem3 int
}

type TimeStamp struct {
	time.Time
}

func (t *TimeStamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(t.Unix(), 10)), nil
}

func (t *TimeStamp) UnmarshalJSON(data []byte) error {
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return t.Time.UnmarshalJSON(data)
	}
	t.Time = time.Unix(i, 0)
	return nil
}

func NewTimeStamp(time time.Time) *TimeStamp {
	return &TimeStamp{
		Time: time,
	}
}

type Element struct {
	ID         int
	Name       string
	Image      string
	Color      int
	Guild      string
	Comment    string
	Creator    string
	CreatedOn  *TimeStamp
	Parents    []int
	Complexity int
	Difficulty int
	UsedIn     int
	TreeSize   int
}

type PollComboData struct {
	Elems  []int
	Result string
	Exists bool
}

type PollSignData struct {
	Elem    int
	NewNote string
	OldNote string
}

type PollImageData struct {
	Elem     int
	NewImage string
	OldImage string
}

type PollCategorizeData struct {
	Elems    []int
	Category string
}

type PollCatImageData struct {
	Category string
	NewImage string
	OldImage string
}

type PollColorData struct {
	Element int
	Color   int
}
type PollCatColorData struct {
	Category string
	Color    int
}

type Poll struct {
	Channel   string
	Message   string
	Guild     string
	Kind      PollType
	Suggestor string

	// Data, pointers to different types with omitempty so that you can selectively have some data
	PollComboData      *PollComboData      `json:"combodata,omitempty"`
	PollSignData       *PollSignData       `json:"signdata,omitempty"`
	PollImageData      *PollImageData      `json:"imagedata,omitempty"`
	PollCategorizeData *PollCategorizeData `json:"catdata,omitempty"` // This is also the uncategorize data
	PollCatImageData   *PollCatImageData   `json:"catimagedata,omitempty"`
	PollColorData      *PollColorData      `json:"colordata,omitempty"`
	PollCatColorData   *PollCatColorData   `json:"catcolordata,omitempty"`

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

type Msg struct {
	Author    *discordgo.User
	ChannelID string
	GuildID   string
}

type Inventory struct {
	Lock *sync.RWMutex `json:"-"`

	Elements map[int]Empty
	MadeCnt  int
	User     string
}

func (i *Inventory) Add(elem int) {
	i.Lock.Lock()
	i.Elements[elem] = Empty{}
	i.Lock.Unlock()
}

func (i *Inventory) Contains(elem int, nolock ...bool) bool {
	if len(nolock) == 0 {
		i.Lock.RLock()
		defer i.Lock.RUnlock()
	}
	_, exists := i.Elements[elem]
	return exists
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
		RWMutex: &sync.RWMutex{},

		UserColors:   make(map[string]int),
		PlayChannels: make(Container),
	}
}

func NewServerData() *ServerData {
	return &ServerData{
		RWMutex: &sync.RWMutex{},

		LastCombs:     make(map[string]Comb),
		PageSwitchers: make(map[string]PageSwitcher),
		ComponentMsgs: make(map[string]ComponentMsg),
		ElementMsgs:   make(map[string]int),
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

func NewInventory(user string, elements map[int]Empty, madecnt int) *Inventory {
	return &Inventory{
		Lock: &sync.RWMutex{},

		Elements: elements,
		User:     user,
		MadeCnt:  madecnt,
	}
}
