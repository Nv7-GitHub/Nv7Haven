package categories

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
)

func (c *Categories) ImageCmd(ctx sevcord.Ctx, cat string, image string) {
	ctx.Acknowledge()

	// Check element
	var name string
	var old string
	err := c.db.QueryRow("SELECT name, image FROM categories WHERE LOWER(name)=$1 AND guild=$2", strings.ToLower(cat), ctx.Guild()).Scan(&name, &old)
	if err != nil {
		c.base.Error(ctx, err, "Category **"+cat+"** doesn't exist!")
		return
	}

	// Make poll
	res := c.polls.CreatePoll(ctx, &types.Poll{
		Kind: types.PollKindCatImage,
		Data: types.PgData{
			"cat": name,
			"new": image,
			"old": old,
		},
	})
	if !res.Ok {
		ctx.Respond(res.Response())
		return
	}

	//Distinguish between adding new image and changing existing image
	var addtext string
	if old != "" {
		addtext = "to edit the"
	} else {
		addtext = "a new"
	}
	// Respond
	ctx.Respond(sevcord.NewMessage(fmt.Sprintf("Suggested %s image for category **%s** 📷", addtext, name)))
}

func (c *Categories) MsgSignCmd(ctx sevcord.Ctx, cat string, mark string) {
	ctx.Acknowledge()

	var name string
	var old string
	err := c.db.QueryRow("SELECT name, comment FROM categories WHERE LOWER(name)=$1 AND guild=$2", strings.ToLower(cat), ctx.Guild()).Scan(&name, &old)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.Respond(sevcord.NewMessage("Element **" + cat + "** doesn't exist! " + types.RedCircle))
			return
		} else {
			c.base.Error(ctx, err)
			return
		}
	}

	// Make poll
	res := c.polls.CreatePoll(ctx, &types.Poll{
		Kind: types.PollKindCatComment,
		Data: types.PgData{
			"cat": name,
			"new": mark,
			"old": old,
		},
	})
	if !res.Ok {
		ctx.Respond(res.Response())
		return
	}
	//Distinguish between new and old
	var addtext string
	if old != types.DefaultMark {
		addtext = "to edit the"
	} else {
		addtext = "a new"
	}
	// Respond
	ctx.Respond(sevcord.NewMessage(fmt.Sprintf("Suggested %s mark for **%s** 🖋️", addtext, name)))
}

func (c *Categories) SignCmd(ctx sevcord.Ctx, opts []any) {
	// Check element
	var name string
	var old string
	err := c.db.QueryRow("SELECT name, comment FROM categories WHERE LOWER(name)=$1 AND guild=$2", strings.ToLower(opts[0].(string)), ctx.Guild()).Scan(&name, &old)
	if err != nil {
		c.base.Error(ctx, err, "Category **"+opts[0].(string)+"** doesn't exist!")
		return
	}

	// Get mark
	ctx.(*sevcord.InteractionCtx).Modal(sevcord.NewModal("Sign Category", func(ctx sevcord.Ctx, s []string) {
		// Make poll
		res := c.polls.CreatePoll(ctx, &types.Poll{
			Kind: types.PollKindCatComment,
			Data: types.PgData{
				"cat": name,
				"new": s[0],
				"old": old,
			},
		})
		if !res.Ok {
			ctx.Respond(res.Response())
			return
		}
		var addtext string
		if old != types.DefaultMark {
			addtext = "to edit the"
		} else {
			addtext = "a new"
		}

		// Respond
		ctx.Respond(sevcord.NewMessage(fmt.Sprintf("Suggested %s mark for category **%s** 🖋️", addtext, name)))
	}).Input(sevcord.NewModalInput("New Comment", types.DefaultMark, sevcord.ModalInputStyleParagraph, 2400)))
}

func (c *Categories) ColorCmd(ctx sevcord.Ctx, opts []any) {
	ctx.Acknowledge()

	// Check hex code
	code := opts[1].(string)
	val, err := strconv.ParseInt(strings.TrimPrefix(code, "#"), 16, 64)
	if err != nil {
		ctx.Respond(sevcord.NewMessage("Invalid hex code! " + types.RedCircle))
		return
	}
	if val < 0 || val > 16777215 {
		ctx.Respond(sevcord.NewMessage("Invalid hex code! " + types.RedCircle))
		return
	}

	// Check category
	var name string
	var old int
	var colorer string
	err = c.db.QueryRow("SELECT name, color,colorer FROM categories WHERE LOWER(name)=$1 AND guild=$2", strings.ToLower(opts[0].(string)), ctx.Guild()).Scan(&name, &old, &colorer)
	if err != nil {
		c.base.Error(ctx, err, "Category **"+opts[0].(string)+"** doesn't exist!")
		return
	}

	// Make poll
	res := c.polls.CreatePoll(ctx, &types.Poll{
		Kind: types.PollKindCatColor,
		Data: types.PgData{
			"cat": name,
			"new": float64(val),
			"old": float64(old),
		},
	})
	if !res.Ok {
		ctx.Respond(res.Response())
		return
	}
	var addtext string
	if colorer != "" {
		addtext = "to edit the"
	} else {
		addtext = "a new"
	}
	// Respond
	ctx.Respond(sevcord.NewMessage(fmt.Sprintf("Suggested %s color for category **%s** 🎨", addtext, name)))
}
