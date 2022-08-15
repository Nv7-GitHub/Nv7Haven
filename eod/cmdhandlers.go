package eod

import (
	"strconv"
	"strings"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/bwmarrin/discordgo"
)

var combs = []string{
	"\n",
	"+",
	",",
	"plus",
}

func (b *EoD) addCmdCounter(m *discordgo.MessageCreate, cmd string) {
	db, res := b.GetDB(m.GuildID)
	if res.Exists {
		if db.Config.CommandStats == nil {
			db.Config.CommandStats = make(map[string]int)
		}
		db.Config.CommandStats[cmd]++
		db.SaveConfig()
	}
}

func (b *EoD) cmdHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := b.newMsgNormal(m)
	rsp := b.newRespNormal(m)

	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if strings.HasPrefix(m.Content, "=") {
		if len(m.Content) < 2 {
			return
		}
		suggestion := m.Content[1:]

		suggestion = strings.TrimSpace(strings.ReplaceAll(suggestion, "\n", ""))
		b.addCmdCounter(m, "suggest")
		b.elements.SuggestCmd(suggestion, true, msg, rsp)
		return
	}

	if strings.HasPrefix(m.Content, "+") && len(m.Content) > 1 {
		m.Content = strings.TrimSpace(m.Content[1:])
		if !b.base.CheckServer(msg, rsp) {
			return
		}

		data, res := b.GetData(msg.GuildID)
		if !res.Exists {
			return
		}
		db, res := b.GetDB(msg.GuildID)
		if !res.Exists {
			return
		}

		// Get last
		comb, res := data.GetComb(msg.Author.ID)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
		if comb.Elem3 == -1 {
			txt := make([]string, len(comb.Elems))
			for i, elem := range comb.Elems {
				el, _ := db.GetElement(elem)
				txt[i] = el.Name
			}
			b.basecmds.Combine(txt, msg, rsp)
			return
		}
		el, _ := db.GetElement(comb.Elem3)

		for _, comb := range combs {
			if strings.Contains(m.Content, comb) {

				// Get parts
				parts := strings.Split(m.Content, comb)
				if len(parts) < 2 {
					return
				}
				for i, part := range parts {
					parts[i] = strings.TrimSpace(strings.Replace(part, "\\", "", -1))
				}

				b.addCmdCounter(m, "combine")
				b.basecmds.Combine(append([]string{el.Name}, parts...), msg, rsp)
				return
			}
		}
		b.addCmdCounter(m, "combine")
		b.basecmds.Combine([]string{el.Name, m.Content}, msg, rsp)
		return
	}

	if strings.HasPrefix(m.Content, "!") {
		db, res := b.GetDB(msg.GuildID)
		if !res.Exists {
			return
		}

		if len(m.Content) < 2 {
			return
		}

		cmd := strings.ToLower(strings.Split(m.Content[1:], " ")[0])
		if cmd == "s" || cmd == "suggest" {
			if len(m.Content) <= len(cmd)+2 {
				return
			}
			suggestion := m.Content[len(cmd)+2:]

			suggestion = strings.TrimSpace(strings.ReplaceAll(suggestion, "\n", ""))
			b.addCmdCounter(m, "suggest")
			b.elements.SuggestCmd(suggestion, true, msg, rsp)
			return
		}

		if cmd == "stats" {
			b.basecmds.StatsCmd(msg, rsp)
			return
		}

		if cmd == "image" || cmd == "img" || cmd == "pic" {
			if len(m.Content) <= len(cmd)+2 {
				return
			}
			suggestion := m.Content[len(cmd)+2:]
			suggestion = strings.TrimSpace(strings.ReplaceAll(suggestion, "\n", ""))

			if len(m.Attachments) < 1 {
				rsp.ErrorMessage(db.Config.LangProperty("MustAttachImage", nil))
				return
			}
			b.addCmdCounter(m, "image")
			b.polls.ImageCmd(suggestion, m.Attachments[0].URL, msg, rsp)
			return
		}

		if cmd == "catimage" || cmd == "catimg" || cmd == "catpic" {
			if len(m.Content) <= len(cmd)+2 {
				return
			}
			suggestion := m.Content[len(cmd)+2:]
			suggestion = strings.TrimSpace(strings.ReplaceAll(suggestion, "\n", ""))

			if len(m.Attachments) < 1 {
				rsp.ErrorMessage(db.Config.LangProperty("MustAttachImage", nil))
				return
			}
			b.addCmdCounter(m, "image")
			b.polls.CatImgCmd(suggestion, m.Attachments[0].URL, msg, rsp)
			return
		}

		if cmd == "hint" || cmd == "h" {
			if len(m.Content) <= len(cmd)+2 {
				b.elements.HintCmd("", false, false, false, msg, rsp)
				return
			}
			suggestion := m.Content[len(cmd)+2:]
			suggestion = strings.TrimSpace(strings.ReplaceAll(suggestion, "\n", ""))

			b.addCmdCounter(m, "hint")
			b.elements.HintCmd(suggestion, true, false, false, msg, rsp)
			return
		}

		if cmd == "invhint" || cmd == "ih" {
			if len(m.Content) <= len(cmd)+2 {
				return
			}
			suggestion := m.Content[len(cmd)+2:]
			suggestion = strings.TrimSpace(strings.ReplaceAll(suggestion, "\n", ""))

			b.addCmdCounter(m, "invhint")
			b.elements.HintCmd(suggestion, true, false, true, msg, rsp)
			return
		}

		if cmd == "addcat" || cmd == "ac" {
			if len(m.Content) <= len(cmd)+2 {
				return
			}
			txt := m.Content[len(cmd)+2:]
			sepPos := strings.Index(txt, "|")
			if sepPos == -1 {
				rsp.ErrorMessage(db.Config.LangProperty("AddCatMustHaveSeparator", nil))
				return
			}

			catName := strings.TrimSpace(txt[:sepPos])
			elems := util.TrimArray(splitByCombs(txt[sepPos+1:]))

			b.addCmdCounter(m, "addcat")
			b.categories.CategoryCmd(elems, catName, msg, rsp)
			return
		}

		if cmd == "rmcat" || cmd == "rc" {
			if len(m.Content) <= len(cmd)+2 {
				return
			}
			txt := m.Content[len(cmd)+2:]
			sepPos := strings.Index(txt, "|")
			if sepPos == -1 {
				rsp.ErrorMessage(db.Config.LangProperty("RmCatMustHaveSeparator", nil))
				return
			}

			catName := strings.TrimSpace(txt[:sepPos])
			elems := util.TrimArray(splitByCombs(txt[sepPos+1:]))

			b.addCmdCounter(m, "rmcat")
			b.categories.RmCategoryCmd(elems, catName, msg, rsp)
			return
		}

		if cmd == "inv" || cmd == "elements" {
			b.addCmdCounter(m, "inv")
			b.elements.InvCmd(m.Author.ID, msg, rsp, "name", "none", false, true)
			return
		}

		if cmd == "lb" || cmd == "top" || cmd == "leaderboard" {
			b.addCmdCounter(m, "lb")
			b.elements.LbCmd(msg, rsp, "count", msg.Author.ID, "")
			return
		}

		if cmd == "cat" || cmd == "🐱" {
			if len(m.Content) <= len(cmd)+2 {
				bot.categories.AllCatCmd("name", false, "", msg, rsp)
				return
			}
			suggestion := m.Content[len(cmd)+2:]
			suggestion = strings.TrimSpace(strings.ReplaceAll(suggestion, "\n", ""))

			b.addCmdCounter(m, "cat")
			b.categories.CatCmd(suggestion, "name", false, "", true, msg, rsp)
			return
		}

		if cmd == "mark" || cmd == "sign" {
			if len(m.Content) <= len(cmd)+2 {
				return
			}
			txt := m.Content[len(cmd)+2:]
			sepPos := strings.Index(txt, "|")
			if sepPos == -1 {
				rsp.ErrorMessage(db.Config.LangProperty("MarkMustHaveSeparator", nil))
				return
			}

			elem := strings.TrimSpace(txt[:sepPos])
			mark := strings.TrimSpace(txt[sepPos+1:])
			b.addCmdCounter(m, "mark")
			b.polls.MarkCmd(elem, mark, msg, rsp)
			return
		}

		if cmd == "info" || cmd == "get" {
			if len(m.Content) <= len(cmd)+2 {
				return
			}
			b.addCmdCounter(m, "info")
			b.elements.InfoCmd(strings.TrimSpace(m.Content[len(cmd)+2:]), msg, rsp)
			return
		}

		if cmd == "ping" {
			// Ping command for text
			b.addCmdCounter(m, "ping")
			rsp.Message(db.Config.LangProperty("PingMessage", time.Since(m.Timestamp).String()))
		}

		if cmd == "optimize" {
			if m.GuildID == "705084182673621033" {
				user, err := b.dg.GuildMember(msg.GuildID, msg.Author.ID)
				if rsp.Error(err) {
					return
				}
				for _, roleID := range user.Roles {
					if roleID == "843348517220777984" {
						b.addCmdCounter(m, "optimize")
						b.optimize(msg, rsp)
					}
				}
			}
		}
	}

	if strings.HasPrefix(m.Content, "?") {
		if len(m.Content) < 2 {
			return
		}
		b.addCmdCounter(m, "info")
		b.elements.InfoCmd(strings.TrimSpace(m.Content[1:]), msg, rsp)
		return
	}

	if strings.HasPrefix(m.Content, "*") && len(m.Content) > 1 {
		if !b.base.CheckServer(msg, rsp) {
			return
		}
		data, res := b.GetData(msg.GuildID)
		if !res.Exists {
			return
		}
		db, res := b.GetDB(msg.GuildID)
		if !res.Exists {
			return
		}

		cont := m.Content[1:]
		split := strings.Contains(cont, " ")
		var items []string
		if split {
			pos := strings.Index(cont, " ")
			items = []string{cont[:pos], cont[pos+1:]}
		} else {
			items = []string{cont}
		}

		length, err := strconv.Atoi(items[0])
		if err != nil {
			return
		}
		if length > types.MaxComboLength {
			length = types.MaxComboLength + 1 // This way it triggers the error message in the combo command
		}
		if length < 2 {
			length = 1
		}

		var last string
		if split {
			last = items[1]
		} else {
			comb, res := data.GetComb(msg.Author.ID)
			if !res.Exists {
				rsp.ErrorMessage(res.Message)
				return
			}

			if comb.Elem3 == -1 {
				txt := make([]string, len(comb.Elems))
				for i, elem := range comb.Elems {
					el, _ := db.GetElement(elem)
					txt[i] = el.Name
				}
				b.addCmdCounter(m, "combine")
				b.basecmds.Combine(txt, msg, rsp)
				return
			}

			el, _ := db.GetElement(comb.Elem3)
			last = el.Name
		}

		elems := make([]string, length)
		for i := range elems {
			elems[i] = last
		}

		b.addCmdCounter(m, "combine")
		b.basecmds.Combine(elems, msg, rsp)
		return
	}

	for _, comb := range combs {
		if strings.Contains(m.Content, comb) {
			if !b.base.CheckServer(msg, rsp) {
				return
			}
			parts := strings.Split(m.Content, comb)
			if len(parts) < 2 {
				return
			}
			for i, part := range parts {
				parts[i] = strings.TrimSpace(strings.Replace(part, "\\", "", -1))
			}
			b.addCmdCounter(m, "combine")
			b.basecmds.Combine(parts, msg, rsp)
			return
		}
	}
}
