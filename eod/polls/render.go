package polls

import (
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
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
	var names []struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}
	err := b.db.Select(&names, "SELECT id, name FROM elements WHERE id = ANY($1)", pq.Array(items))
	if err != nil {
		return sevcord.NewEmbed(), err
	}
	nameMap := make(map[int]string)
	for _, v := range names {
		nameMap[v.ID] = v.Name
	}

	// Generate text
	txt := &strings.Builder{}
	for i, v := range items {
		if i > 0 {
			txt.WriteString(" + ")
		}
		txt.WriteString(nameMap[v])
	}
	txt.WriteString(" = ")
	if ok {
		txt.WriteString(nameMap[int(res)])
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
