package types

import (
	"time"

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

	Parents pq.Int32Array `db:"parents"`
}
