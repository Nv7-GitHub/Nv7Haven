package elements

import (
	"bytes"
	"log"
	"sort"
	"strconv"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

const width = 1800
const spacing = 20

var face font.Face

func init() {
	fnt, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
	face = truetype.NewFace(fnt, &truetype.Options{Size: 48})
}

func (b *Elements) LbImageCmd(m types.Msg, rsp types.Rsp, sorter string) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	rsp.Acknowledge()

	// Calculate number of users to display
	num := 10
	heightSpacing := 100
	db.Config.RLock()
	_, exists := db.Config.PlayChannels[m.ChannelID]
	db.Config.RUnlock()
	if exists {
		num = 30
		heightSpacing = 200
	}

	// Get leaderboard
	invs := make([]*types.Inventory, len(db.Invs()))
	i := 0
	for _, v := range db.Invs() {
		invs[i] = v
		i++
	}
	sortFn := func(a, b int) bool {
		return len(invs[a].Elements) > len(invs[b].Elements)
	}
	if sorter == "made" {
		sortFn = func(a, b int) bool {
			return invs[a].MadeCnt > invs[b].MadeCnt
		}
	}
	sort.Slice(invs, sortFn)
	if len(invs) > num {
		invs = invs[:num]
	}

	// Max vals
	maxElems := 0
	maxMade := 0
	for _, inv := range invs {
		if len(inv.Elements) > maxElems {
			maxElems = len(inv.Elements)
		}

		if inv.MadeCnt > maxMade {
			maxMade = inv.MadeCnt
		}
	}

	// Process
	height := 100 * num

	dc := gg.NewContext(width, height)
	dc.SetFontFace(face)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Draw invs
	for i, inv := range invs {
		// Draw position
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored("#"+strconv.Itoa(i+1), float64(width/2), float64(i*height/len(invs)), 0.5, 1)

		// Draw bars
		space := width/2 - (width / (spacing / 2))

		// Made
		size := (float64(len(inv.Elements)) / float64(maxElems)) * float64(space)
		dc.DrawRectangle(width/2+(width/spacing), float64(i*height/len(invs))+float64(height)/float64(heightSpacing), size, float64(height/len(invs))-float64(height/(heightSpacing/2)))
		dc.SetRGB(0, 1, 0)
		dc.Fill()

		// Suggested
		size = (float64(inv.MadeCnt) / float64(maxMade)) * float64(space)
		dc.DrawRectangle(width/2-width/spacing-size, float64(i*height/len(invs))+float64(height)/float64(heightSpacing), size, float64(height/len(invs))-float64(height/(heightSpacing/2)))
		dc.SetRGB(0, 1, 1)
		dc.Fill()
	}

	// Save
	buf := bytes.NewBuffer(nil)
	err := dc.EncodePNG(buf)
	if rsp.Error(err) {
		return
	}

	rsp.Attachment("**Leaderboard Image**", []*discordgo.File{
		{
			Name:        "image.png",
			Reader:      buf,
			ContentType: "image/png",
		},
	})
}
