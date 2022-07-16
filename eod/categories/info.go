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

const maxprogresslen = 25000

func (b *Categories) progress(els map[int]types.Empty, m types.Msg, db *eodb.DB) *discordgo.MessageEmbedField {
	tree := trees.CalcCatInfo(els, m.Author.ID, db)
	return &discordgo.MessageEmbedField{
		Name:   db.Config.LangProperty("InfoElemProgress", nil),
		Value:  fmt.Sprintf("%s%%", util.FormatFloat(float32(float64(tree.Found)/float64(tree.Total)*100), 2)),
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
		if cat.Comment == "" {
			cat.Comment = "None"
		}

		// Cat info
		emb := &discordgo.MessageEmbed{
			Title: fmt.Sprintf("%s Info", cat.Name),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: cat.Image,
			},
			Fields: []*discordgo.MessageEmbedField{
				{Name: db.Config.LangProperty("InfoComment", nil), Value: cat.Comment, Inline: false},
				{Name: db.Config.LangProperty("ElementCount", nil), Value: strconv.Itoa(len(cat.Elements)), Inline: true},
			},
			Color: cat.Color,
		}
		if len(cat.Elements) < maxprogresslen {
			emb.Fields = append(emb.Fields, b.progress(cat.Elements, m, db))
		}
		if cat.Commenter != "" {
			emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{Name: db.Config.LangProperty("InfoCommenter", nil), Value: fmt.Sprintf("<@%s>", cat.Commenter), Inline: true})
		}
		if cat.Imager != "" {
			emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{Name: db.Config.LangProperty("InfoImager", nil), Value: fmt.Sprintf("<@%s>", cat.Imager), Inline: true})
		}
		if cat.Colorer != "" {
			emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{Name: db.Config.LangProperty("InfoColorer", nil), Value: fmt.Sprintf("<@%s>", cat.Colorer), Inline: true})
		}
		rsp.Embed(emb)
		return
	}

	vcat, res := db.GetVCat(catName)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	if vcat.Comment == "" {
		vcat.Comment = "None"
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
			{Name: db.Config.LangProperty("InfoComment", nil), Value: vcat.Comment, Inline: false},
			{Name: db.Config.LangProperty("ElementCount", nil), Value: strconv.Itoa(len(els)), Inline: true},
			{Name: db.Config.LangProperty("InfoCreator", nil), Value: fmt.Sprintf("<@%s>", vcat.Creator), Inline: true},
		},
		Color: vcat.Color,
	}
	if len(els) < maxprogresslen {
		emb.Fields = append(emb.Fields, b.progress(els, m, db))
	}
	if vcat.Commenter != "" {
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{Name: db.Config.LangProperty("InfoCommenter", nil), Value: fmt.Sprintf("<@%s>", vcat.Commenter), Inline: true})
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
