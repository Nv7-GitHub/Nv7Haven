package elements

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
)

func (e *Elements) InfoSlashCmd(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	e.Info(c, int(opts[0].(int64)))
}

func (e *Elements) InfoMsgCmd(c sevcord.Ctx, val string) {
	e.base.IncrementCommandStat(c, "info")

	c.Acknowledge()
	var id int
	err := e.db.QueryRow("SELECT id FROM elements WHERE guild=$1 AND LOWER(name)=$2", c.Guild(), strings.ToLower(val)).Scan(&id)
	if err != nil {
		e.base.Error(c, err, "Element **"+val+"** doesn't exist!")
		return
	}
	e.Info(c, id)
}

func (e *Elements) Info(c sevcord.Ctx, el int) {
	// Get element
	var elem types.Element
	err := e.db.Get(&elem, "SELECT * FROM elements WHERE id=$1 AND guild=$2", el, c.Guild())
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Check if you have
	description := "**Mark**\n" + elem.Comment
	var have bool
	err = e.db.QueryRow(`SELECT $1=ANY(inv) FROM inventories WHERE guild=$2 AND "user"=$3`, elem.ID, c.Guild(), c.Author().User.ID).Scan(&have)
	if err != nil {
		e.base.Error(c, err)
		return
	}
	if have {
		description = "**You have this.**\n\n" + description
	} else {
		description = "**You don't have this.**\n\n" + description
	}

	// Get stats
	var madewith int
	err = e.db.QueryRow("SELECT COUNT(*) FROM combos WHERE result=$1 AND guild=$2", elem.ID, c.Guild()).Scan(&madewith)
	if err != nil {
		e.base.Error(c, err)
		return
	}
	var usedin int
	err = e.db.QueryRow("SELECT COUNT(*) FROM combos WHERE $1=ANY(els) AND guild=$2", elem.ID, c.Guild()).Scan(&usedin)
	if err != nil {
		e.base.Error(c, err)
		return
	}
	var foundby int
	err = e.db.QueryRow("SELECT COUNT(*) FROM inventories WHERE $1=ANY(inv) AND guild=$2", elem.ID, c.Guild()).Scan(&foundby)
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Element ID
	description = fmt.Sprintf("Element **#%d**\n", elem.ID) + description

	// Embed
	emb := sevcord.NewEmbed().
		Title(elem.Name+" Info").
		Description(description).
		Color(elem.Color).
		AddField("Creator", fmt.Sprintf("<@%s>", elem.Creator), true).
		AddField("Created On", fmt.Sprintf("<t:%d>", elem.CreatedOn.Unix()), true).
		AddField("Tree Size", humanize.Comma(int64(elem.TreeSize)), true).
		AddField("Made With", humanize.Comma(int64(madewith)), true).
		AddField("Used In", humanize.Comma(int64(usedin)), true).
		AddField("Found By", humanize.Comma(int64(foundby)), true)

	// Optional things
	if elem.Image != "" {
		emb = emb.Thumbnail(elem.Image)
	}
	if elem.Commenter != "" {
		emb = emb.AddField("Commenter", fmt.Sprintf("<@%s>", elem.Commenter), true)
	}
	if elem.Colorer != "" {
		emb = emb.AddField("Colorer", fmt.Sprintf("<@%s>", elem.Colorer), true)
	}
	if elem.Imager != "" {
		emb = emb.AddField("Imager", fmt.Sprintf("<@%s>", elem.Imager), true)
	}

	// Respond
	msg := sevcord.NewMessage("").AddEmbed(emb)
	c.Respond(msg)
}
