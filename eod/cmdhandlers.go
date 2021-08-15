package eod

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const maxComboLength = 20

var combs = []string{
	"\n",
	"+",
	",",
	"plus",
}

func (b *EoD) cmdHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := b.newMsgNormal(m)
	rsp := b.newRespNormal(m)

	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if strings.HasPrefix(m.Content, "+") {
		if len(m.Content) < 2 {
			return
		}
		suggestion := m.Content[1:]

		suggestion = strings.TrimSpace(strings.ReplaceAll(suggestion, "\n", ""))
		b.suggestCmd(suggestion, true, msg, rsp)
		return
	}

	if strings.HasPrefix(m.Content, "!") {
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
			b.suggestCmd(suggestion, true, msg, rsp)
			return
		}

		if cmd == "stats" {
			b.statsCmd(msg, rsp)
			return
		}

		if cmd == "image" || cmd == "img" || cmd == "pic" {
			if len(m.Content) <= len(cmd)+2 {
				return
			}
			suggestion := m.Content[len(cmd)+2:]
			suggestion = strings.TrimSpace(strings.ReplaceAll(suggestion, "\n", ""))

			if len(m.Attachments) < 1 {
				rsp.ErrorMessage("You must attach an image!")
				return
			}
			b.imageCmd(suggestion, m.Attachments[0].URL, msg, rsp)
			return
		}

		if cmd == "catimage" || cmd == "catimg" || cmd == "catpic" {
			if len(m.Content) <= len(cmd)+2 {
				return
			}
			suggestion := m.Content[len(cmd)+2:]
			suggestion = strings.TrimSpace(strings.ReplaceAll(suggestion, "\n", ""))

			if len(m.Attachments) < 1 {
				rsp.ErrorMessage("You must attach an image!")
				return
			}
			b.catImgCmd(suggestion, m.Attachments[0].URL, msg, rsp)
			return
		}

		if cmd == "hint" || cmd == "h" {
			if len(m.Content) <= len(cmd)+2 {
				b.hintCmd("", false, false, msg, rsp)
				return
			}
			suggestion := m.Content[len(cmd)+2:]
			suggestion = strings.TrimSpace(strings.ReplaceAll(suggestion, "\n", ""))

			b.hintCmd(suggestion, true, false, msg, rsp)
			return
		}

		if cmd == "addcat" || cmd == "ac" {
			if len(m.Content) <= len(cmd)+2 {
				return
			}
			txt := m.Content[len(cmd)+2:]
			sepPos := strings.Index(txt, "|")
			if sepPos == -1 {
				rsp.ErrorMessage("You must have a \"|\" to seperate the category name and the elements to add!")
				return
			}

			catName := strings.TrimSpace(txt[:sepPos])
			elems := strings.Split(txt[sepPos+1:], ",")
			for i, elem := range elems {
				elems[i] = strings.TrimSpace(elem)
			}

			b.categoryCmd(elems, catName, msg, rsp)
			return
		}

		if cmd == "rmcat" || cmd == "rc" {
			if len(m.Content) <= len(cmd)+2 {
				return
			}
			txt := m.Content[len(cmd)+2:]
			sepPos := strings.Index(txt, "|")
			if sepPos == -1 {
				rsp.ErrorMessage("You must have a \"|\" to seperate the category name and the elements to remove!")
				return
			}

			catName := strings.TrimSpace(txt[:sepPos])
			elems := strings.Split(txt[sepPos+1:], ",")
			for i, elem := range elems {
				elems[i] = strings.TrimSpace(elem)
			}

			b.rmCategoryCmd(elems, catName, msg, rsp)
			return
		}

		if cmd == "inv" {
			b.invCmd(m.Author.ID, msg, rsp, "")
			return
		}

		if cmd == "lb" {
			b.lbCmd(msg, rsp, "count")
			return
		}

		if cmd == "cat" {
			if len(m.Content) <= len(cmd)+2 {
				bot.allCatCmd(catSortAlphabetical, false, "", msg, rsp)
				return
			}
			suggestion := m.Content[len(cmd)+2:]
			suggestion = strings.TrimSpace(strings.ReplaceAll(suggestion, "\n", ""))

			b.catCmd(suggestion, catSortAlphabetical, false, "", msg, rsp)
			return
		}

		if cmd == "mark" || cmd == "sign" {
			if len(m.Content) <= len(cmd)+2 {
				return
			}
			txt := m.Content[len(cmd)+2:]
			sepPos := strings.Index(txt, "|")
			if sepPos == -1 {
				rsp.ErrorMessage("You must have a \"|\" to seperate element name and its new mark!")
				return
			}

			elem := strings.TrimSpace(txt[:sepPos])
			mark := strings.TrimSpace(txt[sepPos+1:])
			b.markCmd(elem, mark, msg, rsp)
			return
		}
		if cmd == "info" || cmd == "get" {
			if len(m.Content) <= len(cmd)+2 {
				return
			}
			b.infoCmd(strings.TrimSpace(m.Content[len(cmd)+2:]), msg, rsp)
			return
		}
	}

	if strings.HasPrefix(m.Content, "?") {
		if len(m.Content) < 2 {
			return
		}
		b.infoCmd(strings.TrimSpace(m.Content[1:]), msg, rsp)
		return
	}

	if strings.HasPrefix(m.Content, "*2") {
		if !b.checkServer(msg, rsp) {
			return
		}
		lock.RLock()
		dat, exists := b.dat[msg.GuildID]
		lock.RUnlock()
		if !exists {
			return
		}

		comb, res := dat.GetComb(msg.Author.ID)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
			return
		}
		if comb.Elem3 != "" {
			b.combine([]string{comb.Elem3, comb.Elem3}, msg, rsp)
			return
		}
		b.combine(comb.Elems, msg, rsp)
		return
	}

	for _, comb := range combs {
		if strings.Contains(m.Content, comb) {
			if !b.checkServer(msg, rsp) {
				return
			}
			parts := strings.Split(m.Content, comb)
			if len(parts) < 2 {
				return
			}
			for i, part := range parts {
				parts[i] = strings.TrimSpace(strings.Replace(part, "\\", "", -1))
			}
			if len(parts) > maxComboLength {
				rsp.ErrorMessage(fmt.Sprintf("You can only combine up to %d elements!", maxComboLength))
				return
			}
			b.combine(parts, msg, rsp)
			return
		}
	}
}
