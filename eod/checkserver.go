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
	if dat.Elements == nil {
		dat.Elements = make(map[string]types.Element)
	}
	if len(dat.Elements) < 4 {
		for _, elem := range starterElements {
			elem.Guild = m.GuildID
			dat.SetElement(elem)
			_, err := b.db.Exec("INSERT INTO eod_elements VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )", elem.Name, elem.Image, elem.Guild, elem.Comment, elem.Creator, int(elem.CreatedOn.Unix()), "" /* Parents */, elem.Complexity, elem.Difficulty, elem.UsedIn)
			rsp.Error(err)
		}
		lock.Lock()
		b.dat[m.GuildID] = dat
		lock.Unlock()
	}

	if dat.Inventories == nil {
		dat.Inventories = make(map[string]types.Container)
	}
	_, res := dat.GetInv(m.Author.ID, true)
	if !res.Exists {
		dat.Lock.Lock()
		dat.Inventories[m.Author.ID] = make(map[string]types.Empty)
		for _, val := range starterElements {
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
		lock.Lock()
		b.dat[m.GuildID] = dat
		lock.Unlock()
	}

	if dat.Categories == nil {
		dat.Lock.Lock()
		dat.Categories = make(map[string]types.Category)
		dat.Lock.Unlock()

		lock.Lock()
		b.dat[m.GuildID] = dat
		lock.Unlock()
	}

	if dat.LastCombs == nil {
		dat.Lock.Lock()
		dat.LastCombs = make(map[string]types.Comb)
		dat.Lock.Unlock()

		lock.Lock()
		b.dat[m.GuildID] = dat
		lock.Unlock()
	}

	if dat.Combos == nil {
		dat.Lock.Lock()
		dat.Combos = make(map[string]string)
		dat.Lock.Unlock()

		lock.Lock()
		b.dat[m.GuildID] = dat
		lock.Unlock()
	}

	if dat.ElementMsgs == nil {
		dat.Lock.Lock()
		dat.ElementMsgs = make(map[string]string)
		dat.Lock.Unlock()

		lock.Lock()
		b.dat[m.GuildID] = dat
		lock.Unlock()
	}

	_, exists = dat.PlayChannels[m.ChannelID]
	return exists
}
