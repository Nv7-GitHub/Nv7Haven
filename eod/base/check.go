package base

import (
	"fmt"

	"github.com/Nv7-Github/sevcord/v2"
)

func (b *Base) CheckCtx(ctx sevcord.Ctx) bool {
	var cnt int
	err := b.db.Select(&cnt, "SELECT COUNT(*) FROM guilds WHERE id=?", ctx.Guild())
	if err != nil {
		b.Error(ctx, err)
		return false
	}

	// Not configured
	if cnt == 0 {
		ctx.Respond(sevcord.NewMessage(fmt.Sprintf("⚠️ This server is not configured! Configure it with </config:%s>.", configCmdId)))
		return false
	}

	return true
}
