package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type property struct {
	Name        string
	Value       int // Money: Value*Upgrades* Hours since last visited
	Upgrades    int
	Cost        int // Initial Price
	UpgradeCost int // Upgrade^1.5 * UpgradeCost + Cost
	ID          string
	Credit      int
}

type prop struct {
	ID       string
	Upgrades int
}

var upgrades = []property{
	property{
		Name:        "Snack Booth",
		Value:       20,
		Cost:        1000,
		UpgradeCost: 200,
		ID:          "snack",
		Credit:      50,
	},
	property{
		Name:        "Homemade Cookie Business",
		Value:       50,
		Cost:        10000,
		UpgradeCost: 600,
		ID:          "cookie",
		Credit:      100,
	},
	property{
		Name:        "Li'l Jon'z Fudge Store",
		Value:       80,
		UpgradeCost: 960,
		Cost:        50000,
		ID:          "fudge",
		Credit:      150,
	},
	property{
		Name:        `|\\/|cDonaIds`,
		Cost:        100000,
		Value:       100,
		UpgradeCost: 1400,
		ID:          "mcd",
		Credit:      200,
	},
	property{
		Name:        "Village Bank",
		Value:       120,
		Cost:        200000,
		UpgradeCost: 1500,
		ID:          "village",
		Credit:      250,
	},
	property{
		Name:        "Vanilla JS Coders",
		Value:       140,
		Cost:        400000,
		UpgradeCost: 1750,
		ID:          "jspain",
		Credit:      350,
	},
	property{
		Name:        "We Use Hacks In Creative",
		Value:       200,
		Cost:        400000,
		UpgradeCost: 2500,
		ID:          "rich",
		Credit:      400,
	},
}

func (b *Bot) properties(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "purchase") {
		var plc string
		_, err := fmt.Sscanf(m.Content, "purchase %s", &plc)
		if b.handle(err, m) {
			return
		}

		isInProperties := false
		var prp property
		for _, property := range upgrades {
			if property.ID == plc {
				isInProperties = true
				prp = property
			}
		}

		if !isInProperties {
			s.ChannelMessageSend(m.ChannelID, "You need to specify a property that exists. Rememeber to use the id and not the name!")
			return
		}

		user, suc := b.getuser(m, m.Author.ID)
		if !suc {
			return
		}
		for _, property := range user.Properties {
			if property.ID == plc {
				s.ChannelMessageSend(m.ChannelID, "You already have this property! Use the upgrade command to upgrade it.")
				return
			}
		}

		if user.Wallet < prp.Cost {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You need %d more coins to buy this property.", prp.Cost-user.Wallet))
			return
		}
		if user.Credit < prp.Credit {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You need %d more credit score to buy this property.", prp.Credit-user.Credit))
			return
		}
		user.Wallet -= prp.Cost
		place := prop{
			ID:       prp.ID,
			Upgrades: 0,
		}
		user.Properties = append(user.Properties, place)
		b.updateuser(m, user)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You bought %s for %d coins!", plc, prp.Cost))
		return
	}
}
