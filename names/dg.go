package names

import (
	_ "embed"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//go:embed token.txt
var token string

const appID = "968719598121783327"
const guild = ""

func (n *Names) initDG() error {
	dg, err := discordgo.New("Bot " + strings.TrimSpace(token))
	if err != nil {
		return err
	}
	n.dg = dg
	_, err = n.dg.ApplicationCommandBulkOverwrite(appID, guild, commands)
	if err != nil {
		return err
	}
	return n.dg.Open()
}
