package discord

import (
	"log"
	"math/rand"
	"sort"
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
		log.Println("insult")
		start := starters[rand.Intn(len(starters))]
		log.Println(start)

		argTypes := make([]string, 0)
		argS := start
		indexes := make([]int, len(replaces))
		isOver := false
		for !isOver {
			isOver = true
			for i, val := range replaces {
				indexes[i] = strings.Index(argS, val)
				if isOver && indexes[i] > -1 {
					isOver = false
				}
			}
			if !isOver {
				sort.Ints(indexes)
				for i, val := range indexes {
					if val > -1 {
						argTypes = append(argTypes, argS[val:val+2])
						if i == len(indexes)-1 {
							argS = argS[val : val+2]
						}
						log.Println(argTypes)
					}
				}
			}
		}

		log.Println(start, argTypes)
	}
}
