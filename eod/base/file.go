package base

import (
	"bytes"
	"compress/gzip"
	"io"

	"github.com/bwmarrin/discordgo"
)

const maxLen = 8000000

func (b *Base) PrepareFile(d *discordgo.File, len int) *discordgo.File {
	if len > maxLen {
		buf := bytes.NewBuffer(nil)
		gz := gzip.NewWriter(buf)
		_, err := io.Copy(gz, d.Reader)
		if err != nil {
			return d
		}
		err = gz.Close()
		if err != nil {
			return d
		}

		d.ContentType = "application/gzip"
		d.Name += ".gz"
		d.Reader = buf
	}
	return d
}
