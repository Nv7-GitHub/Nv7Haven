package eod

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/Nv7-Github/Nv7Haven/eod/achievements"
	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/categories"
	"github.com/Nv7-Github/Nv7Haven/eod/elements"
	"github.com/Nv7-Github/Nv7Haven/eod/pages"
	"github.com/Nv7-Github/Nv7Haven/eod/polls"
	"github.com/Nv7-Github/Nv7Haven/eod/queries"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/jmoiron/sqlx"
)

type Bot struct {
	s  *sevcord.Sevcord
	db *sqlx.DB

	// Modules
	base       *base.Base
	elements   *elements.Elements
	polls      *polls.Polls
	pages      *pages.Pages
	categories *categories.Categories
	queries    *queries.Queries
	users      *achievements.Users
}

func InitEod(db *sqlx.DB, token string) {
	s, err := sevcord.New(token)
	if err != nil {
		panic(err)
	}
	b := Bot{
		s:  s,
		db: db,
	}
	b.Init()

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)
		<-stop
		fmt.Println("Saving command stats...")
		b.base.SaveCommandStats("", nil)
	}()
	b.s.Listen()
}
