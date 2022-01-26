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
		txt += " = " + p.PollComboData.Result

		title := db.Config.LangProperty("NewElemPoll")
		if p.PollComboData.Exists {
			title = db.Config.LangProperty("NewComboPoll")
		}
		return &discordgo.MessageEmbed{
			Title:       title,
			Description: txt + "\n\n" + fmt.Sprintf(db.Config.LangProperty("PollCreatorText"), p.Suggestor),
			Footer: &discordgo.MessageEmbedFooter{
				Text: db.Config.LangProperty("PollFooter"),
			},
		}, nil

	case types.PollSign:
		el, _ := db.GetElement(p.PollSignData.Elem)
		return &discordgo.MessageEmbed{
			Title:       db.Config.LangProperty("NewMarkPoll"),
			Description: fmt.Sprintf("**%s**\n%s\n\n%s\n\n", el.Name, fmt.Sprintf(db.Config.LangProperty("NewNote"), p.PollSignData.NewNote), fmt.Sprintf(db.Config.LangProperty("OldNote"), p.PollSignData.OldNote)) + fmt.Sprintf(db.Config.LangProperty("PollCreatorText"), p.Suggestor),
			Footer: &discordgo.MessageEmbedFooter{
				Text: db.Config.LangProperty("PollFooter"),
			},
		}, nil

	case types.PollImage:
		el, _ := db.GetElement(p.PollImageData.Elem)
		description := fmt.Sprintf("**%s**\n[%s](%s)\n[%s](%s)\n\n", el.Name, db.Config.LangProperty("NewImage"), p.PollImageData.NewImage, db.Config.LangProperty("OldImage"), p.PollImageData.OldImage) + fmt.Sprintf(db.Config.LangProperty("PollCreatorText"), p.Suggestor)
		if p.PollImageData.OldImage == "" {
			description = fmt.Sprintf("**%s**\n[%s](%s)\n\n", el.Name, db.Config.LangProperty("NewImage"), p.PollImageData.NewImage) + fmt.Sprintf(db.Config.LangProperty("PollCreatorText"), p.Suggestor)
		}
		return &discordgo.MessageEmbed{
			Title:       db.Config.LangProperty("ElemImagePoll"),
			Description: description,
			Footer: &discordgo.MessageEmbedFooter{
				Text: db.Config.LangProperty("PollFooter"),
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
		name := db.Config.LangProperty("AddCatPoll")
		if p.Kind == types.PollUnCategorize {
			name = db.Config.LangProperty("RmCatPoll")
		}
		return &discordgo.MessageEmbed{
			Title:       name,
			Description: fmt.Sprintf("%s\n\n%s\n\n", fmt.Sprintf(db.Config.LangProperty("CatPollElems"), strings.Join(names, "\n")), fmt.Sprintf(db.Config.LangProperty("CatPollCat"), p.PollCategorizeData.Category)) + fmt.Sprintf(db.Config.LangProperty("PollCreatorText"), p.Suggestor),
			Footer: &discordgo.MessageEmbedFooter{
				Text: db.Config.LangProperty("PollFooter"),
			},
		}, nil

	case types.PollCatImage:
		description := fmt.Sprintf("**%s**\n[%s](%s)\n[%s](%s)\n\n", p.PollCatImageData.Category, db.Config.LangProperty("NewImage"), p.PollCatImageData.NewImage, db.Config.LangProperty("OldImage"), p.PollCatImageData.OldImage) + fmt.Sprintf(db.Config.LangProperty("PollCreatorText"), p.Suggestor)
		if p.PollCatImageData.OldImage == "" {
			description = fmt.Sprintf("**%s**\n[%s](%s)\n\n", p.PollCatImageData.Category, db.Config.LangProperty("NewImage"), p.PollCatImageData.NewImage) + fmt.Sprintf(db.Config.LangProperty("PollCreatorText"), p.Suggestor)
		}
		return &discordgo.MessageEmbed{
			Title:       db.Config.LangProperty("CatImagePoll"),
			Description: description,
			Footer: &discordgo.MessageEmbedFooter{
				Text: db.Config.LangProperty("PollFooter"),
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: p.PollCatImageData.NewImage,
			},
		}, nil

	case types.PollColor:
		el, _ := db.GetElement(p.PollColorData.Element)
		return &discordgo.MessageEmbed{
			Title:       db.Config.LangProperty("ElemColorPoll"),
			Description: fmt.Sprintf("**%s**\n%s %s\n%s\n\n", el.Name, util.FormatHex(p.PollColorData.Color), db.Config.LangProperty("ShownOnLeft"), fmt.Sprintf(db.Config.LangProperty("OldColor"), util.FormatHex(p.PollColorData.OldColor))) + fmt.Sprintf(db.Config.LangProperty("PollCreatorText"), p.Suggestor),
			Color:       p.PollColorData.Color,
			Footer: &discordgo.MessageEmbedFooter{
				Text: db.Config.LangProperty("PollFooter"),
			},
		}, nil

	case types.PollCatColor:
		emb := &discordgo.MessageEmbed{
			Title:       db.Config.LangProperty("ResetCatColorPoll"),
			Description: fmt.Sprintf("**%s**\n\n", p.PollCatColorData.Category) + fmt.Sprintf(db.Config.LangProperty("PollCreatorText"), p.Suggestor),
			Footer: &discordgo.MessageEmbedFooter{
				Text: db.Config.LangProperty("PollFooter"),
			},
		}
		if p.PollCatColorData.Color != 0 {
			txt := fmt.Sprintf("**%s**\n%s %s\n\n", p.PollCatColorData.Category, util.FormatHex(p.PollCatColorData.Color), db.Config.LangProperty("ShownOnLeft")) + fmt.Sprintf(db.Config.LangProperty("PollCreatorText"), p.Suggestor)
			if p.PollCatColorData.OldColor != 0 {
				txt = fmt.Sprintf("**%s**\n%s %s\n%s\n\n", p.PollCatColorData.Category, util.FormatHex(p.PollCatColorData.Color), db.Config.LangProperty("ShownOnLeft"), fmt.Sprintf(db.Config.LangProperty("OldColor"), util.FormatHex(p.PollCatColorData.OldColor))) + fmt.Sprintf(db.Config.LangProperty("PollCreatorText"), p.Suggestor)
			}
			emb = &discordgo.MessageEmbed{
				Title:       db.Config.LangProperty("SetCatColorPoll"),
				Description: txt,
				Color:       p.PollCatColorData.Color,
				Footer: &discordgo.MessageEmbedFooter{
					Text: db.Config.LangProperty("PollFooter"),
				},
			}
		}
		return emb, nil
	}

	return nil, errors.New("eod: unknown poll type")
}
