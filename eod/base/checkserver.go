package base

import (
	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

func (b *Base) CheckServer(m types.Msg, rsp types.Rsp) bool {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage("No voting or news channel has been set!")
		return false
	}
	if db.Config.VotingChannel == "" {
		rsp.ErrorMessage("No voting channel has been set!")
		return false
	}
	if db.Config.NewsChannel == "" {
		rsp.ErrorMessage("No news channel has been set!")
		return false
	}
	if len(db.Elements) < 4 {
		for _, elem := range types.StarterElements {
			err := db.SaveElement(elem, true)
			if rsp.Error(err) {
				return false
			}
		}
	}

	db.Config.RLock()
	_, exists := db.Config.PlayChannels[m.ChannelID]
	db.Config.RUnlock()
	return exists
}
