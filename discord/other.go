package discord

import (
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var starters = []string{"You utter %s", "You collection of %s", "You are similar to %s", "You are like a %s", "Utter %s. Go %s yourself and consume many %ss", "You are SO many %ss. On an unrelated note, anti %s eradication protocols have been initiated.", "Sometimes you can be confused for %f %s-like entities.", "You are sometimes like a %s", "You are like %d %s %sic toolkits, each %dcm in diameter.", "You are like a %s made of %s"}
var words = []string{"vertices", "newscast", "baryon", "widget", "hyperboloid", "communism", "django", "transport", "apioform"}
var replaces = []string{"%s", "%d", "%f"}

func (b *Bot) other(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if b.startsWith(m, "insult") {
		start := starters[rand.Intn(len(starters))]

		for _, val := range replaces {
			for strings.Contains(start, val) {
				var replace string
				switch val {
				case "%s":
					replace = words[rand.Intn(len(words))]
					break
				case "%d":
					replace = strconv.Itoa(rand.Intn(100000))
					break
				case "%f":
					replace = fmt.Sprintf("%0.4f", float32(rand.Intn(100000))+rand.Float32())
					break
				}
				start = strings.Replace(start, val, replace, 1)
			}
		}

		s.ChannelMessageSend(m.ChannelID, start)
	}

	if b.startsWith(m, "gh") {
		var query string
		_, err := fmt.Sscanf(m.Content, "gh %s", &query)
		if b.handle(err, m) {
			return
		}

		var out ghSearch
		suc := b.req(m, "https://api.github.com/search/repositories?q="+url.PathEscape(query), &out)
		if !suc {
			return
		}

		if len(out.Items) == 0 {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(`No results for "%s" found.`, query))
			return
		}

		item := out.Items[0]

		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{
				Name:    item.Owner.Login,
				IconURL: item.Owner.AvatarURL,
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Repository",
					Value:  fmt.Sprintf("[%s](%s)", item.Name, item.HTMLURL),
					Inline: true,
				},
				{
					Name:   "Language",
					Value:  item.Language,
					Inline: true,
				},
				{
					Name:   "Forks",
					Value:  strconv.Itoa(item.Forks),
					Inline: true,
				},
				{
					Name:   "Watchers",
					Value:  strconv.Itoa(item.Watchers),
					Inline: true,
				},
				{
					Name:   "Open Issues",
					Value:  strconv.Itoa(item.OpenIssuesCount),
					Inline: true,
				},
			},
		}

		if item.License.Name != "" {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "License",
				Value:  item.License.Name,
				Inline: true,
			})
		}

		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		return
	}

	if strings.HasPrefix(m.Content, "setprefix") {
		if !b.isMod(m, m.Author.ID) {
			s.ChannelMessageSend(m.ChannelID, "You need permission `Administrator` to use this command!")
			return
		}
		var prefix string
		_, err := fmt.Sscanf(m.Content, "setprefix %s", &prefix)
		if err != nil {
			prefix = ""
		}
		_, err = b.db.Exec("UPDATE prefixes SET prefix=? WHERE guild=?", prefix, m.GuildID)
		if b.handle(err, m) {
			return
		}
		b.prefixcache[m.GuildID] = prefix
		s.ChannelMessageSend(m.ChannelID, "Successfully update prefix!")
	}

	if b.startsWith(m, "rate") {
		toRate := []byte(m.Content[5:])
		var seed int64 = 0
		for i, val := range toRate {
			seed += int64(math.Pow(10, float64(i))) * int64(val)
		}
		rand.Seed(seed)
		rating := rand.Intn(12)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s is a %d / 10!", toRate, rating))
	}
}

type ghSearch struct {
	Items []gh
}

type gh struct {
	Name            string
	Forks           int
	Language        string
	OpenIssuesCount int `json:"open_issues_count"`
	Watchers        int
	License         ghLicense
	Owner           ghOwner
	HTMLURL         string `json:"html_url"`
}

type ghLicense struct {
	Name string
}

type ghOwner struct {
	Login     string
	AvatarURL string `json:"avatar_url"`
}
