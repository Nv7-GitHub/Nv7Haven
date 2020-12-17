package discord

import (
	"encoding/json"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) exists(m *discordgo.MessageCreate, table string, where string, args ...interface{}) (bool, bool) {
	res, err := b.db.Query("SELECT COUNT(1) FROM "+table+" WHERE "+where+" LIMIT 1", args...)
	if b.handle(err, m) {
		return false, false
	}
	defer res.Close()
	res.Next()

	var count int
	err = res.Scan(&count)
	if b.handle(err, m) {
		return false, false
	}
	return count != 0, true
}

func (b *Bot) handle(err error, m *discordgo.MessageCreate) bool {
	if err != nil {
		b.dg.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
		return true
	}
	return false
}

func (b *Bot) unmarshal(m *discordgo.MessageCreate, data string, out interface{}) {
	err := json.Unmarshal([]byte(data), &out)
	if b.handle(err, m) {
		return
	}
}

type user struct {
	User        string
	Wallet      int
	Bank        int
	Credit      int
	Properties  []string // Places they own
	LastVisited int64
	Metadata    map[string]interface{}
}

func (b *Bot) getuser(m *discordgo.MessageCreate, usr string) (user, bool) {
	res, err := b.db.Query("SELECT * FROM currency WHERE user=?", usr)
	if b.handle(err, m) {
		return user{}, false
	}
	defer res.Close()
	res.Next()
	var name string
	var wallet int
	var bank int
	var credit int
	var props string
	var lastvisited int64
	var met string
	err = res.Scan(&name, &wallet, &bank, &credit, &props, &lastvisited, &met)
	if b.handle(err, m) {
		return user{}, false
	}
	var properties []string
	var metadata map[string]interface{}
	err = json.Unmarshal([]byte(props), &properties)
	if b.handle(err, m) {
		return user{}, false
	}
	err = json.Unmarshal([]byte(met), &metadata)
	if b.handle(err, m) {
		return user{}, false
	}
	return user{
		User:        name,
		Wallet:      wallet,
		Bank:        bank,
		Credit:      credit,
		Properties:  properties,
		LastVisited: lastvisited,
		Metadata:    metadata,
	}, true
}

func (b *Bot) updateuser(m *discordgo.MessageCreate, u user) bool {
	met, err := json.Marshal(u.Metadata)
	if b.handle(err, m) {
		return false
	}
	props, err := json.Marshal(u.Properties)
	if b.handle(err, m) {
		return false
	}
	_, err = b.db.Exec("UPDATE currency SET wallet=?, bank=?, credit=?, properties=?, lastvisited=?, metadata=? WHERE user=?", u.Wallet, u.Bank, u.Credit, props, u.LastVisited, met, u.User)
	if b.handle(err, m) {
		return false
	}
	return true
}

func (b *Bot) checkuser(m *discordgo.MessageCreate) {
	exists, success := b.exists(m, "currency", "user=?", m.Author.ID)
	if !success {
		return
	}
	if !exists {
		_, err := b.db.Exec("INSERT INTO currency VALUES ( ?, ?, ?, ?, ?, ?, ? )", m.Author.ID, 0, 0, 0, "[]", time.Now().Unix()-86400, "{}")
		if b.handle(err, m) {
			return
		}
	}
}
