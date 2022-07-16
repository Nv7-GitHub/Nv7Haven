package base

import (
	"compress/gzip"

	"github.com/bwmarrin/discordgo"
)

const maxLen = 8000000

func (b *Base) PrepareFile(d *discordgo.File, len int) *discordgo.File {
	if len > maxLen {
		d.ContentType = "application/gzip"
		d.Name += ".gz"
		var err error
		d.Reader, err = gzip.NewReader(d.Reader)
		if err != nil {
			panic(err)
		}
	}
	return d
}
