package eod

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
)

var invalidNames = []string{
	"+",
	"@everyone",
	"@here",
	"<@",
	"İ",
	"\n",
}

var charReplace = map[rune]rune{
	'’': '\'',
	'‘': '\'',
	'`': '\'',
	'”': '"',
	'“': '"',
}

const maxSuggestionLength = 240

var remove = []string{"\uFE0E", "\uFE0F", "\u200B", "\u200E", "\u200F", "\u2060", "\u2061", "\u2062", "\u2063", "\u2064", "\u2065", "\u2066", "\u2067", "\u2068", "\u2069", "\u206A", "\u206B", "\u206C", "\u206D", "\u206E", "\u206F", "\u3000", "\uFE00", "\uFE01", "\uFE02", "\uFE03", "\uFE04", "\uFE05", "\uFE06", "\uFE07", "\uFE08", "\uFE09", "\uFE0A", "\uFE0B", "\uFE0C", "\uFE0D"}

func (b *EoD) suggestCmd(suggestion string, autocapitalize bool, m types.Msg, rsp types.Rsp) {
	rsp.Acknowledge()

	if isFoolsMode && !isFool(suggestion) {
		rsp.ErrorMessage(makeFoolResp(suggestion))
		return
	}

	if autocapitalize && strings.ToLower(suggestion) == suggestion {
		suggestion = toTitle(suggestion)
	}

	if strings.HasPrefix(suggestion, "?") {
		rsp.ErrorMessage("Element names can't start with '?'!")
		return
	}
	if len(suggestion) >= maxSuggestionLength {
		rsp.ErrorMessage(fmt.Sprintf("Element names must be under %d characters!", maxSuggestionLength))
		return
	}
	for _, name := range invalidNames {
		if strings.Contains(suggestion, name) {
			rsp.ErrorMessage(fmt.Sprintf("Can't have letters '%s' in an element name!", name))
			return
		}
	}

	// Clean up suggestions with weird quotes
	cleaned := []rune(suggestion)
	for i, char := range cleaned {
		newVal, exists := charReplace[char]
		if exists {
			cleaned[i] = newVal
		}
	}
	suggestion = string(cleaned)
	for _, val := range remove {
		suggestion = strings.ReplaceAll(suggestion, val, "")
	}

	suggestion = strings.TrimSpace(suggestion)
	if len(suggestion) > 1 && suggestion[0] == '#' {
		suggestion = suggestion[1:]
	}
	if len(suggestion) == 0 {
		rsp.Resp("You need to suggest something!")
		return
	}

	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("Guild not set up!")
		return
	}

	_, exists = dat.PlayChannels[m.ChannelID]
	if !exists {
		rsp.ErrorMessage("You can only suggest in play channels!")
		return
	}

	comb, res := dat.GetComb(m.Author.ID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	data := elems2txt(comb.Elems)
	query := "SELECT COUNT(1) FROM eod_combos WHERE guild=? AND elems LIKE ?"

	if isASCII(data) {
		query = "SELECT COUNT(1) FROM eod_combos WHERE guild=CONVERT(? USING utf8mb4) AND CONVERT(elems USING utf8mb4) LIKE CONVERT(? USING utf8mb4) COLLATE utf8mb4_general_ci"
	}

	if isWildcard(data) {
		query = strings.ReplaceAll(query, " LIKE ", "=")
	}

	row := b.db.QueryRow(query, m.GuildID, data)
	var count int
	err := row.Scan(&count)
	if rsp.Error(err) {
		return
	}
	if count != 0 {
		rsp.ErrorMessage("That combo already has a result!")
		return
	}

	el, res := dat.GetElement(suggestion)
	if res.Exists {
		suggestion = el.Name
	}

	err = b.createPoll(types.Poll{
		Channel:   dat.VotingChannel,
		Guild:     m.GuildID,
		Kind:      types.PollCombo,
		Value3:    suggestion,
		Value4:    m.Author.ID,
		Data:      map[string]interface{}{"elems": comb.Elems},
		Upvotes:   0,
		Downvotes: 0,
	})
	if rsp.Error(err) {
		return
	}
	txt := "Suggested **"
	for _, val := range comb.Elems {
		el, _ := dat.GetElement(val)
		txt += el.Name + " + "
		dat.Lock.RUnlock()
	}
	txt = txt[:len(txt)-3]
	if len(comb.Elems) == 1 {
		el, _ := dat.GetElement(comb.Elems[0])
		txt += " + " + el.Name
	}
	txt += " = " + suggestion + "** "

	_, res = dat.GetElement(suggestion)
	dat.Lock.RUnlock()
	if !res.Exists {
		txt += "✨"
	} else {
		txt += "🌟"
	}

	rsp.Message(txt)
}
