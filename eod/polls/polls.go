package polls

import (
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/jmoiron/sqlx"
)

type Polls struct {
	d *sqlx.DB
	b *base.Base
}

func NewPolls(d *sqlx.DB, b *base.Base, s *sevcord.Sevcord) *Polls {
	return &Polls{
		d: d,
		b: b,
	}
}
