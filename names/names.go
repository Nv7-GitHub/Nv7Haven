package names

import (
	"github.com/Nv7-Github/Nv7Haven/db"
	"github.com/bwmarrin/discordgo"
)

type Names struct {
	db *db.DB
	dg *discordgo.Session
}

func (n *Names) Close() {
	n.dg.Close()
}

func NewNames(d *db.DB) (*Names, error) {
	n := &Names{
		db: d,
	}
	err := n.initDG()
	if err != nil {
		return nil, err
	}
	n.dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			cmd := i.ApplicationCommandData()
			switch cmd.Name {
			case "set":
				n.setNameCmd(i)

			case "get":
				n.getNameCmd(cmd.Options[0].UserValue(n.dg).ID, i)

			case "search":
				n.searchCmd(i)

			case "unnamed":
				n.unnamedCmd(i)

			case "View Name":
				n.getNameCmd(cmd.TargetID, i)
			}

		case discordgo.InteractionApplicationCommandAutocomplete:
			n.search(i)
		}
	})
	return n, nil
}
