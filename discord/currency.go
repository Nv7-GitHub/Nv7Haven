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

	if strings.HasPrefix(m.Content, "dep") {
		b.checkuser(m)
		user, success := b.getuser(m, m.Author.ID)
		if !success {
			return
		}

		var dep string
		_, err := fmt.Sscanf(m.Content, "dep %s", &dep)
		if b.handle(err, m) {
			return
		}
		var num int
		if dep == "all" {
			num = user.Wallet
		} else {
			num, err = strconv.Atoi(dep)
			if b.handle(err, m) {
				return
			}
		}
		if user.Wallet < num {
			num = user.Wallet
		}
		user.Bank += num
		user.Wallet -= num
		success = b.updateuser(m, user)
		if !success {
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Deposited %d coins.", num))
	}

	if strings.HasPrefix(m.Content, "with") {
		b.checkuser(m)
		user, success := b.getuser(m, m.Author.ID)
		if !success {
			return
		}

		var with string
		_, err := fmt.Sscanf(m.Content, "with %s", &with)
		if b.handle(err, m) {
			return
		}
		var num int
		if with == "all" {
			num = user.Bank
		} else {
			num, err = strconv.Atoi(with)
			if b.handle(err, m) {
				return
			}
		}
		if user.Bank < num {
			num = user.Bank
		}
		user.Bank -= num
		user.Wallet += num
		success = b.updateuser(m, user)
		if !success {
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Withdrew %d coins.", num))
	}
}
