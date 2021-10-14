package base

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

var StarterElements = []types.Element{
	{
		Name:       "Air",
		Comment:    "The invisible gaseous substance surrounding the earth, a mixture mainly of oxygen and nitrogen.",
		Image:      "https://cdn.discordapp.com/attachments/819077689775882252/819974778106282054/air.png",
		Color:      12764099, // #C2C3C3
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
		Color:      11172162, // #AA7942
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
		Color:      16749824, // #FF9500
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
		Color:      275455, // #0433FF
		Creator:    "",
		Complexity: 0,
		Difficulty: 0,
		CreatedOn:  time.Unix(4, 0),
		Parents:    []string{},
	},
}

func (b *Base) CheckServer(m types.Msg, rsp types.Rsp) bool {
	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
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
	if len(dat.Elements) < 4 {
		for _, elem := range StarterElements {
			elem.Guild = m.GuildID
			dat.SetElement(elem)
			_, err := b.db.Exec("INSERT INTO eod_elements VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )", elem.Name, elem.Image, elem.Color, elem.Guild, elem.Comment, elem.Creator, int(elem.CreatedOn.Unix()), "" /* Parents */, elem.Complexity, elem.Difficulty, elem.UsedIn, elem.TreeSize)
			rsp.Error(err)
		}
		b.lock.Lock()
		b.dat[m.GuildID] = dat
		b.lock.Unlock()
	}

	_, res := dat.GetInv(m.Author.ID, true)
	if !res.Exists {
		dat.Lock.Lock()
		dat.Inventories[m.Author.ID] = make(map[string]types.Empty)
		for _, val := range StarterElements {
			dat.Inventories[m.Author.ID][strings.ToLower(val.Name)] = types.Empty{}
		}
		dat.Lock.Unlock()

		inv, _ := dat.GetInv(m.Author.ID, true)

		data, err := json.Marshal(inv)
		if rsp.Error(err) {
			return false
		}
		_, err = b.db.Exec("INSERT INTO eod_inv VALUES ( ?, ?, ?, ?, ? )", m.GuildID, m.Author.ID, string(data), len(inv), 0) // Guild ID, User ID, inventory, elements found, made by (0 so far)
		rsp.Error(err)
		b.lock.Lock()
		b.dat[m.GuildID] = dat
		b.lock.Unlock()
	}

	_, exists = dat.PlayChannels[m.ChannelID]
	return exists
}
