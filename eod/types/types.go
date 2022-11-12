package types

import (
	"time"

	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
)

type Resp struct {
	Ok      bool
	Message string
}

func Ok() Resp                 { return Resp{Ok: true} }
func Fail(message string) Resp { return Resp{Ok: false, Message: message} }

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

type Config struct {
	Guild         string         `db:"guild"`
	VotingChannel string         `db:"voting"`
	NewsChannel   string         `db:"news"`
	VoteCnt       int            `db:"votecnt"`
	PollCnt       int            `db:"pollcnt"`
	PlayChannels  pq.StringArray `db:"play"`
}

var Sorts = []sevcord.Choice{
	sevcord.NewChoice("ID", "id"),
	sevcord.NewChoice("Name", "name"),
	sevcord.NewChoice("Creator", "creator"),
	sevcord.NewChoice("Created On", "createdon"),
	sevcord.NewChoice("Tree Size", "treesize"),
}

// Creates btns with name prevnext|<params>
func PageSwitchBtns(handler, params string) []sevcord.Component {
	return []sevcord.Component{
		sevcord.NewButton("", sevcord.ButtonStylePrimary, handler, "prev|"+params).WithEmoji(sevcord.ComponentEmojiCustom("leftarrow", "861722690813165598", false)),
		sevcord.NewButton("", sevcord.ButtonStylePrimary, handler, "next|"+params).WithEmoji(sevcord.ComponentEmojiCustom("rightarrow", "861722690926936084", false)),
	}
}
