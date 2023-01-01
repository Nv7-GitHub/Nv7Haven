package elements

import (
	"time"

	"github.com/Nv7-Github/sevcord/v2"
)

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
	c.Respond(sevcord.NewMessage("Successfully edited element " + nameV + "! ✅"))
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
	e.editCmd(c, opts, "creator")
}

func (e *Elements) EditElementCreatedonCmd(c sevcord.Ctx, opts []any) {
	opts[1] = time.Unix(opts[1].(int64), 0)
	e.editCmd(c, opts, "createdon", "creation date")
}
