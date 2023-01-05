package discord

import (
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) giveNumCmd(num int, m msg, rsp rsp) {
	num = b.abs(num)
	res, err := b.db.Query("SELECT COUNT(*) FROM givenum WHERE guild=? AND member=? LIMIT 1", m.GuildID, m.Author.ID)
	if rsp.Error(err) {
		return
	}
	defer res.Close()
	var count int
	if (num > 100) || (num < 0) {
		rsp.ErrorMessage("You need to choose a number from 0-100.")
		return
	}
	res.Next()
	err = res.Scan(&count)
	if rsp.Error(err) {
		return
	}
	if count == 0 {
		_, err = b.db.Exec("INSERT INTO givenum VALUES ( ?, ?, ? )", m.GuildID, m.Author.ID, num)
		if rsp.Error(err) {
			return
		}
		rsp.Resp("Successfully initialized value.")
		return
	}
	_, err = b.db.Exec("UPDATE givenum SET number=? WHERE guild=? AND member=?", num, m.GuildID, m.Author.ID)
	if rsp.Error(err) {
		return
	}
	rsp.Resp("Successfully updated value.")
}

func (b *Bot) getNumCmd(hasMention bool, mention string, m msg, rsp rsp) {
	if !hasMention {
		success, num := b.getNum(m, rsp, m.Author.ID)
		if success {
			rsp.Message(fmt.Sprintf("Your number is %d.", num))
		}
		return
	}
	success, num := b.getNum(m, rsp, mention)
	if success {
		rsp.Message(fmt.Sprintf("<@%s>'s number is %d.", mention, num))
	}
}

func (b *Bot) randselectCmd(m msg, rsp rsp) {
	if b.isUserMod(m, rsp, m.Author.ID) {
		res, err := b.db.Query("SELECT number FROM givenum WHERE guild=?", m.GuildID)
		if rsp.Error(err) {
			return
		}
		defer res.Close()
		nums := make([]int, 0)
		for res.Next() {
			var data int
			err = res.Scan(&data)
			if rsp.Error(err) {
				return
			}
			nums = append(nums, data)
		}
		num := nums[rand.Intn(len(nums))]
		rsp.Message(fmt.Sprintf("The number was %d.", num))
		res, err = b.db.Query("SELECT member FROM givenum WHERE guild=? AND number=?", m.GuildID, num)
		if rsp.Error(err) {
			return
		}
		defer res.Close()
		for res.Next() {
			var memberid string
			err = res.Scan(&memberid)
			if rsp.Error(err) {
				return
			}
			rsp.Message(fmt.Sprintf("<@%s> got it right!", memberid))
		}
		return
	}
	rsp.Message(`You need to have permission "Administrator" to use this command.`)
}

func (b *Bot) giveNum(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	// givenum command
	if b.startsWith(m, "givenum") {
		var num int
		_, err := fmt.Sscanf(m.Content, "givenum %d", &num)
		if b.handle(err, m) {
			return
		}
		b.giveNumCmd(num, b.newMsgNormal(m), b.newRespNormal(m))
		return
	}

	// getnum command
	if b.startsWith(m, "getnum") {
		mention := ""
		if len(m.Mentions) > 0 {
			mention = m.Mentions[0].ID
		}
		b.getNumCmd(len(m.Mentions) > 0, mention, b.newMsgNormal(m), b.newRespNormal(m))
		return
	}

	// randselect command
	if b.startsWith(m, "randselect") {
		b.randselectCmd(b.newMsgNormal(m), b.newRespNormal(m))
		return
	}
}

func (b *Bot) getNum(m msg, rsp rsp, user string) (bool, int) {
	res, err := b.db.Query("SELECT COUNT(*) FROM givenum WHERE guild=? AND member=? LIMIT 1", m.GuildID, user)
	if rsp.Error(err) {
		return false, 0
	}
	defer res.Close()
	var count int
	res.Next()
	err = res.Scan(&count)
	if rsp.Error(err) {
		return false, 0
	}
	if count == 0 {
		rsp.Message("User <@" + user + "> has not chosen a number.")
		return false, 0
	}

	res, err = b.db.Query("SELECT number FROM givenum WHERE guild=? AND member=? LIMIT 1", m.GuildID, user)
	if rsp.Error(err) {
		return false, 0
	}
	defer res.Close()
	res.Next()
	var num int
	err = res.Scan(&num)
	if rsp.Error(err) {
		return false, 0
	}
	return true, num
}
