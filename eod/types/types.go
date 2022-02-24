package types

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/translation"
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
	PollDeleteVCat   = 8

	PageSwitchLdb = 0
	PageSwitchInv = 1
)

type ComponentMsg interface {
	Handler(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type ModalHandler interface {
	Handler(s *discordgo.Session, i *discordgo.InteractionCreate, rsp Rsp)
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
	LanguageFile  string
}

type ServerData struct {
	*sync.RWMutex

	LastCombs     map[string]Comb         // map[userID]comb
	PageSwitchers map[string]PageSwitcher // map[messageid]pageswitcher
	ComponentMsgs map[string]ComponentMsg // map[messageid]componentMsg
	ElementMsgs   map[string]int          // map[messageid]elemname
	Modals        map[string]ModalHandler // map[interactionid]modalHandler, NOTE: interactionid is CustomID
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

	Commenter string
	Colorer   string
	Imager    string
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
	Changed  bool
}

type PollCategorizeData struct {
	Elems    []int
	Category string
	Title    string
}

type PollVCatDeleteData struct {
	Category string
}

type PollCatImageData struct {
	Category string
	NewImage string
	OldImage string
	Changed  bool
}

type PollColorData struct {
	Element  int
	Color    int
	OldColor int
}
type PollCatColorData struct {
	Category string
	Color    int
	OldColor int
}

type Poll struct {
	Channel   string
	Message   string
	Guild     string
	Kind      PollType
	Suggestor string
	CreatedOn *TimeStamp

	// Data, pointers to different types with omitempty so that you can selectively have some data
	PollComboData      *PollComboData      `json:"combodata,omitempty"`
	PollSignData       *PollSignData       `json:"signdata,omitempty"`
	PollImageData      *PollImageData      `json:"imagedata,omitempty"`
	PollCategorizeData *PollCategorizeData `json:"catdata,omitempty"` // This is also the uncategorize data
	PollCatImageData   *PollCatImageData   `json:"catimagedata,omitempty"`
	PollColorData      *PollColorData      `json:"colordata,omitempty"`
	PollCatColorData   *PollCatColorData   `json:"catcolordata,omitempty"`
	PollVCatDeleteData *PollVCatDeleteData `json:"vcatdeldata,omitempty"` // This is also the uncategorize data

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

	Imager  string
	Colorer string
}

type VirtualCategoryRuleType int

const (
	VirtualCategoryRuleRegex        VirtualCategoryRuleType = 0
	VirtualCategoryRuleInvFilter    VirtualCategoryRuleType = 1
	VirtualCategoryRuleSetOperation VirtualCategoryRuleType = 2
	VirtualCategoryRuleAllElements  VirtualCategoryRuleType = 3
)

type VirtualCategoryData map[string]interface{}

type CategoryOperation string

const (
	CatOpUnion     CategoryOperation = "union"
	CatOpIntersect CategoryOperation = "intersect"
	CatOpDiff      CategoryOperation = "difference"
)

type VirtualCategory struct {
	Name    string
	Guild   string
	Creator string

	Image   string
	Color   int
	Imager  string
	Colorer string

	Rule VirtualCategoryRuleType
	Data VirtualCategoryData

	Cache map[int]Empty `json:"-"`
}

type Msg struct {
	Author    *discordgo.User
	ChannelID string
	GuildID   string
}

type Inventory struct {
	Lock *sync.RWMutex `json:"-"`

	Elements      map[int]Empty
	MadeCnt       int
	ImagedCnt     int
	SignedCnt     int
	ColoredCnt    int
	CatImagedCnt  int
	CatColoredCnt int
	UsedCnt       int
	User          string
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
	RawEmbed(emb *discordgo.MessageEmbed, components ...discordgo.MessageComponent) string
	Resp(msg string, components ...discordgo.MessageComponent)
	Acknowledge()
	DM(msg string)
	Attachment(text string, files []*discordgo.File)
	Modal(modal *discordgo.InteractionResponseData, handler ModalHandler)
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		RWMutex: &sync.RWMutex{},

		UserColors:   make(map[string]int),
		PlayChannels: make(Container),
		LanguageFile: translation.DefaultLang,
	}
}

func NewServerData() *ServerData {
	return &ServerData{
		RWMutex: &sync.RWMutex{},

		LastCombs:     make(map[string]Comb),
		PageSwitchers: make(map[string]PageSwitcher),
		ComponentMsgs: make(map[string]ComponentMsg),
		ElementMsgs:   make(map[string]int),
		Modals:        make(map[string]ModalHandler),
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
