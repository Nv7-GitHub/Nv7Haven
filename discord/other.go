package discord

import (
	"fmt"
	"log"
	"math/rand"
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

	if strings.HasPrefix(m.Content, "insult") {
		start := starters[rand.Intn(len(starters))]
		log.Println(start)

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
}
