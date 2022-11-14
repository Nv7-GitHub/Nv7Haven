package pages

import "github.com/Nv7-Github/sevcord/v2"

// Creates btns with name prevnext|<params>
func PageSwitchBtns(handler, params string) []sevcord.Component {
	return []sevcord.Component{
		sevcord.NewButton("", sevcord.ButtonStylePrimary, handler, "prev|"+params).WithEmoji(sevcord.ComponentEmojiCustom("leftarrow", "861722690813165598", false)),
		sevcord.NewButton("", sevcord.ButtonStylePrimary, handler, "next|"+params).WithEmoji(sevcord.ComponentEmojiCustom("rightarrow", "861722690926936084", false)),
	}
}

func ApplyPage(param string, page, pagecnt int) int {
	switch param {
	case "prev":
		page--

	case "next":
		page++
	}
	if page < 0 {
		page = pagecnt - 1
	}
	if page >= pagecnt {
		page = 0
	}
	return page
}
