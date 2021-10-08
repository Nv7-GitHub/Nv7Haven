package eod

import (
	"fmt"
	"log"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

const defaultPageLength = 10
const playPageLength = 30

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

func (b *EoD) newPageSwitcher(ps types.PageSwitcher, m types.Msg, rsp types.Rsp) {
	rsp.Acknowledge()
	// Get emojis for guild to find their ID
	/*ems, _ := b.dg.GuildEmojis("819077688371314718")
	for _, emoji := range ems {
		fmt.Println(emoji)
	}*/

	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("Guild isn't setup yet!")
		return
	}

	ps.Guild = m.GuildID
	ps.Page = 0
	ps.PageLength = defaultPageLength
	_, exists = dat.PlayChannels[m.ChannelID]
	if exists {
		ps.PageLength = playPageLength
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

	if dat.PageSwitchers == nil {
		dat.PageSwitchers = make(map[string]types.PageSwitcher)
	}

	dat.SavePageSwitcher(id, ps)

	lock.Lock()
	b.dat[m.GuildID] = dat
	lock.Unlock()
}

func (b *EoD) pageSwitchHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	lock.RLock()
	dat, exists := b.dat[i.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}

	ps, res := dat.GetPageSwitcher(i.Message.ID)
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
		color, _ = b.getColor(i.GuildID, i.Member.User.ID)
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
		log.SetOutput(discordlogs)
		log.Println(err)
	}
	dat.PageSwitchers[i.Message.ID] = ps

	lock.Lock()
	b.dat[i.GuildID] = dat
	lock.Unlock()
}
