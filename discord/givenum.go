package discord

import (
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) giveNum(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// givenum command
	if b.startsWith(m, "givenum") {
		var num int
		_, err := fmt.Sscanf(m.Content, "givenum %d", &num)
		if b.handle(err, m) {
			return
		}
		num = b.abs(num)
		res, err := b.db.Query("SELECT COUNT(1) FROM givenum WHERE guild=? AND member=? LIMIT 1", m.GuildID, m.Author.ID)
		defer res.Close()
		if b.handle(err, m) {
			return
		}
		var count int
		if (num > 100) || (num < 0) {
			s.ChannelMessageSend(m.ChannelID, "You need to choose a number from 0-100.")
			return
		}
		res.Next()
		err = res.Scan(&count)
		if b.handle(err, m) {
			return
		}
		if count == 0 {
			_, err = b.db.Exec("INSERT INTO givenum VALUES ( ?, ?, ? )", m.GuildID, m.Author.ID, num)
			if b.handle(err, m) {
				return
			}
			s.ChannelMessageSend(m.ChannelID, "Successfully initialized value.")
			return
		}
		_, err = b.db.Exec("UPDATE givenum SET number=? WHERE guild=? AND member=?", num, m.GuildID, m.Author.ID)
		if b.handle(err, m) {
			return
		}
		s.ChannelMessageSend(m.ChannelID, "Successfully updated value.")
		return
	}

	// getnum command
	if b.startsWith(m, "getnum") {
		success, num := b.getNum(m, m.Author.ID)
		if success {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Your number is %d.", num))
		}
		for _, user := range m.Mentions {
			success, num = b.getNum(m, user.ID)
			if success {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>'s number is %d.", user.ID, num))
			}
		}
		return
	}

	// randselect command
	if b.startsWith(m, "randselect") {
		if b.isMod(m, m.Author.ID) {
			res, err := b.db.Query("SELECT number FROM givenum WHERE guild=?", m.GuildID)
			if b.handle(err, m) {
				return
			}
			defer res.Close()
			nums := make([]int, 0)
			for res.Next() {
				var data int
				err = res.Scan(&data)
				if b.handle(err, m) {
					return
				}
				nums = append(nums, data)
			}
			num := nums[rand.Intn(len(nums))]
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The number was %d.", num))
			res, err = b.db.Query("SELECT member FROM givenum WHERE guild=? AND number=?", m.GuildID, num)
			if b.handle(err, m) {
				return
			}
			defer res.Close()
			for res.Next() {
				var memberid string
				err = res.Scan(&memberid)
				if b.handle(err, m) {
					return
				}
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> got it right!", memberid))
			}
			return
		}
		s.ChannelMessageSend(m.ChannelID, `You need to have permission "Administrator" to use this command.`)
		return
	}
}

func (b *Bot) getNum(m *discordgo.MessageCreate, user string) (bool, int) {
	res, err := b.db.Query("SELECT COUNT(1) FROM givenum WHERE guild=? AND member=? LIMIT 1", m.GuildID, user)
	defer res.Close()
	if b.handle(err, m) {
		return false, 0
	}
	var count int
	res.Next()
	err = res.Scan(&count)
	if b.handle(err, m) {
		return false, 0
	}
	if count == 0 {
		b.dg.ChannelMessageSend(m.ChannelID, "User <@"+user+"> has not chosen a number.")
		return false, 0
	}

	res, err = b.db.Query("SELECT number FROM givenum WHERE guild=? AND member=? LIMIT 1", m.GuildID, user)
	defer res.Close()
	if b.handle(err, m) {
		return false, 0
	}
	res.Next()
	var num int
	err = res.Scan(&num)
	if b.handle(err, m) {
		return false, 0
	}
	return true, num
}
