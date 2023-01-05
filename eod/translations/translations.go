package translations

import (
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/jmoiron/sqlx"
)

type Translations struct {
	db   *sqlx.DB
	base *base.Base
	s    *sevcord.Sevcord
}

func (t *Translations) AutoComplete(ctx sevcord.Ctx, val any) []sevcord.Choice {
	res := []string{"English", "Cresorium"}

	choices := make([]sevcord.Choice, len(res))
	for i := range res {
		choices[i] = sevcord.NewChoice(res[i], res[i])
	}
	return choices
}

func (t *Translations) Init() {
	t.s.RegisterSlashCommand(sevcord.NewSlashCommand(
		"translate",
		"Set the personal language of the bot!",
		t.SetTranslate,
		sevcord.NewOption("language", "The language to translate to!", sevcord.OptionKindString, true).
			AutoComplete(t.AutoComplete),
	))
}

func NewTranslations(d *sqlx.DB, b *base.Base, s *sevcord.Sevcord) *Translations {
	t := &Translations{
		db:   d,
		base: b,
		s:    s,
	}
	t.Init()
	return t
}
