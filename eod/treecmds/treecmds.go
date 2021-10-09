package treecmds

import (
	"sync"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

type TreeCmds struct {
	lock *sync.RWMutex
	dat  map[string]types.ServerData
	base *base.Base
	dg   *discordgo.Session
}

func NewTreeCmds(dat map[string]types.ServerData, dg *discordgo.Session, base *base.Base, lock *sync.RWMutex) *TreeCmds {
	return &TreeCmds{
		lock: lock,
		dat:  dat,
		base: base,
		dg:   dg,
	}
}
