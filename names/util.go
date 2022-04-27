package names

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (n *Names) Err(msg string, i *discordgo.Interaction) {
	n.dg.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   1 << 6,
			Content: fmt.Sprintf("**Error**: %s 🔴", msg),
		},
	})
}
