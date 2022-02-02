package elements

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/finnbear/moderation"
)

var invalidNames = []string{
	"@everyone",
	"@here",
	"<@",
	"İ",
	"\n",
}

var charReplace = map[rune]rune{
	'’': '\'',
	'‘': '\'',
	'”': '"',
	'“': '"',
}

const maxSuggestionLength = 240

var remove = []string{"\uFE0E", "\uFE0F", "\u200B", "\u200E", "\u200F", "\u2060", "\u2061", "\u2062", "\u2063", "\u2064", "\u2065", "\u2066", "\u2067", "\u2068", "\u2069", "\u206A", "\u206B", "\u206C", "\u206D", "\u206E", "\u206F", "\u3000", "\uFE00", "\uFE01", "\uFE02", "\uFE03", "\uFE04", "\uFE05", "\uFE06", "\uFE07", "\uFE08", "\uFE09", "\uFE0A", "\uFE0B", "\uFE0C", "\uFE0D"}

func (b *Elements) SuggestCmd(suggestion string, autocapitalize bool, m types.Msg, rsp types.Rsp) {
	rsp.Acknowledge()

	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	data, _ := b.GetData(m.GuildID)

	if base.IsFoolsMode && !base.IsFool(suggestion) {
		rsp.ErrorMessage(base.MakeFoolResp(suggestion))
		return
	}

	if autocapitalize && strings.ToLower(suggestion) == suggestion {
		suggestion = util.ToTitle(suggestion)
	}

	if strings.HasPrefix(suggestion, "?") {
		rsp.ErrorMessage(db.Config.LangProperty("ElemNameCannotStartWithQuestionMark"))
		return
	}
	if len(suggestion) >= maxSuggestionLength {
		rsp.ErrorMessage(fmt.Sprintf(db.Config.LangProperty("ElemNameMaxLength"), maxSuggestionLength))
		return
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

	for _, name := range invalidNames {
		if strings.Contains(suggestion, name) {
			rsp.ErrorMessage(fmt.Sprintf(db.Config.LangProperty("ElemNameForbiddenChar"), name))
			return
		}
	}

	suggestion = strings.TrimSpace(suggestion)
	if len(suggestion) > 1 && suggestion[0] == '#' {
		suggestion = suggestion[1:]
	}
	if len(suggestion) == 0 {
		rsp.Resp(db.Config.LangProperty("NoSuggestElemName"))
		return
	}

	// Check if play channel
	db.Config.RLock()
	_, exists := db.Config.PlayChannels[m.ChannelID]
	db.Config.RUnlock()
	if !exists {
		rsp.ErrorMessage(db.Config.LangProperty("MustSuggestInPlayChannel"))
		return
	}

	// Check if exists
	comb, res := data.GetComb(m.Author.ID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	_, res = db.GetCombo(comb.Elems)
	if res.Exists {
		rsp.ErrorMessage(db.Config.LangProperty("ComboHasResult"))
		return
	}

	// Offensive language filter (if enabled by server)
	if moderation.Scan(suggestion).Is(moderation.Offensive) && db.Config.SwearFilter {
		rsp.ErrorMessage(db.Config.LangProperty("NoInappropriateSuggest"))
		return
	}

	// Check if result exists
	el, res := db.GetElementByName(suggestion)
	if res.Exists {
		suggestion = el.Name
	}

	err := b.polls.CreatePoll(types.Poll{
		Channel:   db.Config.VotingChannel,
		Guild:     m.GuildID,
		Kind:      types.PollCombo,
		Suggestor: m.Author.ID,

		PollComboData: &types.PollComboData{
			Elems:  comb.Elems,
			Result: suggestion,
			Exists: res.Exists,
		},
	})
	if rsp.Error(err) {
		return
	}

	txt := "**"
	for _, val := range comb.Elems {
		el, _ := db.GetElement(val)
		txt += el.Name + " + "
	}
	txt = txt[:len(txt)-3]
	if len(comb.Elems) == 1 {
		el, _ := db.GetElement(comb.Elems[0])
		txt += " + " + el.Name
	}
	txt += " = " + suggestion + "**"

	txt = fmt.Sprintf(db.Config.LangProperty("SuggestedElem"), txt)

	if !res.Exists {
		txt += " ✨"
	} else {
		txt += " 🌟"
	}

	id := rsp.Message(txt)
	if res.Exists {
		data.SetMsgElem(id, el.ID)
	}
}
