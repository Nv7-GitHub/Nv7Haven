package polls

import (
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
)

const footer = "You can change your vote"

func (b *Polls) makePollEmbed(p *types.Poll) (sevcord.EmbedBuilder, error) {
	switch p.Kind {
	case types.PollKindCombo:
		return b.makeComboEmbed(p)

		// TODO: The rest

	default:
		return sevcord.NewEmbed(), nil // Impossible
	}
}

func (b *Polls) makeComboEmbed(p *types.Poll) (sevcord.EmbedBuilder, error) {
	// Get title
	title := "Element"
	res, ok := p.Data["result"].(float64)
	if ok {
		title = "Combination"
	}

	// Get list of element names to fetch
	items := util.Map(p.Data["els"].([]float64), func(a float64) int {
		return int(a)
	})
	if ok {
		items = append(items, int(res))
	}
	names, err := b.base.GetNames(items)
	if err != nil {
		return sevcord.NewEmbed(), err
	}
	if ok {
		items = items[:len(items)-1]
	}

	// Generate text
	txt := &strings.Builder{}
	for i := range items {
		if i > 0 {
			txt.WriteString(" + ")
		}
		txt.WriteString(names[i])
	}
	txt.WriteString(" = ")
	if ok {
		txt.WriteString(names[len(names)-1])
	} else {
		txt.WriteString(p.Data["result"].(string))
	}

	return sevcord.NewEmbed().
		Title(title).
		Description(makeMessage(txt.String(), p)).
		Footer(footer, ""), nil
}

func makeMessage(description string, p *types.Poll) string {
	return description + "\n\nSuggested By <@" + p.Creator + ">"
}
