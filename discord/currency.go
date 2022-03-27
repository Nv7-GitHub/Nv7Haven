package discord

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const leftArrow = "⬅️"
const rightArrow = "➡️"

type reactionMsg struct {
	Type     reactionMsgType
	Metadata map[string]any
	Handler  func(*discordgo.MessageReactionAdd)
}

func (b *Bot) ldbPageSwitcher(r *discordgo.MessageReactionAdd) {
	pg := b.pages[r.MessageID]
	var page int
	if r.Emoji.Name == leftArrow {
		page = pg.Metadata["page"].(int) - 1
	} else {
		page = pg.Metadata["page"].(int) + 1
	}
	if ((page * 10) > pg.Metadata["count"].(int)) || (page < 0) {
		b.dg.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.Name, r.UserID)
		b.pages[r.MessageID] = pg
		return
	}
	pg.Metadata["page"] = page

	res, err := b.db.Query("SELECT user, wallet FROM currency WHERE guilds LIKE ? ORDER BY wallet DESC LIMIT 10 OFFSET ?", "%"+r.GuildID+"%", page*10)
	if err != nil {
		return
	}
	defer res.Close()

	var ldb string
	var user string
	var wallet int
	var usr *discordgo.User
	i := 1 + (page * 10)
	for res.Next() {
		err = res.Scan(&user, &wallet)
		if err != nil {
			return
		}
		usr, err = b.dg.User(user)
		if err != nil {
			return
		}

		ldb += fmt.Sprintf("%d. %s#%s - %d\n", i, usr.Username, usr.Discriminator, wallet)
		i++
	}

	gld, err := b.dg.Guild(r.GuildID)
	if err != nil {
		return
	}
	b.dg.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Richest users in %s", gld.Name),
		Description: ldb,
	})
	b.dg.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.Name, r.UserID)
	b.pages[r.MessageID] = pg
}

func (b *Bot) currencyBasics(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if b.startsWith(m, "daily") {
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
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You still need to wait %0.2f hours.", 24-(float32(diff)/3600)))
				return
			}
		}

		user.Wallet += 2500
		user.Metadata["lastdaily"] = time.Now().Unix()
		success = b.updateuser(m, user)
		if !success {
			return
		}
		s.ChannelMessageSend(m.ChannelID, "Congrats on the 2,500 coins! Come back in 24 hours to get more!")
		return
	}

	if b.startsWith(m, "bal") {
		b.checkuser(m)
		id := m.Author.ID
		person := "You have"
		describer := "your"
		if len(m.Mentions) > 0 {
			id = m.Mentions[0].ID
			b.checkuserwithid(m, id)
			person = "<@" + id + "> has"
			describer = "their"
		}
		user, success := b.getuser(m, id)
		if !success {
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s %d coins in %s wallet and %d coins in the bank. %s credit is %d.", person, user.Wallet, describer, user.Bank, strings.Title(describer), user.Credit))
		return
	}

	if b.startsWith(m, "dep") {
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
			num = b.abs(num)
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

	if b.startsWith(m, "with") {
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
			num = b.abs(num)
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

	if b.startsWith(m, "ldb") {
		count := b.db.QueryRow("SELECT COUNT(1) FROM currency WHERE guilds LIKE ?", "%"+m.GuildID+"%")
		var num int
		err := count.Scan(&num)
		if b.handle(err, m) {
			return
		}

		res, err := b.db.Query("SELECT user, wallet FROM currency WHERE guilds LIKE ? ORDER BY wallet DESC LIMIT 10", "%"+m.GuildID+"%")
		if b.handle(err, m) {
			return
		}
		defer res.Close()

		var ldb string
		var user string
		var wallet int
		var usr *discordgo.User
		i := 1
		for res.Next() {
			err = res.Scan(&user, &wallet)
			if b.handle(err, m) {
				return
			}
			usr, err = s.User(user)
			if b.handle(err, m) {
				return
			}

			ldb += fmt.Sprintf("%d. %s#%s - %d\n", i, usr.Username, usr.Discriminator, wallet)
			i++
		}

		gld, err := s.Guild(m.GuildID)
		if b.handle(err, m) {
			return
		}
		msg, _ := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("Richest users in %s", gld.Name),
			Description: ldb,
		})
		s.MessageReactionAdd(m.ChannelID, msg.ID, leftArrow)
		s.MessageReactionAdd(m.ChannelID, msg.ID, rightArrow)
		b.pages[msg.ID] = reactionMsg{
			Type: ldbPageSwitcher,
			Metadata: map[string]any{
				"page":  0,
				"count": num,
			},
			Handler: b.ldbPageSwitcher,
		}
		return
	}

	if b.startsWith(m, "credup") {
		user, suc := b.getuser(m, m.Author.ID)
		if !suc {
			return
		}

		var numVal string
		_, err := fmt.Sscanf(m.Content, "credup %s", &numVal)
		if b.handle(err, m) {
			return
		}

		var num int
		if numVal == "max" {
			price := 0
			for price < user.Wallet {
				numoff := num + user.Credit
				price = (numoff * numoff) - (user.Credit * user.Credit)
				num++
			}
			if num > 0 {
				num--
			}
		} else {
			num, err = strconv.Atoi(numVal)
			if b.handle(err, m) {
				return
			}
			num = b.abs(num)
		}

		numoff := num + user.Credit
		price := (numoff * numoff) - (user.Credit * user.Credit)
		if user.Wallet < price {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You need %d more coins to upgrade your credit %d levels.", price-user.Wallet, num))
			return
		}

		user.Wallet -= price
		user.Credit += num
		suc = b.updateuser(m, user)
		if !suc {
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You upgraded your credit by %d levels!", num))
		return
	}

	if b.startsWith(m, "donate") {
		b.checkuser(m)
		if !(len(m.Mentions) > 0) {
			s.ChannelMessageSend(m.ChannelID, "You need to mention the person you are going to donate to!")
			return
		}
		if m.Mentions[0].ID == m.Author.ID {
			s.ChannelMessageSend(m.ChannelID, "You can't donate to yourself!")
			return
		}
		b.checkuserwithid(m, m.Mentions[0].ID)

		user1, suc := b.getuser(m, m.Author.ID)
		if !suc {
			return
		}

		var num int
		_, err := fmt.Sscanf(m.Content, "donate %d", &num)
		if b.handle(err, m) {
			return
		}
		num = b.abs(num)

		if user1.Wallet < num {
			s.ChannelMessageSend(m.ChannelID, "You don't have that much money to give!")
			return
		}

		user2, suc := b.getuser(m, m.Mentions[0].ID)
		if !suc {
			return
		}

		user1.Wallet -= num
		user2.Wallet += num

		b.updateuser(m, user1)
		b.updateuser(m, user2)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Successfully donated %d coins to <@%s>!", num, m.Mentions[0].ID))
	}
}
