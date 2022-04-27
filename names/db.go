package names

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (n *Names) setNameCmd(i *discordgo.InteractionCreate) {
	// Get guild owner
	info, err := n.dg.Guild(i.GuildID)
	if err != nil {
		n.Err(err.Error(), i.Interaction)
		return
	}
	if info.OwnerID != i.Member.User.ID {
		n.Err("You must be owner to run this command!", i.Interaction)
		return
	}

	var user string
	var name string
	d := i.ApplicationCommandData()
	for _, par := range d.Options {
		switch par.Name {
		case "user":
			user = par.UserValue(n.dg).ID

		case "name":
			name = par.StringValue()
		}
	}

	// Update in DB
	_, err = n.db.Exec("REPLACE INTO names_discord (guild, user, name) VALUES(?, ?, ?)", i.GuildID, user, name)
	if err != nil {
		n.Err(err.Error(), i.Interaction)
		return
	}

	n.dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   1 << 6,
			Content: fmt.Sprintf("Updated <@%s>!", user),
		},
	})
}

func (n *Names) getNameCmd(user string, i *discordgo.InteractionCreate) {
	// Get count
	var cnt int
	err := n.db.QueryRow("SELECT COUNT(1) FROM names_discord WHERE guild=? AND user=?", i.GuildID, user).Scan(&cnt)
	if err != nil {
		n.Err(err.Error(), i.Interaction)
		return
	}
	if cnt == 0 {
		// Name doesn't exist
		n.Err("Name not found!", i.Interaction)
		return
	}

	// Get name
	var name string
	err = n.db.QueryRow("SELECT name FROM names_discord WHERE guild=? AND user=?", i.GuildID, user).Scan(&name)
	if err != nil {
		n.Err(err.Error(), i.Interaction)
		return
	}

	// Return
	n.dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   1 << 6,
			Content: name,
		},
	})
}
