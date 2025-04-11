package elements

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
)

func (e *Elements) ImageCmd(c sevcord.Ctx, id int, image string) {
	c.Acknowledge()

	// Check element
	var elem string
	var old string
	err := e.db.QueryRow("SELECT name, image FROM elements WHERE id=$1 AND guild=$2", id, c.Guild()).Scan(&elem, &old)
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Make poll
	res := e.polls.CreatePoll(c, &types.Poll{
		Kind: types.PollKindImage,
		Data: types.PgData{
			"elem": float64(id),
			"new":  image,
			"old":  old,
		},
	})
	if !res.Ok {
		c.Respond(res.Response())
		return
	}
	// Distinguish between adding new image and changing existing image
	var addtext string
	if old != "" {
		addtext = "to edit the"
	} else {
		addtext = "a new"
	}
	// Respond
	c.Respond(sevcord.NewMessage(fmt.Sprintf("Suggested %s image for **%s** 📷", addtext, elem)))
}

func (e *Elements) SignCmd(c sevcord.Ctx, opts []any) {
	// Check element
	var name string
	var old string
	err := e.db.QueryRow("SELECT name, comment FROM elements WHERE id=$1 AND guild=$2", opts[0].(int64), c.Guild()).Scan(&name, &old)
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Get mark
	c.(*sevcord.InteractionCtx).Modal(sevcord.NewModal("Sign Element", func(c sevcord.Ctx, s []string) {
		// Make poll
		res := e.polls.CreatePoll(c, &types.Poll{
			Kind: types.PollKindComment,
			Data: types.PgData{
				"elem": float64(opts[0].(int64)),
				"new":  s[0],
				"old":  old,
			},
		})
		if !res.Ok {
			c.Respond(res.Response())
			return
		}
		//distinguish between new and changing
		var addtext string
		if old != types.DefaultMark {
			addtext = "to edit the"
		} else {
			addtext = "a new"
		}
		// Respond
		c.Respond(sevcord.NewMessage(fmt.Sprintf("Suggested %s note for **%s** 🖋️", addtext, name)))
	}).Input(sevcord.NewModalInput("New Comment", types.DefaultMark, sevcord.ModalInputStyleParagraph, 2400)))
}

func (e *Elements) MsgSignCmd(c sevcord.Ctx, elem string, mark string) {
	c.Acknowledge()

	var name string
	var old string
	var id int
	err := e.db.QueryRow("SELECT id, name, comment FROM elements WHERE LOWER(name)=$1 AND guild=$2", strings.ToLower(elem), c.Guild()).Scan(&id, &name, &old)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Respond(sevcord.NewMessage("Element **" + elem + "** doesn't exist! " + types.RedCircle))
			return
		} else {
			e.base.Error(c, err)
			return
		}
	}

	// Make poll
	res := e.polls.CreatePoll(c, &types.Poll{
		Kind: types.PollKindComment,
		Data: types.PgData{
			"elem": float64(id),
			"new":  mark,
			"old":  old,
		},
	})
	if !res.Ok {
		c.Respond(res.Response())
		return
	}
	// distinguish between new and changing
	var addtext string
	if old != types.DefaultMark {
		addtext = "to edit the"
	} else {
		addtext = "a new"
	}
	// Respond
	c.Respond(sevcord.NewMessage(fmt.Sprintf("Suggested %s note for **%s** 🖋️", addtext, name)))
}

func (e *Elements) ColorCmd(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	// Check hex code
	code := opts[1].(string)
	val, err := strconv.ParseInt(strings.TrimPrefix(code, "#"), 16, 64)
	if err != nil {
		c.Respond(sevcord.NewMessage("Invalid hex code! " + types.RedCircle))
		return
	}
	if val < 0 || val > 16777215 {
		c.Respond(sevcord.NewMessage("Invalid hex code! " + types.RedCircle))
		return
	}

	// Check element
	var name string
	var old int
	var colorer string
	err = e.db.QueryRow("SELECT name, color,colorer FROM elements WHERE id=$1 AND guild=$2", opts[0].(int64), c.Guild()).Scan(&name, &old, &colorer)
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Make poll
	res := e.polls.CreatePoll(c, &types.Poll{
		Kind: types.PollKindColor,
		Data: types.PgData{
			"elem": float64(opts[0].(int64)),
			"new":  float64(val),
			"old":  float64(old),
		},
	})
	if !res.Ok {
		c.Respond(res.Response())
		return
	}
	// distinguish between new and changing
	var addtext string
	if colorer != "" {
		addtext = "to edit the"
	} else {
		addtext = "a new"
	}
	// Respond
	c.Respond(sevcord.NewMessage(fmt.Sprintf("Suggested %s color for **%s** 🎨", addtext, name)))
}
