package categories

import (
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/polls"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

type Categories struct {
	dat   map[string]types.ServerData
	lock  *sync.RWMutex
	base  *base.Base
	dg    *discordgo.Session
	polls *polls.Polls
}

func NewCategories(dat map[string]types.ServerData, base *base.Base, dg *discordgo.Session, polls *polls.Polls, lock *sync.RWMutex) *Categories {
	return &Categories{
		dat:   dat,
		lock:  lock,
		base:  base,
		dg:    dg,
		polls: polls,
	}
}
