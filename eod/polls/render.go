package polls

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
)

const footer = "You can change your vote"

const maxCategorizePollElems = 20

func (b *Polls) makePollEmbed(p *types.Poll) (sevcord.EmbedBuilder, error) {
	switch p.Kind {
	case types.PollKindCombo:
		return b.makeComboEmbed(p)

	case types.PollKindImage:
		title := "Add Image"
		oldImage := ""
		if p.Data["old"] != "" {
			oldImage = "[Old Image](" + p.Data["old"].(string) + ")\n"
			title = "Change Image"
		}
		name, err := b.base.GetName(p.Guild, int(p.Data["elem"].(float64)))
		if err != nil {
			return sevcord.NewEmbed(), err
		}
		return sevcord.NewEmbed().
				Title(title).
				Description(makeMessage(fmt.Sprintf("**%s**\n%s[New Image](%s)", name, oldImage, p.Data["new"].(string)), p)).
				Footer(footer, "").
				Thumbnail(p.Data["new"].(string)),
			nil

	case types.PollKindCatImage:
		title := "Add Image"
		oldImage := ""
		if p.Data["old"] != "" {
			oldImage = "[Old Image](" + p.Data["old"].(string) + ")\n"
			title = "Change Image"
		}
		return sevcord.NewEmbed().
				Title(title).
				Description(makeMessage(fmt.Sprintf("**%s**\n%s[New Image](%s)", p.Data["cat"].(string), oldImage, p.Data["new"].(string)), p)).
				Footer(footer, "").
				Thumbnail(p.Data["new"].(string)),
			nil

	case types.PollKindQueryImage:
		title := "Add Image"
		oldImage := ""
		if p.Data["old"] != "" {
			oldImage = "[Old Image](" + p.Data["old"].(string) + ")\n"
			title = "Change Image"
		}
		return sevcord.NewEmbed().
				Title(title).
				Description(makeMessage(fmt.Sprintf("**%s**\n%s[New Image](%s)", p.Data["query"].(string), oldImage, p.Data["new"].(string)), p)).
				Footer(footer, "").
				Thumbnail(p.Data["new"].(string)),
			nil

	case types.PollKindComment:
		name, err := b.base.GetName(p.Guild, int(p.Data["elem"].(float64)))
		if err != nil {
			return sevcord.NewEmbed(), err
		}
		return sevcord.NewEmbed().
			Title("Sign Element").
			Description(makeMessage(fmt.Sprintf("**%s**\nNew Note: %s\n\nOld Note: %s", name, p.Data["new"].(string), p.Data["old"].(string)), p)).
			Footer(footer, ""), nil

	case types.PollKindCatComment:
		return sevcord.NewEmbed().
			Title("Sign Category").
			Description(makeMessage(fmt.Sprintf("**%s**\nNew Note: %s\n\nOld Note: %s", p.Data["cat"].(string), p.Data["new"].(string), p.Data["old"].(string)), p)).
			Footer(footer, ""), nil

	case types.PollKindQueryComment:
		return sevcord.NewEmbed().
			Title("Sign Query").
			Description(makeMessage(fmt.Sprintf("**%s**\nNew Note: %s\n\nOld Note: %s", p.Data["query"].(string), p.Data["new"].(string), p.Data["old"].(string)), p)).
			Footer(footer, ""), nil

	case types.PollKindColor:
		name, err := b.base.GetName(p.Guild, int(p.Data["elem"].(float64)))
		if err != nil {
			return sevcord.NewEmbed(), err
		}
		return sevcord.NewEmbed().
			Title("Set Color").
			Description(makeMessage(fmt.Sprintf("**%s**\nNew Color: %s (shown on left)\n\nOld Color: %s", name, util.FormatHex(int(p.Data["new"].(float64))), util.FormatHex(int(p.Data["old"].(float64)))), p)).
			Color(int(p.Data["new"].(float64))).
			Footer(footer, ""), nil

	case types.PollKindCatColor:
		return sevcord.NewEmbed().
			Title("Set Category Color").
			Description(makeMessage(fmt.Sprintf("**%s**\nNew Color: %s (shown on left)\n\nOld Color: %s", p.Data["cat"].(string), util.FormatHex(int(p.Data["new"].(float64))), util.FormatHex(int(p.Data["old"].(float64)))), p)).
			Color(int(p.Data["new"].(float64))).
			Footer(footer, ""), nil

	case types.PollKindQueryColor:
		return sevcord.NewEmbed().
			Title("Set Query Color").
			Description(makeMessage(fmt.Sprintf("**%s**\nNew Color: %s (shown on left)\n\nOld Color: %s", p.Data["query"].(string), util.FormatHex(int(p.Data["new"].(float64))), util.FormatHex(int(p.Data["old"].(float64)))), p)).
			Color(int(p.Data["new"].(float64))).
			Footer(footer, ""), nil

	case types.PollKindCategorize, types.PollKindUncategorize:
		elems := util.Map(p.Data["elems"].([]any), func(v any) int { return int(v.(float64)) })
		moreTxt := ""
		if len(elems) > maxCategorizePollElems {
			moreTxt = fmt.Sprintf("\nAnd %d more...", len(elems)-maxCategorizePollElems)
			elems = elems[:maxCategorizePollElems]
		}
		names, err := b.base.GetNames(elems, p.Guild)
		if err != nil {
			return sevcord.NewEmbed(), err
		}
		title := "Categorize"
		if p.Kind == types.PollKindUncategorize {
			title = "Un-Categorize"
		}
		return sevcord.NewEmbed().
			Title(title).
			Description(makeMessage(fmt.Sprintf("**%s**\nElements:\n%s%s", p.Data["cat"].(string), strings.Join(names, "\n"), moreTxt), p)).
			Footer(footer, ""), nil

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
	items := util.Map(p.Data["els"].([]any), func(a any) int {
		return int(a.(float64))
	})
	if ok {
		items = append(items, int(res))
	}
	names, err := b.base.GetNames(items, p.Guild)
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
