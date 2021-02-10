package discord

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) specials(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if b.startsWith(m, "rob") {
		b.checkuser(m)

		if !(len(m.Mentions) > 0) {
			s.ChannelMessageSend(m.ChannelID, "You need to mention the person you are going to rob!")
			return
		}

		b.checkuserwithid(m, m.Mentions[0].ID)

		user1, suc := b.getuser(m, m.Author.ID)
		if !suc {
			return
		}

		ups, exists := user1.Properties["rob"]
		if !exists {
			s.ChannelMessageSend(m.ChannelID, "You need property `rob` to Rob people!")
			return
		}

		user2, suc := b.getuser(m, m.Mentions[0].ID)
		if !suc {
			return
		}

		var num int
		_, err := fmt.Sscanf(m.Content, "rob %d", &num)
		if b.handle(err, m) {
			return
		}
		num = b.abs(num)

		if user2.Wallet < num {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User <@%s> doesn't even have %d coins!", m.Mentions[0].ID, num))
			return
		}

		if user1.Wallet < num {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("If you are going to rob someone of %d coins, you need to have that many coins in case you get caught.", num))
			return
		}

		if (num / 10) > 0 {
			num -= rand.Intn(num / 10) // loss
		}

		// backfiring
		backNum := ups - 2
		randChance := int(math.Pow(float64(ups), 1.5))
		if randChance != 0 {
			backNum += rand.Intn(randChance)
		}

		if backNum < 0 {
			user1.Wallet -= num
			user2.Wallet += num
			b.updateuser(m, user1)
			b.updateuser(m, user2)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Oh no! You got caught and had to give %d coins to the person you were stealing from! Try upgrading property `rob` to reduce the chances of backfiring!", num))
			return
		}
		user1.Wallet += num
		user2.Wallet -= num
		b.updateuser(m, user1)
		b.updateuser(m, user2)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Everything went perfectly and you just stole %d coins!", num))
		return
	}
}
