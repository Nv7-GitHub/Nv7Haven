package discord

import (
	"encoding/json"
	"math"
	"net/http"
	"strings"
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
	Guilds      []string
	Wallet      int
	Bank        int
	Credit      int
	Properties  map[string]int // Places they own
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
	var glds string
	var wallet int
	var bank int
	var credit int
	var props string
	var lastvisited int64
	var met string
	err = res.Scan(&name, &glds, &wallet, &bank, &credit, &props, &lastvisited, &met)
	if b.handle(err, m) {
		return user{}, false
	}
	var guilds []string
	err = json.Unmarshal([]byte(glds), &guilds)
	if b.handle(err, m) {
		return user{}, false
	}
	var properties map[string]int
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
		Guilds:      guilds,
		Wallet:      wallet,
		Bank:        bank,
		Credit:      credit,
		Properties:  properties,
		LastVisited: lastvisited,
		Metadata:    metadata,
	}, true
}

func (b *Bot) updateuser(m *discordgo.MessageCreate, u user) bool {
	glds, err := json.Marshal(u.Guilds)
	if b.handle(err, m) {
		return false
	}
	met, err := json.Marshal(u.Metadata)
	if b.handle(err, m) {
		return false
	}
	props, err := json.Marshal(u.Properties)
	if b.handle(err, m) {
		return false
	}
	_, err = b.db.Exec("UPDATE currency SET guilds=?, wallet=?, bank=?, credit=?, properties=?, lastvisited=?, metadata=? WHERE user=?", glds, u.Wallet, u.Bank, u.Credit, props, u.LastVisited, met, u.User)
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
		_, err := b.db.Exec("INSERT INTO currency VALUES ( ?, ?, ?, ?, ?, ?, ?, ? )", m.Author.ID, "[\""+m.GuildID+"\"]", 0, 0, 0, "{}", time.Now().Unix(), "{}")
		if b.handle(err, m) {
			return
		}
	} else {
		user, success := b.getuser(m, m.Author.ID)
		if !success {
			return
		}
		isInGuild := false
		for _, guild := range user.Guilds {
			if guild == m.GuildID {
				isInGuild = true
			}
		}
		if !isInGuild {
			user.Guilds = append(user.Guilds, m.GuildID)
			b.updateuser(m, user)
		}
	}
}

func (b *Bot) checkuserwithid(m *discordgo.MessageCreate, id string) {
	exists, success := b.exists(m, "currency", "user=?", id)
	if !success {
		return
	}
	if !exists {
		_, err := b.db.Exec("INSERT INTO currency VALUES ( ?, ?, ?, ?, ?, ?, ?, ? )", id, "[\""+m.GuildID+"\"]", 0, 0, 0, "{}", time.Now().Unix(), "{}")
		if b.handle(err, m) {
			return
		}
	} else {
		user, success := b.getuser(m, id)
		if !success {
			return
		}
		isInGuild := false
		for _, guild := range user.Guilds {
			if guild == m.GuildID {
				isInGuild = true
			}
		}
		if !isInGuild {
			user.Guilds = append(user.Guilds, m.GuildID)
			b.updateuser(m, user)
		}
	}
}

func (b *Bot) abs(val int) int {
	return int(math.Abs(float64(val)))
}

func (b *Bot) isMod(m *discordgo.MessageCreate, ID string) bool {
	mem, err := b.dg.GuildMember(m.GuildID, ID)
	if b.handle(err, m) {
		return false
	}

	// Nv7#0582
	if mem.User.ID == "567132457820749842" {
		return true
	}

	roles, err := b.dg.GuildRoles(m.GuildID)
	if b.handle(err, m) {
		return false
	}
	for _, roleID := range mem.Roles {
		for _, role := range roles {
			if role.ID == roleID && ((role.Permissions & discordgo.PermissionAdministrator) == discordgo.PermissionAdministrator) {
				return true
			}
		}
	}
	return false
}

func (b *Bot) getServerData(m *discordgo.MessageCreate, id string) map[string]interface{} {
	exists, success := b.exists(m, "serverdata", "id=?", id)
	if !success {
		return map[string]interface{}{}
	}
	if !exists {
		_, err := b.db.Exec("INSERT INTO serverdata VALUES ( ?, ? )", id, "{}")
		if b.handle(err, m) {
			return map[string]interface{}{}
		}
	} else {
		row := b.db.QueryRow("SELECT data FROM serverdata WHERE id=?", id)
		var data string
		err := row.Scan(&data)
		if b.handle(err, m) {
			return map[string]interface{}{}
		}
		var out map[string]interface{}
		err = json.Unmarshal([]byte(data), &out)
		if b.handle(err, m) {
			return map[string]interface{}{}
		}
		return out
	}
	return map[string]interface{}{}
}

func (b *Bot) updateServerData(m *discordgo.MessageCreate, id string, data map[string]interface{}) {
	dat, err := json.Marshal(data)
	if b.handle(err, m) {
		return
	}
	_, err = b.db.Exec("UPDATE serverdata SET data=? WHERE id=?", string(dat), id)
	if b.handle(err, m) {
		return
	}
}

func (b *Bot) req(m *discordgo.MessageCreate, url string, out interface{}) bool {
	res, err := http.Get(url)
	if b.handle(err, m) {
		return false
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&out)
	if b.handle(err, m) {
		return false
	}
	return true
}

func (b *Bot) checkprefix(m *discordgo.MessageCreate) {
	_, ex := b.prefixcache[m.GuildID]
	if !ex {
		exists, suc := b.exists(m, "prefixes", "guild=?", m.GuildID)
		if !suc {
			return
		}
		if !exists {
			b.db.Exec("INSERT INTO prefixes VALUES ( ?, ? )", m.GuildID, "")
			b.prefixcache[m.GuildID] = ""
		} else {
			row := b.db.QueryRow("SELECT prefix FROM prefixes WHERE guild=?", m.GuildID)
			var prefix string
			err := row.Scan(&prefix)
			if b.handle(err, m) {
				return
			}
			b.prefixcache[m.GuildID] = prefix
		}
	}
}

func (b *Bot) startsWith(m *discordgo.MessageCreate, cmd string) bool {
	prefix, exists := b.prefixcache[m.GuildID]
	if !exists {
		b.checkprefix(m)
	}
	if strings.HasPrefix(m.Content, prefix+cmd) {
		m.Content = m.Content[len(prefix):]
		return true
	}
	return false
}
