package nv7haven

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (d *Nv7Haven) getIP(c *fiber.Ctx) error {
	return c.SendString(c.IPs()[0])
}

func (n *Nv7Haven) getURL(c *fiber.Ctx) error {
	id := c.Params("id")
	link := fmt.Sprintf("https://www.youtube.com/get_video_info?video_id=%s", id)
	resp, err := http.Get(link)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	u, err := url.Parse("?" + string(data))
	if err != nil {
		return err
	}
	query := u.Query()

	dat := query.Get("player_response")
	var d ytResponse
	err = json.Unmarshal([]byte(dat), &d)
	if err != nil {
		return err
	}
	dats := d.StreamingData.Formats

	out := make([]ytOut, len(dats))
	for i, format := range dats {
		duration, err := strconv.Atoi(format.Duration)
		if err != nil {
			return err
		}
		out[i].URL, err = url.PathUnescape(format.URL)
		if err != nil {
			return err
		}
		out[i].Quality = format.Quality
		out[i].Size = (format.Bitrate * duration) / 8
		out[i].SizeFormatted = FormatByteSize(out[i].Size / 1024)
	}

	return c.JSON(out)
}

type ytResponse struct {
	StreamingData ytStreamingData `json:"streamingData"`
}

type ytStreamingData struct {
	Formats []ytFormat `json:"formats"`
}

type ytFormat struct {
	URL      string `json:"url"`
	Quality  string `json:"qualityLabel"`
	Duration string `json:"approxDurationMs"`
	Bitrate  int    `json:"bitrate"`
}

type ytOut struct {
	Size          int
	SizeFormatted string
	URL           string
	Quality       string
}
