package elements

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
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

func (b *Elements) SuggestCmd(suggestion string, autocapitalize bool, m types.Msg, rsp types.Rsp) {
	rsp.Acknowledge()

	if base.IsFoolsMode && !base.IsFool(suggestion) {
		rsp.ErrorMessage(base.MakeFoolResp(suggestion))
		return
	}

	if autocapitalize && strings.ToLower(suggestion) == suggestion {
		suggestion = util.ToTitle(suggestion)
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

	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
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

	data := util.Elems2Txt(comb.Elems)
	_, res = dat.GetCombo(data)
	if res.Exists {
		rsp.ErrorMessage("That combo already has a result!")
		return
	}

	el, res := dat.GetElement(suggestion)
	if res.Exists {
		suggestion = el.Name
	}

	err := b.polls.CreatePoll(types.Poll{
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
	}
	txt = txt[:len(txt)-3]
	if len(comb.Elems) == 1 {
		el, _ := dat.GetElement(comb.Elems[0])
		txt += " + " + el.Name
	}
	txt += " = " + suggestion + "** "

	_, res = dat.GetElement(suggestion)
	if !res.Exists {
		txt += "✨"
	} else {
		txt += "🌟"
	}

	id := rsp.Message(txt)
	dat.SetMsgElem(id, suggestion)

	b.lock.Lock()
	b.dat[m.GuildID] = dat
	b.lock.Unlock()
}
