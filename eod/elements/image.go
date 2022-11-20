package elements

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
)

func (e *Elements) ImageCmd(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	// Check element
	var elem string
	var old string
	err := e.db.QueryRow("SELECT name, image FROM elements WHERE id=$1 AND guild=$2", opts[0].(int64), c.Guild()).Scan(&elem, &old)
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Check image
	if !strings.HasPrefix(opts[1].(*sevcord.SlashCommandAttachment).ContentType, "image") {
		c.Respond(sevcord.NewMessage("The attachment must be an image! " + types.RedCircle))
		return
	}

	// Make poll
	e.polls.CreatePoll(c, &types.Poll{
		Kind: types.PollKindImage,
		Data: types.PollData{
			"elem": float64(opts[0].(int64)),
			"new":  opts[1].(*sevcord.SlashCommandAttachment).URL,
			"old":  old,
		},
	})

	// Respond
	c.Respond(sevcord.NewMessage(fmt.Sprintf("Suggested an image for **%s** 📷", elem)))
}
