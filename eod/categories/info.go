package categories

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/bwmarrin/discordgo"
)

func (b *Categories) progress(els map[int]types.Empty, m types.Msg, db *eodb.DB) *discordgo.MessageEmbedField {
	found := 0
	total := 0
	for el := range els {
		suc, _, tree := trees.CalcElemInfo(el, m.Author.ID, db)
		if !suc {
			continue
		}

		found += tree.Found
		total += tree.Total
	}
	return &discordgo.MessageEmbedField{
		Name:   db.Config.LangProperty("InfoElemProgress", nil),
		Value:  fmt.Sprintf("%s%%", util.FormatFloat(float32(float64(found)/float64(total)*100), 2)),
		Inline: true,
	}
}

func (b *Categories) InfoCmd(catName string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}
	rsp.Acknowledge()

	cat, exists := db.GetCat(catName)
	if exists.Exists {
		// Cat info
		emb := &discordgo.MessageEmbed{
			Title: fmt.Sprintf("%s Info", cat.Name),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: cat.Image,
			},
			Fields: []*discordgo.MessageEmbedField{
				{Name: db.Config.LangProperty("ElementCount", nil), Value: strconv.Itoa(len(cat.Elements))},
				b.progress(cat.Elements, m, db),
			},
			Color: cat.Color,
		}
		if cat.Imager != "" {
			emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{Name: db.Config.LangProperty("InfoImager", nil), Value: fmt.Sprintf("<@%s>", cat.Imager)})
		}
		if cat.Colorer != "" {
			emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{Name: db.Config.LangProperty("InfoColorer", nil), Value: fmt.Sprintf("<@%s>", cat.Colorer)})
		}
		rsp.Embed(emb)
		return
	}

	vcat, res := db.GetVCat(catName)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	els, res := b.base.CalcVCat(vcat, db, true)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	// VCat info
	emb := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s Info", vcat.Name),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: vcat.Image,
		},
		Fields: []*discordgo.MessageEmbedField{
			{Name: db.Config.LangProperty("ElementCount", nil), Value: strconv.Itoa(len(els)), Inline: true},
			{Name: db.Config.LangProperty("InfoCreator", nil), Value: fmt.Sprintf("<@%s>", vcat.Creator), Inline: true},
			b.progress(cat.Elements, m, db),
		},
		Color: vcat.Color,
	}
	if vcat.Imager != "" {
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{Name: db.Config.LangProperty("InfoImager", nil), Value: fmt.Sprintf("<@%s>", vcat.Imager), Inline: true})
	}
	if vcat.Colorer != "" {
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{Name: db.Config.LangProperty("InfoColorer", nil), Value: fmt.Sprintf("<@%s>", vcat.Colorer), Inline: true})
	}

	// TODO: Translate fields below
	emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{Name: "Kind", Value: vcat.Rule.String(), Inline: true}) // TODO: Translate)
	switch vcat.Rule {
	case types.VirtualCategoryRuleRegex:
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{Name: "Regex", Value: "`" + vcat.Data["regex"].(string) + "`", Inline: true})

	case types.VirtualCategoryRuleInvFilter:
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{Name: "User", Value: fmt.Sprintf("<@%s>", vcat.Data["user"].(string)), Inline: true})
		filter := "None"
		if vcat.Data["filter"] == "madeby" {
			filter = "Made By"
		}
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{Name: "Filter", Value: filter, Inline: true})

	case types.VirtualCategoryRuleSetOperation:
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{Name: "Operation", Value: strings.ToTitle(vcat.Data["operation"].(string)), Inline: true})
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{Name: "Left", Value: vcat.Data["lhs"].(string), Inline: true})
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{Name: "Right", Value: vcat.Data["rhs"].(string), Inline: true})
	}
	rsp.Embed(emb)
}
