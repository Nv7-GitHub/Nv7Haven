package discord

import (
	"bytes"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/golang/freetype"
)

const fontFile = "discord/memes/Arial.ttf"
const dpi = 100

var reg = regexp.MustCompile(`genmeme ([A-Za-z1-9]+) (.+)`)

func drawMeme(fileName, text string, x, y int, size float64) (image.Image, error) {
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
	Size float64
}

var memes = map[string]memeGen{
	"stroke": memeGen{
		File: "discord/memes/stroke.png",
		X:    355,
		Y:    380,
		Size: 12,
	},
	"violence": memeGen{
		File: "discord/memes/violence.png",
		X:    700,
		Y:    10,
		Size: 24,
	},
	"humanity": memeGen{
		File: "discord/memes/humanity.png",
		X:    630,
		Y:    600,
		Size: 18,
	},
	"abandon": memeGen{
		File: "discord/memes/abandon.png",
		X:    50,
		Y:    285,
		Size: 10,
	},
}

func (b *Bot) memeGen(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "listmemes") {
		var text string
		for k := range memes {
			text += k + "\n"
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title:       "Nv7 Bot's Memes",
			Description: text,
		})
		return
	}

	if strings.HasPrefix(m.Content, "genmeme") {
		match := reg.FindAllStringSubmatch(m.Content, -1)
		if (len(match) == 0) || (len(match[0]) < 3) {
			s.ChannelMessageSend(m.ChannelID, "Does not fit format `genmeme <name> <text>`")
			return
		}

		name := match[0][1]
		text := match[0][2]

		template, exists := memes[name]
		if !exists {
			s.ChannelMessageSend(m.ChannelID, "That meme template doesn't exist!")
			return
		}

		out, err := drawMeme(template.File, text, template.X, template.Y, template.Size)
		if b.handle(err, m) {
			return
		}

		final := bytes.NewBuffer([]byte{})
		err = png.Encode(final, out)
		if b.handle(err, m) {
			return
		}

		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Files: []*discordgo.File{
				&discordgo.File{
					Name:        "meme.png",
					ContentType: "image/png",
					Reader:      final,
				},
			},
		})
		return
	}
}
