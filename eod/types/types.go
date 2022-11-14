package types

import (
	"database/sql/driver"
	"encoding/json"
	"sync"
	"time"

	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
)

// Resp
type Resp struct {
	Ok      bool
	Message string
}

func Ok() Resp                 { return Resp{Ok: true} }
func Fail(message string) Resp { return Resp{Ok: false, Message: message} }

// Element
type Element struct {
	ID        int       `db:"id"`
	Guild     string    `db:"guild"`
	Name      string    `db:"name"`
	Image     string    `db:"image"`
	Color     int       `db:"color"`
	Comment   string    `db:"comment"`
	Creator   string    `db:"creator"`
	CreatedOn time.Time `db:"createdon"`

	Commenter string `db:"commenter"`
	Colorer   string `db:"colorer"`
	Imager    string `db:"imager"`

	Parents  pq.Int32Array `db:"parents"`
	TreeSize int           `db:"treesize"`
}

// Guilds
type Config struct {
	Guild         string         `db:"guild"`
	VotingChannel string         `db:"voting"`
	NewsChannel   string         `db:"news"`
	VoteCnt       int            `db:"votecnt"`
	PollCnt       int            `db:"pollcnt"`
	PlayChannels  pq.StringArray `db:"play"`
}

type ServerMem struct {
	sync.RWMutex
	CombCache map[string][]int // map[userid][]id
}

// Polls
type PollKind string

const (
	PollKindCombo         PollKind = "combo"
	PollKindCategorize    PollKind = "cat"
	PollKindUncategorized PollKind = "rmcat"
	PollKindComment       PollKind = "comment" // Use "Sign", not "Comment"
	PollKindCatComment    PollKind = "catcomment"
	PollKindImage         PollKind = "img"
	PollKindCatImage      PollKind = "catimg"
	PollKindColor         PollKind = "color"
	PollKindCatColor      PollKind = "catcolor"
)

type PollData map[string]interface{}

func (p PollData) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p PollData) Scan(v interface{}) error {
	return json.Unmarshal(v.([]byte), &p)
}

type Poll struct {
	// Filed in by CreatePoll
	Guild     string    `db:"guild"`
	Message   string    `db:"message"`
	Channel   string    `db:"channel"`
	Creator   string    `db:"creator"`
	CreatedOn time.Time `db:"createdon"`

	Upvotes   int `db:"upvotes"`
	Downvotes int `db:"downvotes"`

	// Required
	Kind PollKind `db:"kind"`
	Data PollData `db:"data"`
}

// Discord util
var Sorts = []sevcord.Choice{
	sevcord.NewChoice("ID", "id"),
	sevcord.NewChoice("Name", "name"),
	sevcord.NewChoice("Creator", "creator"),
	sevcord.NewChoice("Created On", "createdon"),
	sevcord.NewChoice("Tree Size", "treesize"),
}

var SortSql = map[string]string{
	"id":        "id",
	"name":      "name",
	"creator":   "creator",
	"createdon": "createdon",
	"treesize":  "treesize DESC",
}

// Consts
const RedCircle = "🔴"
const Check = "<:eodCheck:765333533362225222>"
