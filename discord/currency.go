package discord

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) currencyBasics(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "daily") {
		b.checkuser(m)
		user, success := b.getuser(m, m.Author.ID)
		if !success {
			return
		}
		_, exists := user.Metadata["lastdaily"]
		if !exists {
			user.Metadata["lastdaily"] = time.Now().Unix()
		} else {
			diff := time.Now().Unix() - user.Metadata["lastdaily"].(int64)
			if (diff) < 86400 { // less than a day
				s.ChannelMessageSend(m.ChannelID, "You still need to wait "+strconv.Itoa(int(diff/3600))+" hours.")
				return
			}
		}

		user.Wallet += 2500
		success = b.updateuser(m, user)
		if !success {
			return
		}
		s.ChannelMessageSend(m.ChannelID, "Congrats on the 2,500 coins! Come back in 24 hours to get more!")
	}

	if strings.HasPrefix(m.Content, "bal") {
		b.checkuser(m)
		user, success := b.getuser(m, m.Author.ID)
		if !success {
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You have %d coins in your wallet and %d coins in the bank.", user.Wallet, user.Bank))
	}
}
