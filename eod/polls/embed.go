package polls

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/bwmarrin/discordgo"
)

func (b *Polls) GetPollEmbed(dat types.ServerDat, p types.OldPoll) (*discordgo.MessageEmbed, error) {
	switch p.Kind {
	case types.PollCombo:
		txt := ""
		elems, ok := p.Data["elems"].([]string)
		if !ok {
			elemDat := p.Data["elems"].([]interface{})
			elems = make([]string, len(elemDat))
			for i, val := range elemDat {
				elems[i] = val.(string)
			}
		}
		if len(elems) < 1 {
			return nil, errors.New("error: combo must have at least one element")
		}
		for _, val := range elems {
			el, _ := dat.GetElement(val)
			txt += el.Name + " + "
		}
		txt = txt[:len(txt)-2]
		if len(elems) == 1 {
			el, _ := dat.GetElement(elems[0])
			txt += " + " + el.Name
		}
		txt += " = " + p.Value3

		title := "Element"
		if p.Data["exists"].(bool) {
			title = "Combination"
		}
		return &discordgo.MessageEmbed{
			Title:       title,
			Description: txt + "\n\n" + "Suggested by <@" + p.Value4 + ">",
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
		}, nil

	case types.PollSign:
		return &discordgo.MessageEmbed{
			Title:       "Sign Note",
			Description: fmt.Sprintf("**%s**\nNew Note: %s\n\nOld Note: %s\n\nSuggested by <@%s>", p.Value1, p.Value2, p.Value3, p.Value4),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
		}, nil

	case types.PollImage:
		description := fmt.Sprintf("**%s**\n[New Image](%s)\n[Old Image](%s)\n\nSuggested by <@%s>", p.Value1, p.Value2, p.Value3, p.Value4)
		if p.Value3 == "" {
			description = fmt.Sprintf("**%s**\n[New Image](%s)\n\nSuggested by <@%s>", p.Value1, p.Value2, p.Value4)
		}
		return &discordgo.MessageEmbed{
			Title:       "Add Image",
			Description: description,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: p.Value2,
			},
		}, nil

	case types.PollCategorize:
		data, ok := p.Data["elems"].([]string)
		if !ok {
			dat := p.Data["elems"].([]interface{})
			data = make([]string, len(dat))
			for i, val := range dat {
				data[i] = val.(string)
			}
		}
		p.Data["elems"] = data
		return &discordgo.MessageEmbed{
			Title:       "Categorize",
			Description: fmt.Sprintf("Elements:\n**%s**\n\nCategory: **%s**\n\nSuggested By <@%s>", strings.Join(p.Data["elems"].([]string), "\n"), p.Value1, p.Value4),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
		}, nil
	case types.PollUnCategorize:
		data, ok := p.Data["elems"].([]string)
		if !ok {
			dat := p.Data["elems"].([]interface{})
			data = make([]string, len(dat))
			for i, val := range dat {
				data[i] = val.(string)
			}
		}
		p.Data["elems"] = data
		return &discordgo.MessageEmbed{
			Title:       "Un-Categorize",
			Description: fmt.Sprintf("Elements:\n**%s**\n\nCategory: **%s**\n\nSuggested By <@%s>", strings.Join(p.Data["elems"].([]string), "\n"), p.Value1, p.Value4),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
		}, nil

	case types.PollCatImage:
		description := fmt.Sprintf("**%s**\n[New Image](%s)\n[Old Image](%s)\n\nSuggested by <@%s>", p.Value1, p.Value2, p.Value3, p.Value4)
		if p.Value3 == "" {
			description = fmt.Sprintf("**%s**\n[New Image](%s)\n\nSuggested by <@%s>", p.Value1, p.Value2, p.Value4)
		}
		return &discordgo.MessageEmbed{
			Title:       "Add Category Image",
			Description: description,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: p.Value2,
			},
		}, nil

	case types.PollColor:
		return &discordgo.MessageEmbed{
			Title:       "Set Color",
			Description: fmt.Sprintf("**%s**\n%s (Shown on Left)\n\nSuggested by <@%s>", p.Value1, util.FormatHex(p.Data["color"].(int)), p.Value4),
			Color:       p.Data["color"].(int),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
		}, nil

	case types.PollCatColor:
		emb := &discordgo.MessageEmbed{
			Title:       "Reset Category Color",
			Description: fmt.Sprintf("**%s**)\n\nSuggested by <@%s>", p.Value1, p.Value4),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
		}
		if p.Data["color"].(int) != 0 {
			emb = &discordgo.MessageEmbed{
				Title:       "Set Category Color",
				Description: fmt.Sprintf("**%s**\n%s (Shown on Left)\n\nSuggested by <@%s>", p.Value1, util.FormatHex(p.Data["color"].(int)), p.Value4),
				Color:       p.Data["color"].(int),
				Footer: &discordgo.MessageEmbedFooter{
					Text: "You can change your vote",
				},
			}
		}
		return emb, nil
	}

	return nil, errors.New("eod: unknown poll type")
}
