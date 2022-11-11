package base

import (
	"fmt"

	"github.com/Nv7-Github/sevcord/v2"
)

func (b *Base) CheckCtx(ctx sevcord.Ctx) bool {
	var cnt int
	err := b.db.Select(&[]*int{&cnt}, "SELECT COUNT(*) FROM config WHERE guild=$1", ctx.Guild())
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

	return true
}
