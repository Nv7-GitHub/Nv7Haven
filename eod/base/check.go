package base

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/lib/pq"
)

func (b *Base) CheckCtx(ctx sevcord.Ctx, cmd string) bool {
	var cnt int
	err := b.db.QueryRow("SELECT COUNT(*) FROM config WHERE guild=$1 AND config IS NOT NULL", ctx.Guild()).Scan(&cnt)
	if err != nil {
		b.Error(ctx, err)
		return false
	}

	// Not configured
	if cnt == 0 {
		// Check if empty row
		err := b.db.QueryRow("SELECT COUNT(*) FROM config WHERE guild=$1", ctx.Guild()).Scan(&cnt)
		if err != nil {
			b.Error(ctx, err)
			return false
		}

		if cnt == 0 {
			// Make empty row
			_, err = b.db.Exec(`INSERT INTO config (guild, play, pollcnt) VALUES ($1, $2, $3)`, ctx.Guild(), pq.Array([]string{}), 0)
			if err != nil {
				b.Error(ctx, err)
				return false
			}

			// Check if starters are there
			var exists bool
			err = b.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM elements WHERE guild=$1 LIMIT 1)`, ctx.Guild()).Scan(&exists)
			if err != nil {
				b.Error(ctx, err)
				return false
			}
			if !exists { // No starters, insert
				_, err = b.db.NamedExec("INSERT INTO elements (id, guild, name, image, color, comment, creator, createdon, parents, treesize, commenter, colorer, imager) VALUES (:id, :guild, :name, :image, :color, :comment, :creator, :createdon, :parents, :treesize, :commenter, :colorer, :imager)", types.Starters(ctx.Guild()))
				if err != nil {
					b.Error(ctx, err)
					return false
				}
			}
		}

		if cmd == "config" {
			return true // Pass to config command
		}

		ctx.Acknowledge()
		ctx.Respond(sevcord.NewMessage(fmt.Sprintf("⚠️ This server is not configured! Configure it with </config:%s>.", configCmdId)))
		return false
	}

	// Check if user has account
	var exists bool
	err = b.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM inventories WHERE "user"=$2 AND guild=$1)`, ctx.Guild(), ctx.Author().User.ID).Scan(&exists)
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

var notAllowed = []string{
	"\n",
	"<@",
	"\t",
	"`",
	"@everyone",
	"@here",
	"<t:",
	"</",
}

var charReplace = map[rune]rune{
	'’': '\'',
	'‘': '\'',
	'`': '\'',
	'”': '"',
	'“': '"',
}

// CheckName checks the validity of a name & returns a cleaned up version + error
func CheckName(name string) (string, types.Resp) {
	for _, v := range notAllowed {
		if strings.Contains(name, v) {
			return "", types.Fail("A name may not contain '" + v + "'!")
		}
	}
	for k, v := range charReplace {
		if strings.ContainsRune(name, k) {
			name = strings.ReplaceAll(name, string(k), string(v))
		}
	}
	return name, types.Ok()
}
