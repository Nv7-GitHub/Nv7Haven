package discord

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) properties(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if b.startsWith(m, "props") {
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

	if b.startsWith(m, "prop") {
		b.checkuser(m)
		var id string
		fmt.Sscanf(m.Content, "prop %s", &id)
		prop, exists := b.props[id]
		if !exists {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("There aren't any properties with ID `%s`!", id))
			return
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: fmt.Sprintf("Property Info: %s", prop.Name),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Price",
					Value: strconv.Itoa(prop.Cost),
				},
				{
					Name:  "Credit Required",
					Value: strconv.Itoa(prop.Credit),
				},
				{
					Name:  "Money/Upgrade",
					Value: strconv.Itoa(prop.Value),
				},
			},
		})
		return
	}

	if b.startsWith(m, "inv") {
		b.checkuser(m)
		id := m.Author.ID
		if len(m.Mentions) > 0 {
			id = m.Mentions[0].ID
			b.checkuserwithid(m, id)
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
		for id, ups := range user.Properties {
			text += fmt.Sprintf("`%s` - %d upgrades\n\n", id, ups)
		}
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("%s's Properties", usr.Username),
			Description: text,
		})
		return
	}

	if b.startsWith(m, "purchase") {
		b.checkuser(m)
		var plc string
		_, err := fmt.Sscanf(m.Content, "purchase %s", &plc)
		if b.handle(err, m) {
			return
		}

		prp, exists := b.props[plc]

		if !exists {
			s.ChannelMessageSend(m.ChannelID, "You need to specify a property that exists. Rememeber to use the id and not the name!")
			return
		}

		user, suc := b.getuser(m, m.Author.ID)
		if !suc {
			return
		}
		_, exists = user.Properties[plc]
		if exists {
			s.ChannelMessageSend(m.ChannelID, "You already have this property! Use the upgrade command to upgrade it.")
			return
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
		user.Properties[prp.ID] = 0
		user.LastVisited = time.Now().Unix()
		b.updateuser(m, user)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You bought %s for %d coins!", plc, prp.Cost))
		return
	}

	if b.startsWith(m, "upgrade") {
		b.checkuser(m)
		var plc string
		var numVal string
		_, err := fmt.Sscanf(m.Content, "upgrade %s %s", &plc, &numVal)
		if b.handle(err, m) {
			return
		}

		info, exists := b.props[plc]
		if !exists {
			s.ChannelMessageSend(m.ChannelID, "That property doesn't exist!")
			return
		}

		user, suc := b.getuser(m, m.Author.ID)
		if !suc {
			return
		}

		ups, exists := user.Properties[plc]
		if !exists {
			s.ChannelMessageSend(m.ChannelID, "You don't have that property!")
			return
		}

		upgradeCost := 0
		numUpgrades := 0
		if numVal == "max" {
			for upgradeCost+int(math.Pow(float64(ups+numUpgrades), 1.5)*float64(info.UpgradeCost)) < user.Wallet {
				upgradeCost += int(math.Pow(float64(ups+numUpgrades), 1.5) * float64(info.UpgradeCost))
				numUpgrades++
			}
		} else {
			num, err := strconv.Atoi(numVal)
			if b.handle(err, m) {
				return
			}
			num = b.abs(num)
			numUpgrades = num
			for i := 0; i < num; i++ {
				upgradeCost += int(math.Pow(float64(ups+i), 1.5) * float64(info.UpgradeCost))
			}
		}

		if user.Wallet < upgradeCost {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You need %d more coins to upgrade property %s!!", upgradeCost-user.Wallet, info.Name))
			return
		}

		user.Properties[plc] += numUpgrades
		user.Wallet -= upgradeCost

		suc = b.updateuser(m, user)
		if !suc {
			return
		}

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Upgraded %s to Level %d!", info.Name, user.Properties[plc]))
		return
	}

	if b.startsWith(m, "collect") {
		b.checkuser(m)
		user, suc := b.getuser(m, m.Author.ID)
		if !suc {
			return
		}
		moneyCollected := 0
		for id, ups := range user.Properties {
			val := b.props[id].Value

			coll := float32(val*ups) * (float32(time.Now().Unix()-user.LastVisited) / 3600)
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
