package base

import (
	"fmt"

	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
)

func (b *Base) CheckCtx(ctx sevcord.Ctx) bool {
	var cnt int
	err := b.db.QueryRow("SELECT COUNT(*) FROM config WHERE guild=$1", ctx.Guild()).Scan(&cnt)
	if err != nil {
		b.Error(ctx, err)
		return false
	}

	// Not configured
	if cnt == 0 {
		ctx.Acknowledge()
		ctx.Respond(sevcord.NewMessage(fmt.Sprintf("⚠️ This server is not configured! Configure it with </config:%s>.", configCmdId)))
		return false
	}

	// Check if user has account
	var exists bool
	err = b.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM inventories WHERE "user"='$2 AND guild=$1)`, ctx.Guild(), ctx.Author().User.ID).Scan(&exists)
	if err != nil {
		b.Error(ctx, err)
		return false
	}
	if !exists {
		// Make account
		_, err = b.db.Exec(`INSERT INTO inventories (guild, "user", inv) VALUES ($1, $2, $3)`, ctx.Guild(), ctx.Author().User.ID, pq.Array([]int{1, 2, 3, 4}))
		if err != nil {
			b.Error(ctx, err)
			return false
		}
	}

	return true
}
