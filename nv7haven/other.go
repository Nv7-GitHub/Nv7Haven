package nv7haven

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (d *Nv7Haven) getIP(c *fiber.Ctx) error {
	return c.SendString(c.IPs()[0])
}

func (n *Nv7Haven) httpGet(c *fiber.Ctx) error {
	resp, err := http.Get(string(c.Body()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.SendStream(resp.Body)
}

func (n *Nv7Haven) getURL(c *fiber.Ctx) error {
	id := c.Params("id")
	var jsonData = []byte(`{ "context": { "client": { "hl": "en", "clientName": "WEB", "clientVersion": "2.20210721.00.00" } }, "videoId": "`+id+`" }`)
	request, error := http.NewRequest("POST", "https://youtubei.googleapis.com/youtubei/v1/player?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8", bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	
	client := &http.Client{}
	resp, error := client.Do(request)
	if error != nil {
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
		out[i].MimeType = format.MimeType
		out[i].Size = (format.Bitrate * duration) / 8
		out[i].SizeFormatted = FormatByteSize(out[i].Size / 1024)
	}

	return c.JSON(ytResp{
		Results:   out,
		Title:     strings.ReplaceAll(d.Details.Title, "+", " "),
		Thumbnail: d.Details.Thumbnail.Thumbnails[len(d.Details.Thumbnail.Thumbnails)-1].URL,
	})
}

type ytResponse struct {
	StreamingData ytStreamingData `json:"streamingData"`
	Details       ytDetails       `json:"videoDetails"`
}

type ytStreamingData struct {
	Formats []ytFormat `json:"formats"`
}

type ytFormat struct {
	URL      string `json:"url"`
	Quality  string `json:"qualityLabel"`
	Duration string `json:"approxDurationMs"`
	Bitrate  int    `json:"bitrate"`
	MimeType string `json:"mimeType"`
}

type ytDetails struct {
	Title     string         `json:"title"`
	Thumbnail ytThumbnailDat `json:"thumbnail"`
}

type ytThumbnailDat struct {
	Thumbnails []ytThumbnail `json:"thumbnails"`
}

type ytThumbnail struct {
	URL string `json:"url"`
}

type ytOut struct {
	Size          int
	SizeFormatted string
	URL           string
	MimeType      string
	Quality       string
}

type ytResp struct {
	Results   []ytOut
	Thumbnail string
	Title     string
}
