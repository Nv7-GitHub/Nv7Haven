package queries

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
)

func (q *Queries) ImageCmd(ctx sevcord.Ctx, opts []any) {
	ctx.Acknowledge()

	// Check element
	var name string
	var old string
	err := q.db.QueryRow("SELECT name, image FROM queries WHERE LOWER(name)=$1 AND guild=$2", strings.ToLower(opts[0].(string)), ctx.Guild()).Scan(&name, &old)
	if err != nil {
		q.base.Error(ctx, err)
		return
	}

	// Check image
	if !strings.HasPrefix(opts[1].(*sevcord.SlashCommandAttachment).ContentType, "image") {
		ctx.Respond(sevcord.NewMessage("The attachment must be an image! " + types.RedCircle))
		return
	}

	// Make poll
	q.polls.CreatePoll(ctx, &types.Poll{
		Kind: types.PollKindQueryImage,
		Data: types.PgData{
			"query": name,
			"new":   opts[1].(*sevcord.SlashCommandAttachment).URL,
			"old":   old,
		},
	})

	// Respond
	ctx.Respond(sevcord.NewMessage(fmt.Sprintf("Suggested an image for query **%s** 📷", name)))
}

func (q *Queries) SignCmd(ctx sevcord.Ctx, opts []any) {
	// Check element
	var name string
	var old string
	err := q.db.QueryRow("SELECT name, comment FROM queries WHERE LOWER(name)=$1 AND guild=$2", strings.ToLower(opts[0].(string)), ctx.Guild()).Scan(&name, &old)
	if err != nil {
		q.base.Error(ctx, err)
		return
	}

	// Get mark
	ctx.(*sevcord.InteractionCtx).Modal(sevcord.NewModal("Sign Query", func(ctx sevcord.Ctx, s []string) {
		// Make poll
		q.polls.CreatePoll(ctx, &types.Poll{
			Kind: types.PollKindQueryComment,
			Data: types.PgData{
				"query": name,
				"new":   s[0],
				"old":   old,
			},
		})

		// Respond
		ctx.Respond(sevcord.NewMessage(fmt.Sprintf("Suggested a note for query **%s** 🖋️", name)))
	}).Input(sevcord.NewModalInput("New Comment", "None", sevcord.ModalInputStyleParagraph, 2400)))
}

func (q *Queries) ColorCmd(ctx sevcord.Ctx, opts []any) {
	ctx.Acknowledge()

	// Check hex code
	code := opts[1].(string)
	if !strings.HasPrefix(code, "#") {
		ctx.Respond(sevcord.NewMessage("Invalid hex code! " + types.RedCircle))
		return
	}
	val, err := strconv.ParseInt(strings.TrimPrefix(code, "#"), 16, 64)
	if err != nil {
		q.base.Error(ctx, err)
		return
	}
	if val < 0 || val > 16777215 {
		ctx.Respond(sevcord.NewMessage("Invalid hex code! " + types.RedCircle))
		return
	}

	// Check element
	var name string
	var old int
	err = q.db.QueryRow("SELECT name, color FROM queries WHERE LOWER(name)=$1 AND guild=$2", strings.ToLower(opts[0].(string)), ctx.Guild()).Scan(&name, &old)
	if err != nil {
		q.base.Error(ctx, err)
		return
	}

	// Make poll
	q.polls.CreatePoll(ctx, &types.Poll{
		Kind: types.PollKindQueryColor,
		Data: types.PgData{
			"query": name,
			"new":   float64(val),
			"old":   float64(old),
		},
	})

	// Respond
	ctx.Respond(sevcord.NewMessage(fmt.Sprintf("Suggested a color for query **%s** 🎨", name)))
}
