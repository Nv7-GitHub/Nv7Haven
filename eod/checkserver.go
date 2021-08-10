package eod

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

var starterElements = []types.Element{
	{
		Name:       "Air",
		Comment:    "The invisible gaseous substance surrounding the earth, a mixture mainly of oxygen and nitrogen.",
		Creator:    "",
		Complexity: 0,
		Difficulty: 0,
		CreatedOn:  time.Unix(1, 0),
		Parents:    []string{},
	},
	{
		Name:       "Earth",
		Comment:    "The substance of the land surface; soil.",
		Creator:    "",
		Complexity: 0,
		Difficulty: 0,
		CreatedOn:  time.Unix(2, 0),
		Parents:    []string{},
	},
	{
		Name:       "Fire",
		Comment:    "Combustion or burning, in which substances combine chemically with oxygen from the air and typically give out bright light, heat, and smoke.",
		Creator:    "",
		Complexity: 0,
		Difficulty: 0,
		CreatedOn:  time.Unix(3, 0),
		Parents:    []string{},
	},
	{
		Name:       "Water",
		Comment:    "A colorless, transparent, odorless liquid that forms the seas, lakes, rivers, and rain and is the basis of the fluids of living organisms.",
		Creator:    "",
		Complexity: 0,
		Difficulty: 0,
		CreatedOn:  time.Unix(4, 0),
		Parents:    []string{},
	},
}

func (b *EoD) checkServer(m types.Msg, rsp types.Rsp) bool {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("No voting or news channel has been set!")
		return false
	}
	if dat.VotingChannel == "" {
		rsp.ErrorMessage("No voting channel has been set!")
		return false
	}
	if dat.NewsChannel == "" {
		rsp.ErrorMessage("No news channel has been set!")
		return false
	}
	if dat.ElemCache == nil {
		dat.ElemCache = make(map[string]types.Element)
	}
	if len(dat.ElemCache) < 4 {
		for _, elem := range starterElements {
			elem.Guild = m.GuildID
			elem.CreatedOn = time.Now()
			dat.Lock.Lock()
			dat.ElemCache[strings.ToLower(elem.Name)] = elem
			dat.Lock.Unlock()
			_, err := b.db.Exec("INSERT INTO eod_elements VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )", elem.Name, elem.Image, elem.Guild, elem.Comment, elem.Creator, int(elem.CreatedOn.Unix()), "" /* Parents */, elem.Complexity, elem.Difficulty, elem.UsedIn)
			rsp.Error(err)
		}
		lock.Lock()
		b.dat[m.GuildID] = dat
		lock.Unlock()
	}

	if dat.InvCache == nil {
		dat.InvCache = make(map[string]map[string]types.Empty)
	}
	dat.Lock.RLock()
	_, exists = dat.InvCache[m.Author.ID]
	dat.Lock.RUnlock()
	if !exists {
		dat.Lock.Lock()
		dat.InvCache[m.Author.ID] = make(map[string]types.Empty)
		for _, val := range starterElements {
			dat.InvCache[m.Author.ID][strings.ToLower(val.Name)] = types.Empty{}
		}
		dat.Lock.Unlock()

		dat.Lock.RLock()
		inv := dat.InvCache[m.Author.ID]
		dat.Lock.RUnlock()

		data, err := json.Marshal(inv)
		if rsp.Error(err) {
			return false
		}
		_, err = b.db.Exec("INSERT INTO eod_inv VALUES ( ?, ?, ?, ?, ? )", m.GuildID, m.Author.ID, string(data), len(inv), 0) // Guild ID, User ID, inventory, elements found, made by (0 so far)
		rsp.Error(err)
		lock.Lock()
		b.dat[m.GuildID] = dat
		lock.Unlock()
	}

	if dat.CatCache == nil {
		dat.CatCache = make(map[string]types.Category)
	}

	_, exists = dat.PlayChannels[m.ChannelID]
	return exists
}
