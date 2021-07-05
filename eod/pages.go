package eod

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

const defaultPageLength = 10
const playPageLength = 30

var btnRow = discordgo.ActionsRow{
	Components: []discordgo.MessageComponent{
		discordgo.Button{
			Label:    "Previous Page",
			CustomID: "prev",
		},
		discordgo.Button{
			Label:    "Next Page",
			CustomID: "next",
		},
	},
}

func (b *EoD) newPageSwitcher(ps pageSwitcher, m msg, rsp rsp) {
	rsp.Acknowledge()

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
	_, exists = dat.playChannels[m.ChannelID]
	if exists {
		ps.PageLength = playPageLength
	}

	cont, _, length, err := ps.PageGetter(ps)
	if rsp.Error(err) {
		return
	}
	id := rsp.Embed(&discordgo.MessageEmbed{
		Title:       ps.Title,
		Description: cont,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: ps.Thumbnail,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Page %d/%d", ps.Page+1, length+1),
		},
	}, btnRow)

	if dat.pageSwitchers == nil {
		dat.pageSwitchers = make(map[string]pageSwitcher)
	}

	fmt.Println(id)
	dat.lock.Lock()
	dat.pageSwitchers[id] = ps
	dat.lock.Unlock()

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

	dat.lock.RLock()
	ps, exists := dat.pageSwitchers[i.Message.ID]
	dat.lock.RUnlock()
	if !exists {
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

	color, _ := b.getColor(i.GuildID, i.Member.User.ID)
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
						Text: fmt.Sprintf("Page %d/%d", ps.Page+1, length+1),
					},
					Color: color,
				},
			},
			Components: []discordgo.MessageComponent{btnRow},
		},
	})
	if err != nil {
		log.Println("failed to update page switcher:", err)
	}
	dat.pageSwitchers[i.Message.ID] = ps

	lock.Lock()
	b.dat[i.GuildID] = dat
	lock.Unlock()
}
