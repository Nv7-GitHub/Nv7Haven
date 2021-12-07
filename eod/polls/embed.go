package polls

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/bwmarrin/discordgo"
)

func (b *Polls) GetPollEmbed(db *eodb.DB, p types.Poll) (*discordgo.MessageEmbed, error) {
	switch p.Kind {
	case types.PollCombo:
		txt := ""
		elems := p.PollComboData.Elems
		if len(elems) < 1 {
			return nil, errors.New("error: combo must have at least one element")
		}
		for _, val := range elems {
			el, _ := db.GetElement(val)
			txt += el.Name + " + "
		}
		txt = txt[:len(txt)-2]
		res, _ := db.GetElement(p.PollComboData.Result)
		txt += " = " + res.Name

		title := "Element"
		if p.PollComboData.Exists {
			title = "Combination"
		}
		return &discordgo.MessageEmbed{
			Title:       title,
			Description: txt + "\n\n" + "Suggested by <@" + p.Suggestor + ">",
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
		}, nil

	case types.PollSign:
		el, _ := db.GetElement(p.PollSignData.Elem)
		return &discordgo.MessageEmbed{
			Title:       "Sign Note",
			Description: fmt.Sprintf("**%s**\nNew Note: %s\n\nOld Note: %s\n\nSuggested by <@%s>", el.Name, p.PollSignData.NewNote, p.PollSignData.OldNote, p.Suggestor),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
		}, nil

	case types.PollImage:
		el, _ := db.GetElement(p.PollImageData.Elem)
		description := fmt.Sprintf("**%s**\n[New Image](%s)\n[Old Image](%s)\n\nSuggested by <@%s>", el.Name, p.PollImageData.NewImage, p.PollImageData.OldImage, p.Suggestor)
		if p.PollImageData.OldImage == "" {
			description = fmt.Sprintf("**%s**\n[New Image](%s)\n\nSuggested by <@%s>", el.Name, p.PollImageData.NewImage, p.Suggestor)
		}
		return &discordgo.MessageEmbed{
			Title:       "Add Image",
			Description: description,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: p.PollImageData.NewImage,
			},
		}, nil

	case types.PollCategorize, types.PollUnCategorize:
		elems := p.PollCategorizeData.Elems
		names := make([]string, len(elems))
		for i, v := range elems {
			el, _ := db.GetElement(v)
			names[i] = el.Name
		}
		name := "Categorize"
		if p.Kind == types.PollUnCategorize {
			name = "Un-Categorize"
		}
		return &discordgo.MessageEmbed{
			Title:       name,
			Description: fmt.Sprintf("Elements:\n**%s**\n\nCategory: **%s**\n\nSuggested By <@%s>", strings.Join(names, "\n"), p.PollCategorizeData.Category, p.Suggestor),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
		}, nil

	case types.PollCatImage:
		description := fmt.Sprintf("**%s**\n[New Image](%s)\n[Old Image](%s)\n\nSuggested by <@%s>", p.PollCatImageData.Category, p.PollCatImageData.NewImage, p.PollCatImageData.OldImage, p.Suggestor)
		if p.PollCatImageData.OldImage == "" {
			description = fmt.Sprintf("**%s**\n[New Image](%s)\n\nSuggested by <@%s>", p.PollCatImageData.Category, p.PollCatImageData.NewImage, p.Suggestor)
		}
		return &discordgo.MessageEmbed{
			Title:       "Add Category Image",
			Description: description,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: p.PollImageData.NewImage,
			},
		}, nil

	case types.PollColor:
		el, _ := db.GetElement(p.PollColorData.Element)
		return &discordgo.MessageEmbed{
			Title:       "Set Color",
			Description: fmt.Sprintf("**%s**\n%s (Shown on Left)\n\nSuggested by <@%s>", el.Name, util.FormatHex(p.PollColorData.Color), p.Suggestor),
			Color:       p.PollColorData.Color,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
		}, nil

	case types.PollCatColor:
		emb := &discordgo.MessageEmbed{
			Title:       "Reset Category Color",
			Description: fmt.Sprintf("**%s**)\n\nSuggested by <@%s>", p.PollCatColorData.Category, p.Suggestor),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "You can change your vote",
			},
		}
		if p.PollCatColorData.Color != 0 {
			emb = &discordgo.MessageEmbed{
				Title:       "Set Category Color",
				Description: fmt.Sprintf("**%s**\n%s (Shown on Left)\n\nSuggested by <@%s>", p.PollCatColorData.Category, util.FormatHex(p.PollCatColorData.Color), p.Suggestor),
				Color:       p.PollCatColorData.Color,
				Footer: &discordgo.MessageEmbedFooter{
					Text: "You can change your vote",
				},
			}
		}
		return emb, nil
	}

	return nil, errors.New("eod: unknown poll type")
}
