package discord

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type property struct {
	Name        string
	Value       int // Money: Value*Upgrades* Hours since last visited
	Cost        int // Initial Price
	UpgradeCost int // Upgrade^1.5 * UpgradeCost
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

	if strings.HasPrefix(m.Content, "props") {
		b.checkuser(m)
		var text string
		for _, prop := range upgrades {
			text += fmt.Sprintf("**%s** - id `%s` - %d coins\n\n", prop.Name, prop.ID, prop.Cost)
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title:       "Available Properties",
			Description: text,
		})
		return
	}

	if strings.HasPrefix(m.Content, "inv") {
		id := m.Author.ID
		if len(m.Mentions) > 0 {
			id = m.Mentions[0].ID
			exists, suc := b.exists(m, "currency", "user=?", id)
			if !suc {
				return
			}
			if !exists {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User <@%s> has never used this bot's currency commands.", id))
				return
			}
		}
		user, suc := b.getuser(m, id)
		if !suc {
			return
		}

		usr, err := s.User(id)
		if b.handle(err, m) {
			return
		}

		var text string
		for _, prop := range user.Properties {
			text += fmt.Sprintf("`%s` - %d upgrades\n\n", prop.ID, prop.Upgrades)
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("%s's Properties", usr.Username),
			Description: text,
		})
		return
	}

	if strings.HasPrefix(m.Content, "purchase") {
		b.checkuser(m)
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
		user.LastVisited = time.Now().Unix()
		b.updateuser(m, user)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You bought %s for %d coins!", plc, prp.Cost))
		return
	}

	if strings.HasPrefix(m.Content, "upgrade") {
		b.checkuser(m)
		var plc string
		_, err := fmt.Sscanf(m.Content, "upgrade %s", &plc)
		if b.handle(err, m) {
			return
		}

		var info property
		suc := false
		for _, property := range upgrades {
			if property.ID == plc {
				suc = true
				info = property
				break
			}
		}
		if !suc {
			s.ChannelMessageSend(m.ChannelID, "That property doesn't exist!")
			return
		}

		user, suc := b.getuser(m, m.Author.ID)
		if !suc {
			return
		}

		var place prop
		index := -1
		for i, plac := range user.Properties {
			if plac.ID == plc {
				index = i
				place = plac
				break
			}
		}
		if index < 0 {
			s.ChannelMessageSend(m.ChannelID, "You don't have that property!")
			return
		}

		upgradeCost := int(math.Pow(float64(place.Upgrades), 1.5) * float64(info.UpgradeCost))
		if user.Wallet < upgradeCost {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You need %d more coins to upgrade property %s!!", upgradeCost-user.Wallet, info.Name))
			return
		}

		user.Properties[index].Upgrades++
		user.Wallet -= upgradeCost

		suc = b.updateuser(m, user)
		if !suc {
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Upgraded %s to Level %d!", info.Name, user.Properties[index].Upgrades))
		return
	}

	if strings.HasPrefix(m.Content, "collect") {
		b.checkuser(m)
		user, suc := b.getuser(m, m.Author.ID)
		if !suc {
			return
		}
		moneyCollected := 0
		for _, prop := range user.Properties {
			var val int
			for _, upgrd := range upgrades {
				if upgrd.ID == prop.ID {
					val = upgrd.Value
				}
			}

			coll := float32(val*prop.Upgrades) * (float32(time.Now().Unix()-user.LastVisited) / 3600)
			moneyCollected += int(coll)
		}
		user.Wallet += moneyCollected
		user.LastVisited = time.Now().Unix()
		suc = b.updateuser(m, user)
		if !suc {
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You collected %d coins!", moneyCollected))
		return
	}
}
