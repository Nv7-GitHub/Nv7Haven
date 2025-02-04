package achievements

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/sevcord/v2"
	"github.com/bwmarrin/discordgo"
)

func (u *Users) Profile(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	var userID string
	if opts[0] != nil && opts[0] != "" {
		userID = opts[0].(*discordgo.User).ID
	} else {
		userID = c.Author().User.ID
	}
	u.ProfileHandler(c, userID)
}
func (u *Users) ProfileHandler(c sevcord.Ctx, params string) {

	parts := strings.Split(params, "|")
	user, err := c.Dg().GuildMember(c.Guild(), parts[0])
	if err != nil {
		return
	}
	var invsize int
	err = u.db.QueryRow(`SELECT array_length(inv, 1) FROM inventories WHERE guild=$1 AND "user"=$2`, c.Guild(), parts[0]).Scan(&invsize)
	if err != nil {
		u.base.Error(c, err)
		return

	}
	var madesize int
	err = u.db.QueryRow("SELECT COUNT(*) FROM elements WHERE guild=$1 AND creator=$2", c.Guild(), parts[0]).Scan(&madesize)
	if err != nil {
		u.base.Error(c, err)
		return
	}
	var votecnt int

	err = u.db.QueryRow(`SELECT votecnt FROM inventories WHERE guild=$1 AND "user"=$2`, c.Guild(), parts[0]).Scan(&votecnt)
	if err != nil {
		u.base.Error(c, err)
		return
	}
	var achievementcnt int

	err = u.db.QueryRow(`SELECT cardinality(achievements) FROM achievers WHERE guild=$1 AND "user"=$2`, c.Guild(), parts[0]).Scan(&achievementcnt)
	if err != nil {
		u.base.Error(c, err)
		return
	}

	//get user avatar
	img := user.User.AvatarURL("128")
	emb := sevcord.NewEmbed().
		Title(user.User.Username + "'s Profile")
	emb = emb.AddField("User", user.Mention(), false)
	emb = emb.AddField("🔍 Elements Found", fmt.Sprintf("%d", invsize), true)
	emb = emb.AddField("💡 Elements Made", fmt.Sprintf("%d", madesize), true)
	emb = emb.AddField("🗳️ Votes Cast", fmt.Sprintf("%d", votecnt), true)
	emb = emb.AddField("🏆 Achievements Found", fmt.Sprintf("%d", achievementcnt), false)

	emb = emb.Thumbnail(img)
	c.Respond(sevcord.NewMessage("").AddEmbed(emb))
}
