package types

import (
	"time"

	"github.com/lib/pq"
)

func Starters(guild string) []Element {
	return []Element{
		{
			Guild: guild,

			Name:      "Air",
			ID:        1,
			Comment:   "The invisible gaseous substance surrounding the earth, a mixture mainly of oxygen and nitrogen.",
			Image:     "https://cdn.discordapp.com/attachments/819077689775882252/819974778106282054/air.png",
			Color:     12764099, // #C2C3C3
			Creator:   "819076922867712031",
			CreatedOn: time.Unix(1, 0),
			Parents:   pq.Int32Array{},

			Commenter: "819076922867712031",
			Colorer:   "819076922867712031",
			Imager:    "819076922867712031",
			TreeSize:  1,
		},
		{
			Guild: guild,

			Name:      "Earth",
			ID:        2,
			Comment:   "The substance of the land surface; soil.",
			Image:     "https://cdn.discordapp.com/attachments/819078122963861525/820507498737172490/Earth-Science-Facts-for-Kids-All-About-Whats-in-Soil-Image-of-Soil.png",
			Color:     11172162, // #AA7942
			Creator:   "819076922867712031",
			CreatedOn: time.Unix(2, 0),
			Parents:   pq.Int32Array{},

			Commenter: "819076922867712031",
			Colorer:   "819076922867712031",
			Imager:    "819076922867712031",
			TreeSize:  1,
		},
		{
			Guild: guild,

			Name:      "Fire",
			ID:        3,
			Comment:   "Combustion or burning, in which substances combine chemically with oxygen from the air and typically give out bright light, heat, and smoke.",
			Image:     "https://cdn.discordapp.com/attachments/819078122963861525/820508007795916820/fire-flame-flames-heat-burn-hot-blaze-fiery-burning.png",
			Color:     16749824, // #FF9500
			Creator:   "819076922867712031",
			CreatedOn: time.Unix(3, 0),
			Parents:   pq.Int32Array{},

			Commenter: "819076922867712031",
			Colorer:   "819076922867712031",
			Imager:    "819076922867712031",
			TreeSize:  1,
		},
		{
			Guild: guild,

			Name:      "Water",
			ID:        4,
			Comment:   "A colorless, transparent, odorless liquid that forms the seas, lakes, rivers, and rain and is the basis of the fluids of living organisms.",
			Image:     "https://cdn.discordapp.com/attachments/819078122963861525/820513012074151947/water.png",
			Color:     275455, // #0433FF
			Creator:   "819076922867712031",
			CreatedOn: time.Unix(4, 0),
			Parents:   pq.Int32Array{},

			Commenter: "819076922867712031",
			Colorer:   "819076922867712031",
			Imager:    "819076922867712031",
			TreeSize:  1,
		},
	}
}
