package eod

import (
	"fmt"
	"strings"
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

func (b *EoD) suggestCmd(suggestion string, autocapitalize bool, m msg, rsp rsp) {
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

	_, exists = dat.playChannels[m.ChannelID]
	if !exists {
		rsp.ErrorMessage("You can only suggest in play channels!")
		return
	}

	if dat.combCache == nil {
		dat.combCache = make(map[string]comb)
	}
	comb, exists := dat.combCache[m.Author.ID]
	if !exists {
		rsp.ErrorMessage("You haven't combined anything!")
		return
	}

	data := elems2txt(comb.elems)
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

	dat.lock.RLock()
	el, exists := dat.elemCache[strings.ToLower(suggestion)]
	dat.lock.RUnlock()
	if exists {
		suggestion = el.Name
	}

	err = b.createPoll(poll{
		Channel:   dat.votingChannel,
		Guild:     m.GuildID,
		Kind:      pollCombo,
		Value3:    suggestion,
		Value4:    m.Author.ID,
		Data:      map[string]interface{}{"elems": comb.elems},
		Upvotes:   0,
		Downvotes: 0,
	})
	if rsp.Error(err) {
		return
	}
	txt := "Suggested **"
	for _, val := range comb.elems {
		dat.lock.RLock()
		txt += dat.elemCache[strings.ToLower(val)].Name + " + "
		dat.lock.RUnlock()
	}
	txt = txt[:len(txt)-3]
	if len(comb.elems) == 1 {
		dat.lock.RLock()
		txt += " + " + dat.elemCache[strings.ToLower(comb.elems[0])].Name
		dat.lock.RUnlock()
	}
	txt += " = " + suggestion + "** "

	dat.lock.RLock()
	_, exists = dat.elemCache[strings.ToLower(suggestion)]
	dat.lock.RUnlock()
	if !exists {
		txt += "✨"
	} else {
		txt += "🌟"
	}

	rsp.Message(txt)
}
