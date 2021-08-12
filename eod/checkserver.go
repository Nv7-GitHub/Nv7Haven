package eod

import (
	"encoding/json"
	"strings"
	"time"
)

var starterElements = []element{
	{
		Name:       "Air",
		Comment:    "The invisible gaseous substance surrounding the earth, a mixture mainly of oxygen and nitrogen.",
		Image:      "https://cdn.discordapp.com/attachments/819077689775882252/819974778106282054/air.png",
		Creator:    "",
		Complexity: 0,
		Difficulty: 0,
		CreatedOn:  time.Unix(1, 0),
		Parents:    []string{},
	},
	{
		Name:       "Earth",
		Comment:    "The substance of the land surface; soil.",
		Image:      "https://cdn.discordapp.com/attachments/819078122963861525/820507498737172490/Earth-Science-Facts-for-Kids-All-About-Whats-in-Soil-Image-of-Soil.png",
		Creator:    "",
		Complexity: 0,
		Difficulty: 0,
		CreatedOn:  time.Unix(2, 0),
		Parents:    []string{},
	},
	{
		Name:       "Fire",
		Comment:    "Combustion or burning, in which substances combine chemically with oxygen from the air and typically give out bright light, heat, and smoke.",
		Image:      "https://cdn.discordapp.com/attachments/819078122963861525/820508007795916820/fire-flame-flames-heat-burn-hot-blaze-fiery-burning.png",
		Creator:    "",
		Complexity: 0,
		Difficulty: 0,
		CreatedOn:  time.Unix(3, 0),
		Parents:    []string{},
	},
	{
		Name:       "Water",
		Comment:    "A colorless, transparent, odorless liquid that forms the seas, lakes, rivers, and rain and is the basis of the fluids of living organisms.",
		Image:      "https://cdn.discordapp.com/attachments/819078122963861525/820513012074151947/water.png",
		Creator:    "",
		Complexity: 0,
		Difficulty: 0,
		CreatedOn:  time.Unix(4, 0),
		Parents:    []string{},
	},
}

func (b *EoD) checkServer(m msg, rsp rsp) bool {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("No voting or news channel has been set!")
		return false
	}
	if dat.votingChannel == "" {
		rsp.ErrorMessage("No voting channel has been set!")
		return false
	}
	if dat.newsChannel == "" {
		rsp.ErrorMessage("No news channel has been set!")
		return false
	}
	if dat.elemCache == nil {
		dat.elemCache = make(map[string]element)
	}
	if len(dat.elemCache) < 4 {
		for _, elem := range starterElements {
			elem.Guild = m.GuildID
			elem.CreatedOn = time.Now()
			dat.lock.Lock()
			dat.elemCache[strings.ToLower(elem.Name)] = elem
			dat.lock.Unlock()
			_, err := b.db.Exec("INSERT INTO eod_elements VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )", elem.Name, elem.Image, elem.Guild, elem.Comment, elem.Creator, int(elem.CreatedOn.Unix()), "" /* Parents */, elem.Complexity, elem.Difficulty, elem.UsedIn)
			rsp.Error(err)
		}
		lock.Lock()
		b.dat[m.GuildID] = dat
		lock.Unlock()
	}

	if dat.invCache == nil {
		dat.invCache = make(map[string]map[string]empty)
	}
	dat.lock.RLock()
	_, exists = dat.invCache[m.Author.ID]
	dat.lock.RUnlock()
	if !exists {
		dat.lock.Lock()
		dat.invCache[m.Author.ID] = make(map[string]empty)
		for _, val := range starterElements {
			dat.invCache[m.Author.ID][strings.ToLower(val.Name)] = empty{}
		}
		dat.lock.Unlock()

		dat.lock.RLock()
		inv := dat.invCache[m.Author.ID]
		dat.lock.RUnlock()

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

	if dat.catCache == nil {
		dat.catCache = make(map[string]category)
	}

	_, exists = dat.playChannels[m.ChannelID]
	return exists
}
