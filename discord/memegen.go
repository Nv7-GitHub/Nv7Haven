package discord

import (
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/golang/freetype"
)

const size = 10
const fontFile = "discord/memes/Arial.ttf"
const dpi = 100

func drawMeme(fileName, text string, x, y int) (image.Image, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	meme, err := png.Decode(file)
	if err != nil {
		return nil, err
	}
	file.Close()

	c := freetype.NewContext()

	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		return nil, err
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	pt := freetype.Pt(x, y+int(c.PointToFixed(size)>>6))

	c.SetDPI(dpi)
	c.SetFont(f)
	c.SetFontSize(size)
	c.SetClip(meme.Bounds())
	c.SetDst(meme.(*image.RGBA))
	c.SetSrc(image.Black)

	_, err = c.DrawString(text, pt)
	if err != nil {
		return nil, err
	}
	return meme, nil
}

type memeGen struct {
	File string
	X    int
	Y    int
}

var memes = map[string]memeGen{
	"stroke": memeGen{
		File: "discord/memes/stroke.png",
		X:    355,
		Y:    380,
	},
}

func (b *Bot) memeGen(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "memelist") {
		var text string
		for k := range memes {
			text += k + "\n"
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title:       "Nv7 Bot's Memes",
			Description: text,
		})
	}
}
