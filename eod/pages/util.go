package pages

import "github.com/Nv7-Github/sevcord/v2"

// Creates btns with name prevnext|<params>
func PageSwitchBtns(handler, params string) []sevcord.Component {
	return []sevcord.Component{
		sevcord.NewButton("", sevcord.ButtonStylePrimary, handler, "prev|"+params).WithEmoji(sevcord.ComponentEmojiCustom("leftarrow", "861722690813165598", false)),
		sevcord.NewButton("", sevcord.ButtonStylePrimary, handler, "next|"+params).WithEmoji(sevcord.ComponentEmojiCustom("rightarrow", "861722690926936084", false)),
	}
}
