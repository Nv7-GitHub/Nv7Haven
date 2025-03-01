package base

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/sevcord/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/lib/pq"
)

func (b *Base) Give(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	user := opts[0].(*discordgo.User).ID
	q, ok := b.CalcQuery(c, opts[1].(string))
	if !ok {
		return
	}

	// Combine to inv
	_, err := b.db.Exec(`UPDATE inventories SET inv=inv | $1 WHERE guild=$2 AND "user"=$3`, pq.Array(q.Elements), c.Guild(), user)
	if err != nil {
		b.Error(c, err)
		return
	}

	// Respond
	c.Respond(sevcord.NewMessage(fmt.Sprintf("Succesfully gave elements to <@%s>!", user)))
}
func (b *Base) GiveAchievement(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	user := opts[0].(*discordgo.User).ID

	var newachievement int
	b.db.QueryRow(`SELECT id FROM achievements WHERE guild=$1 AND LOWER(name)=$2`, c.Guild(), strings.ToLower(opts[1].(string))).Scan(&newachievement)
	var newachievementarray []int
	newachievementarray = append(newachievementarray, newachievement)
	_, err := b.db.Exec(`UPDATE achievers SET achievements=achievements | $1 WHERE guild=$2 AND "user"=$3`, pq.Array(newachievementarray), c.Guild(), user)
	if err != nil {
		b.Error(c, err)
		return
	}
	c.Respond(sevcord.NewMessage(fmt.Sprintf("Succesfully gave achievement to <@%s>!", user)))
}
