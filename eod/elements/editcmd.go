package elements

import (
	"fmt"
	"log"
	"time"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/bwmarrin/discordgo"
)

func (e *Elements) editNewsMessage(c sevcord.Ctx, message string) {
	var news string
	err := e.db.QueryRow(`SELECT news FROM config WHERE guild=$1`, c.Guild()).Scan(&news)
	if err != nil {
		log.Println("news err", err)
		return
	}
	_, err = c.Dg().ChannelMessageSend(news, fmt.Sprintf("🔨 "+message))
	if err != nil {
		log.Println("news err", err)
	}
}

func (e *Elements) editCmd(c sevcord.Ctx, opts []any, field string, name ...string) {
	c.Acknowledge()
	_, err := e.db.Exec("UPDATE elements SET "+field+"=$3 WHERE id=$1 AND guild=$2", opts[0].(int64), c.Guild(), opts[1])
	if err != nil {
		e.base.Error(c, err)
		return
	}
	nameV := field
	if len(name) > 0 {
		nameV = name[0]
	}
	// Get element
	var elem types.Element
	err = e.db.Get(&elem, "SELECT * FROM elements WHERE id=$1 AND guild=$2", int(opts[0].(int64)), c.Guild())
	if err != nil {
		e.base.Error(c, err)
		return
	}
	nameE := elem.Name
	nameC := elem.Creator
	c.Respond(sevcord.NewMessage("Successfully edited element " + nameV + "! ✅"))
	e.editNewsMessage(c, fmt.Sprintf("Edited Element %s - **%s** (By <@%s>) (element by <@%s>) - Element **#%d** ", util.Capitalize(nameV), nameE, c.Author().User.ID, nameC, opts[0].(int64)))
}

func (e *Elements) EditElementNameCmd(c sevcord.Ctx, opts []any) {
	e.editCmd(c, opts, "name")
}

func (e *Elements) EditElementImageCmd(c sevcord.Ctx, opts []any) {
	e.editCmd(c, opts, "image")
}

func (e *Elements) EditElementColorCmd(c sevcord.Ctx, opts []any) {
	e.editCmd(c, opts, "color")
}

func (e *Elements) EditElementCommentCmd(c sevcord.Ctx, opts []any) {
	e.editCmd(c, opts, "comment")
}

func (e *Elements) EditElementCreatorCmd(c sevcord.Ctx, opts []any) {
	opts[1] = opts[1].(*discordgo.User).ID
	e.editCmd(c, opts, "creator")
}

func (e *Elements) EditElementCreatedonCmd(c sevcord.Ctx, opts []any) {
	opts[1] = time.Unix(opts[1].(int64), 0)
	e.editCmd(c, opts, "createdon", "creation date")
}
