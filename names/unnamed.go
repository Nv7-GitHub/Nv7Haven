package names

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (n *Names) unnamedCmd(i *discordgo.InteractionCreate) {
	// Check guild owner
	info, err := n.dg.Guild(i.GuildID)
	if err != nil {
		n.Err(err.Error(), i.Interaction)
		return
	}
	if info.OwnerID != i.Member.User.ID {
		n.Err("You must be owner to run this command!", i.Interaction)
		return
	}

	// Get ids of already in DB
	res, err := n.db.Query("SELECT user FROM names_discord WHERE guild=?", i.GuildID)
	if err != nil {
		n.Err(err.Error(), i.Interaction)
		return
	}
	ids := make(map[string]struct{})
	for res.Next() {
		var id string
		err = res.Scan(&id)
		if err != nil {
			n.Err(err.Error(), i.Interaction)
			return
		}
		ids[id] = struct{}{}
	}

	// Get members
	mem, err := n.dg.GuildMembers(i.GuildID, "", 1000)
	if err != nil {
		n.Err(err.Error(), i.Interaction)
		return
	}

	// Find unnamed
	unnamed := make([]string, 0)
	for _, v := range mem {
		_, exists := ids[v.User.ID]
		if !exists {
			unnamed = append(unnamed, v.User.ID)
		}
	}

	// Make list
	content := &strings.Builder{}
	for _, id := range unnamed {
		fmt.Fprintf(content, "<@%s>", id)
	}

	n.dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: 1 << 6,
			Embeds: []*discordgo.MessageEmbed{{
				Title:       "Unnamed",
				Description: content.String(),
			}},
		},
	})
}
