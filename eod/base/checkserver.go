package base

import (
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

var StarterElements = []types.Element{
	{
		Name:       "Air",
		ID:         1,
		Comment:    "The invisible gaseous substance surrounding the earth, a mixture mainly of oxygen and nitrogen.",
		Image:      "https://cdn.discordapp.com/attachments/819077689775882252/819974778106282054/air.png",
		Color:      12764099, // #C2C3C3
		Creator:    "",
		Complexity: 0,
		Difficulty: 0,
		CreatedOn:  types.NewTimeStamp(time.Unix(1, 0)),
		Parents:    []int{},
	},
	{
		Name:       "Earth",
		ID:         2,
		Comment:    "The substance of the land surface; soil.",
		Image:      "https://cdn.discordapp.com/attachments/819078122963861525/820507498737172490/Earth-Science-Facts-for-Kids-All-About-Whats-in-Soil-Image-of-Soil.png",
		Color:      11172162, // #AA7942
		Creator:    "",
		Complexity: 0,
		Difficulty: 0,
		CreatedOn:  types.NewTimeStamp(time.Unix(2, 0)),
		Parents:    []int{},
	},
	{
		Name:       "Fire",
		ID:         3,
		Comment:    "Combustion or burning, in which substances combine chemically with oxygen from the air and typically give out bright light, heat, and smoke.",
		Image:      "https://cdn.discordapp.com/attachments/819078122963861525/820508007795916820/fire-flame-flames-heat-burn-hot-blaze-fiery-burning.png",
		Color:      16749824, // #FF9500
		Creator:    "",
		Complexity: 0,
		Difficulty: 0,
		CreatedOn:  types.NewTimeStamp(time.Unix(3, 0)),
		Parents:    []int{},
	},
	{
		Name:       "Water",
		ID:         4,
		Comment:    "A colorless, transparent, odorless liquid that forms the seas, lakes, rivers, and rain and is the basis of the fluids of living organisms.",
		Image:      "https://cdn.discordapp.com/attachments/819078122963861525/820513012074151947/water.png",
		Color:      275455, // #0433FF
		Creator:    "",
		Complexity: 0,
		Difficulty: 0,
		CreatedOn:  types.NewTimeStamp(time.Unix(4, 0)),
		Parents:    []int{},
	},
}

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
		for _, elem := range StarterElements {
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
