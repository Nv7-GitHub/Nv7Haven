package names

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const max = 25

func (n *Names) search(i *discordgo.InteractionCreate) {
	d := i.ApplicationCommandData()
	val := d.Options[0].StringValue()

	// Query
	res, err := n.db.Query("SELECT name from names_discord WHERE guild=? AND name LIKE ? ORDER BY name", i.GuildID, "%"+val+"%")
	if err != nil {
		n.Err(err.Error(), i.Interaction)
		return
	}

	// Make pars
	pars := make([]*discordgo.ApplicationCommandOptionChoice, 0)
	for res.Next() {
		var nm string
		err = res.Scan(&nm)
		if err != nil {
			n.Err(err.Error(), i.Interaction)
			return
		}
		pars = append(pars, &discordgo.ApplicationCommandOptionChoice{
			Name:  nm,
			Value: nm,
		})
	}

	// Limit
	if len(pars) > max {
		pars = pars[:max]
	}

	n.dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: pars,
		},
	})
}

func (n *Names) searchCmd(i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Options[0].StringValue()

	// Check if exists
	var cnt int
	err := n.db.QueryRow("SELECT COUNT(1) FROM names_discord WHERE guild=? AND name=?", i.GuildID, name).Scan(&cnt)
	if err != nil {
		n.Err(err.Error(), i.Interaction)
		return
	}
	if cnt == 0 {
		// Name doesn't exist
		n.Err("Name not found!", i.Interaction)
		return
	}

	// Get user
	var usr string
	err = n.db.QueryRow("SELECT user FROM names_discord WHERE guild=? AND name=?", i.GuildID, name).Scan(&usr)
	if err != nil {
		n.Err(err.Error(), i.Interaction)
		return
	}

	// Return
	n.dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   1 << 6,
			Content: fmt.Sprintf("<@%s>", usr),
		},
	})
}
