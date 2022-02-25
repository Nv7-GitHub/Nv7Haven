package types

import "time"

const X = "❌"
const Check = "<:eodCheck:765333533362225222>" // ✅
const RedCircle = "🔴"
const NewText = "🆕"
const UpArrow = "⬆️"
const DownArrow = "⬇️"

var MaxComboLength = 21

const MaxTries = 20
const AutocompleteResults = 25

var StarterElements = []Element{
	{
		Name:       "Air",
		ID:         1,
		Comment:    "The invisible gaseous substance surrounding the earth, a mixture mainly of oxygen and nitrogen.",
		Image:      "https://cdn.discordapp.com/attachments/819077689775882252/819974778106282054/air.png",
		Color:      12764099, // #C2C3C3
		Creator:    "819076922867712031",
		Complexity: 0,
		Difficulty: 0,
		CreatedOn:  NewTimeStamp(time.Unix(1, 0)),
		Parents:    []int{},
		Air:        1,
		Earth:      0,
		Fire:       0,
		Water:      0,
	},
	{
		Name:       "Earth",
		ID:         2,
		Comment:    "The substance of the land surface; soil.",
		Image:      "https://cdn.discordapp.com/attachments/819078122963861525/820507498737172490/Earth-Science-Facts-for-Kids-All-About-Whats-in-Soil-Image-of-Soil.png",
		Color:      11172162, // #AA7942
		Creator:    "819076922867712031",
		Complexity: 0,
		Difficulty: 0,
		CreatedOn:  NewTimeStamp(time.Unix(2, 0)),
		Parents:    []int{},
		Air:        0,
		Earth:      1,
		Fire:       0,
		Water:      0,
	},
	{
		Name:       "Fire",
		ID:         3,
		Comment:    "Combustion or burning, in which substances combine chemically with oxygen from the air and typically give out bright light, heat, and smoke.",
		Image:      "https://cdn.discordapp.com/attachments/819078122963861525/820508007795916820/fire-flame-flames-heat-burn-hot-blaze-fiery-burning.png",
		Color:      16749824, // #FF9500
		Creator:    "819076922867712031",
		Complexity: 0,
		Difficulty: 0,
		CreatedOn:  NewTimeStamp(time.Unix(3, 0)),
		Parents:    []int{},
		Air:        0,
		Earth:      0,
		Fire:       1,
		Water:      0,
	},
	{
		Name:       "Water",
		ID:         4,
		Comment:    "A colorless, transparent, odorless liquid that forms the seas, lakes, rivers, and rain and is the basis of the fluids of living organisms.",
		Image:      "https://cdn.discordapp.com/attachments/819078122963861525/820513012074151947/water.png",
		Color:      275455, // #0433FF
		Creator:    "819076922867712031",
		Complexity: 0,
		Difficulty: 0,
		CreatedOn:  NewTimeStamp(time.Unix(4, 0)),
		Parents:    []int{},
		Air:        0,
		Earth:      0,
		Fire:       0,
		Water:      1,
	},
}
