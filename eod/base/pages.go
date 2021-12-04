package base

import (
	"fmt"
	"log"

	"github.com/Nv7-Github/Nv7Haven/eod/logs"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

const DefaultPageLength = 10
const PlayPageLength = 30

var btnRow = discordgo.ActionsRow{
	Components: []discordgo.MessageComponent{
		discordgo.Button{
			Emoji: discordgo.ComponentEmoji{
				Name:     "leftarrow",
				ID:       "861722690813165598",
				Animated: false,
			},
			CustomID: "prev",
		},
		discordgo.Button{
			Emoji: discordgo.ComponentEmoji{
				Name:     "rightarrow",
				ID:       "861722690926936084",
				Animated: false,
			},
			CustomID: "next",
		},
	},
}

func (b *Base) NewPageSwitcher(ps types.PageSwitcher, m types.Msg, rsp types.Rsp) {
	rsp.Acknowledge()
	// Get emojis for guild to find their ID
	/*ems, _ := b.dg.GuildEmojis("819077688371314718")
	for _, emoji := range ems {
		fmt.Println(emoji)
	}*/

	data, res := b.GetData(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}
	db, _ := b.GetDB(m.GuildID)

	ps.Guild = m.GuildID
	ps.Page = 0
	ps.PageLength = DefaultPageLength

	db.Config.RLock()
	_, exists := db.Config.PlayChannels[m.ChannelID]
	db.Config.RUnlock()
	if exists {
		ps.PageLength = PlayPageLength
	}

	cont, _, length, err := ps.PageGetter(ps)
	if rsp.Error(err) {
		return
	}

	footerTxt := fmt.Sprintf("Page %d/%d", ps.Page+1, length+1)
	if ps.Footer != "" {
		footerTxt += " • " + ps.Footer
	}

	id := rsp.Embed(&discordgo.MessageEmbed{
		Title:       ps.Title,
		Description: cont,
		Color:       ps.Color,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: ps.Thumbnail,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: footerTxt,
		},
	}, btnRow)

	data.SavePageSwitcher(id, ps)
}

func (b *Base) PageSwitchHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data, res := b.GetData(i.GuildID)
	if !res.Exists {
		return
	}

	ps, res := data.GetPageSwitcher(i.Message.ID)
	if !res.Exists {
		return
	}

	resp := i.MessageComponentData()
	if resp.CustomID == "next" {
		ps.Page++
	} else if resp.CustomID == "prev" {
		ps.Page--
	} else {
		return
	}

	cont, page, length, err := ps.PageGetter(ps)
	if err != nil {
		return
	}
	if page != ps.Page {
		ps.Page = page
		cont, _, length, err = ps.PageGetter(ps)
		if err != nil {
			return
		}
	}

	footerTxt := fmt.Sprintf("Page %d/%d", ps.Page+1, length+1)
	if ps.Footer != "" {
		footerTxt += " • " + ps.Footer
	}

	color := ps.Color
	if color == 0 {
		color, _ = b.GetColor(i.GuildID, i.Member.User.ID)
	}
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       ps.Title,
					Description: cont,
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: ps.Thumbnail,
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text: footerTxt,
					},
					Color: color,
				},
			},
			Components: []discordgo.MessageComponent{btnRow},
		},
	})
	if err != nil {
		log.SetOutput(logs.DiscordLogs)
		log.Println(err)
	}
	data.SavePageSwitcher(i.Message.ID, ps)
}
