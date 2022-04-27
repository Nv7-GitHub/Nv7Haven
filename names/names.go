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
	return n, nil
}
