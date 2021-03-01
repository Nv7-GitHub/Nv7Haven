package discord

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type redditResp struct {
	Data redditData
}

type redditData struct {
	Children []previewData
	After    string
}

type previewData struct {
	Data meme
}

type redditVideo struct {
	Media media `json:"reddit_video"`
}

type meme struct {
	URL       string
	Title     string
	Permalink string
	IsVideo   bool        `json:"is_video"`
	Media     redditVideo `json:"media"`
}

type media struct {
	FallbackURL string `json:"fallback_url"`
}

func (b *Bot) makeMemeEmbed(m meme, msg msg) *discordgo.MessageEmbed {
	mE := &discordgo.MessageEmbed{
		URL:   m.Permalink,
		Type:  discordgo.EmbedTypeImage,
		Title: m.Title,
	}
	if m.IsVideo {
		mE.Video = &discordgo.MessageEmbedVideo{
			URL: m.Media.Media.FallbackURL,
		}
		b.dg.ChannelMessageSend(msg.ChannelID, m.Media.Media.FallbackURL)
	} else if !strings.Contains(m.URL, "youtu") {
		mE.Image = &discordgo.MessageEmbedImage{
			URL: m.URL,
		}
	} else {
		b.dg.ChannelMessageSend(msg.ChannelID, m.URL)
	}
	return mE
}

func (b *Bot) memes(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if b.startsWith(m, "meme") {
		b.memeCommand(b.newMsgNormal(m), b.newRespNormal(m))
	}

	if b.startsWith(m, "cmeme") {
		b.cmemeCommand(b.newMsgNormal(m), b.newRespNormal(m))
	}

	if b.startsWith(m, "pmeme") {
		b.pmemeCommand(b.newMsgNormal(m), b.newRespNormal(m))
	}
}

func (b *Bot) startMemes(m msg, rsp rsp) (bool, bool) {
	hasResponded := false
	if len(b.memedat) == 0 { // first time since startup
		rsp.Message("Sorry, this is the first time someone has asked for a meme since the server started. It may take a moment for us to download the memes from reddit.")
		hasResponded = true
		success := b.loadMemes(rsp)
		if !success {
			b.dg.ChannelMessageSend(m.ChannelID, "Failed to load memes")
			return false, false
		}
	}
	if (time.Now().Sub(b.memerefreshtime)).Hours() >= 1 { // its been an hour
		go b.loadMemes(rsp)
	}
	return hasResponded, true
}

func (b *Bot) pmemeCommand(m msg, rsp rsp) {
	hasResponded, suc := b.startMemes(m, rsp)
	if !suc {
		return
	}

	// send message
	unique := false
	var randnum int
	if len(b.pmemecache[m.GuildID]) == len(b.pmemedat) {
		b.pmemecache[m.GuildID] = make(map[int]empty, 0)
	}
	for !unique {
		randnum = rand.Intn(len(b.pmemedat))
		unique = true
		_, exists := b.pmemecache[m.GuildID]
		if !exists {
			b.pmemecache[m.GuildID] = make(map[int]empty, 0)
			unique = true
			break
		}
		_, unique = b.pmemecache[m.GuildID][randnum]
		unique = !unique
	}
	b.pmemecache[m.GuildID][randnum] = empty{}
	meme := b.pmemedat[randnum]
	emb := b.makeMemeEmbed(meme, m)
	if hasResponded {
		b.dg.ChannelMessageSendEmbed(m.ChannelID, emb)
		return
	}
	rsp.Embed(emb)
}

func (b *Bot) memeCommand(m msg, rsp rsp) {
	hasResponded, suc := b.startMemes(m, rsp)
	if !suc {
		return
	}

	// send message
	unique := false
	var randnum int
	if len(b.memecache[m.GuildID]) == len(b.memedat) {
		b.memecache[m.GuildID] = make(map[int]empty, 0)
	}
	for !unique {
		randnum = rand.Intn(len(b.memedat))
		unique = true
		_, exists := b.memecache[m.GuildID]
		if !exists {
			b.memecache[m.GuildID] = make(map[int]empty, 0)
			unique = true
			break
		}
		_, unique = b.memecache[m.GuildID][randnum]
		unique = !unique
	}
	b.memecache[m.GuildID][randnum] = empty{}
	meme := b.memedat[randnum]
	emb := b.makeMemeEmbed(meme, m)
	if hasResponded {
		b.dg.ChannelMessageSendEmbed(m.ChannelID, emb)
		return
	}
	rsp.Embed(emb)
}

func (b *Bot) cmemeCommand(m msg, rsp rsp) {
	hasResponded, suc := b.startMemes(m, rsp)
	if !suc {
		return
	}

	// send message
	unique := false
	var randnum int
	if len(b.cmemecache[m.GuildID]) == len(b.cmemedat) {
		b.cmemecache[m.GuildID] = make(map[int]empty, 0)
	}
	for !unique {
		randnum = rand.Intn(len(b.cmemedat))
		unique = true
		_, exists := b.cmemecache[m.GuildID]
		if !exists {
			b.cmemecache[m.GuildID] = make(map[int]empty, 0)
			unique = true
			break
		}
		_, unique = b.cmemecache[m.GuildID][randnum]
		unique = !unique
	}
	b.cmemecache[m.GuildID][randnum] = empty{}
	meme := b.cmemedat[randnum]
	emb := b.makeMemeEmbed(meme, m)
	if hasResponded {
		b.dg.ChannelMessageSendEmbed(m.ChannelID, emb)
		return
	}
	rsp.Embed(emb)
}

func (b *Bot) loadMemes(rsp rsp) bool {
	b.memerefreshtime = time.Now()
	b.memecache = make(map[string]map[int]empty, 0)
	b.cmemecache = make(map[string]map[int]empty, 0)
	b.pmemecache = make(map[string]map[int]empty, 0)
	var suc bool
	b.memedat, suc = b.downloadMeme(rsp, "memes")
	if !suc {
		return false
	}
	b.cmemedat, suc = b.downloadMeme(rsp, "cleanmemes")
	if !suc {
		return false
	}
	b.pmemedat, suc = b.downloadMeme(rsp, "ProgrammerHumor")
	if !suc {
		return false
	}
	return true
}

func (b *Bot) downloadMeme(rsp rsp, subreddit string) ([]meme, bool) {
	children := make([]previewData, 0)
	after := ""
	for len(children) < 200 {
		// Download
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://reddit.com/r/"+subreddit+"/hot.json?after="+after, nil)
		if rsp.Error(err) {
			return nil, false
		}
		req.Header.Set("User-Agent", "Nv7 Bot")
		res, err := client.Do(req)
		if rsp.Error(err) {
			return nil, false
		}
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if rsp.Error(err) {
			return nil, false
		}

		// Process
		var dat redditResp
		err = json.Unmarshal(data, &dat)
		if rsp.Error(err) {
			return nil, false
		}
		children = append(children, dat.Data.Children...)
		after = dat.Data.After
	}

	memedat := make([]meme, len(children))
	for i, val := range children {
		memedat[i] = val.Data
		memedat[i].Permalink = "https://reddit.com" + memedat[i].Permalink
	}
	return memedat, true
}
