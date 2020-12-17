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
			diff := time.Now().Unix() - int64(user.Metadata["lastdaily"].(float64))
			if (diff) < 86400 { // less than a day
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You still need to wait %0.2f hours.", float32(diff)/3600))
				return
			}
		}

		user.Wallet += 2500
		success = b.updateuser(m, user)
		if !success {
			return
		}
		s.ChannelMessageSend(m.ChannelID, "Congrats on the 2,500 coins! Come back in 24 hours to get more!")
		return
	}

	if strings.HasPrefix(m.Content, "bal") {
		b.checkuser(m)
		user, success := b.getuser(m, m.Author.ID)
		if !success {
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You have %d coins in your wallet and %d coins in the bank.", user.Wallet, user.Bank))
		return
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
		return
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
		return
	}

	if strings.HasPrefix(m.Content, "ldb") {
		res, err := b.db.Query("SELECT user, wallet FROM currency WHERE guilds LIKE ? ORDER BY wallet DESC LIMIT 10", "%"+m.GuildID+"%")
		if b.handle(err, m) {
			return
		}

		var ldb string
		var user string
		var wallet int
		var usr *discordgo.User
		for res.Next() {
			err = res.Scan(&user, &wallet)
			if b.handle(err, m) {
				return
			}
			usr, err = s.User(user)
			if b.handle(err, m) {
				return
			}

			ldb += fmt.Sprintf("%s#%s - %d\n", usr.Username, usr.Discriminator, wallet)
		}

		gld, err := s.Guild(m.GuildID)
		if b.handle(err, m) {
			return
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("Richest users in %s", gld.Name),
			Description: ldb,
		})
		return
	}

	if strings.HasPrefix(m.Content, "credup") {
		user, suc := b.getuser(m, m.Author.ID)
		if !suc {
			return
		}

		var num int
		_, err := fmt.Sscanf(m.Content, "credup %d", &num)
		if b.handle(err, m) {
			return
		}

		price := (num * num) - (user.Credit * user.Credit)
		if user.Wallet < price {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You need %d more coins to upgrade your credit %d levels.", price-user.Wallet, num))
			return
		}

		user.Wallet -= price
		user.Credit += price
		suc = b.updateuser(m, user)
		if !suc {
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You upgraded your credit by %d levels!", num))
		return
	}
}
