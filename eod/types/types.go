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
	PollKindCombo        PollKind = "combo"
	PollKindCategorize   PollKind = "cat"
	PollKindUncategorize PollKind = "rmcat"
	PollKindComment      PollKind = "comment" // Use "Sign", not "Comment"
	PollKindCatComment   PollKind = "catcomment"
	PollKindQueryComment PollKind = "querycomment"
	PollKindImage        PollKind = "img"
	PollKindCatImage     PollKind = "catimg"
	PollKindQueryImage   PollKind = "queryimg"
	PollKindColor        PollKind = "color"
	PollKindCatColor     PollKind = "catcolor"
	PollKindQueryColor   PollKind = "querycolor"
	PollKindQuery        PollKind = "query"
	PollKindDelQuery     PollKind = "delquery"
)

type PgData map[string]interface{}

func (p PgData) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p PgData) Scan(v interface{}) error {
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
	Data PgData   `db:"data"`
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

// Queries

type QueryKind string

const (
	QueryKindElement    QueryKind = "element"
	QueryKindCategory   QueryKind = "cat"
	QueryKindProducts   QueryKind = "products"
	QueryKindParents    QueryKind = "parents"
	QueryKindInventory  QueryKind = "inv"
	QueryKindElements   QueryKind = "elements"
	QueryKindRegex      QueryKind = "regex"
	QueryKindComparison QueryKind = "compare"
)

type Query struct {
	Guild   string `db:"guild"`
	Name    string `db:"name"`
	Creator string `db:"creator"`

	Image   string `db:"image"`
	Color   int    `db:"color"`
	Comment string `db:"comment"`

	Imager    string `db:"imager"`
	Colorer   string `db:"colorer"`
	Commenter string `db:"commenter"`

	CreatedOn time.Time `db:"createdon"`

	Kind QueryKind `db:"kind"`
	Data PgData    `db:"data"`

	// After calc
	Elements []int
}

// Categories

type Category struct {
	Guild   string `db:"guild"`
	Name    string `db:"name"`
	Comment string `db:"comment"`
	Image   string `db:"image"`
	Color   int    `db:"color"`

	Commenter string `db:"commenter"`
	Imager    string `db:"imager"`
	Colorer   string `db:"colorer"`

	Elements pq.Int32Array `db:"elements"`
}

// Consts
const RedCircle = "🔴"
const Check = "<:eodCheck:765333533362225222>"
const NoCheck = "❌"
