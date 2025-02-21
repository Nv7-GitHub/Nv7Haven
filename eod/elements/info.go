package elements

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/dustin/go-humanize"
)

func (e *Elements) InfoSlashCmd(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	e.Info(c, int(opts[0].(int64)))
}

func (e *Elements) InfoSlashCmdName(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	var id int
	err := e.db.QueryRow("SELECT id FROM elements WHERE guild=$1 AND LOWER(name)=$2", c.Guild(), strings.ToLower(opts[0].(string))).Scan(&id)
	if err != nil {
		e.base.Error(c, err, "Element **"+opts[0].(string)+"** doesn't exist!")
		return
	}
	e.Info(c, id)
}

func (e *Elements) InfoMsgCmd(c sevcord.Ctx, val string) {
	c.Acknowledge()

	var id int
	if strings.HasPrefix(val, "#") {
		var err error
		id, err = strconv.Atoi(val[1:])
		if err != nil {
			c.Respond(sevcord.NewMessage("Invalid element ID! " + types.RedCircle))
			return
		}
	} else {
		err := e.db.QueryRow("SELECT id FROM elements WHERE guild=$1 AND LOWER(name)=$2", c.Guild(), strings.ToLower(val)).Scan(&id)
		if err != nil {
			e.base.Error(c, err, "Element **"+val+"** doesn't exist!")
			return
		}
	}
	e.Info(c, id)
}

const catInfoCount = 3

func (e *Elements) Info(c sevcord.Ctx, el int) {
	c.Acknowledge()
	// Get element
	c.Acknowledge()
	var elem types.Element
	err := e.db.Get(&elem, "SELECT * FROM elements WHERE id=$1 AND guild=$2", el, c.Guild())
	if err != nil {
		e.base.Error(c, err)
		return
	}

	// Check if you have
	description := "**📝 Mark**\n" + elem.Comment
	var have bool
	err = e.db.QueryRow(`SELECT $1=ANY(inv) FROM inventories WHERE guild=$2 AND "user"=$3`, elem.ID, c.Guild(), c.Author().User.ID).Scan(&have)
	if err != nil {
		e.base.Error(c, err)
		return
	}
	if have {
		description = "📫 **You have this.**\n\n" + description
	} else {
		description = "📪 **You don't have this.**\n\n" + description
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

	// Do categories
	var categories []string
	err = e.db.Select(&categories, "SELECT name FROM categories WHERE $1=ANY(elements) AND guild=$2", elem.ID, c.Guild())
	if err != nil {
		e.base.Error(c, err)
		return
	}
	cnt := len(categories)
	if cnt > catInfoCount {
		categories = categories[:catInfoCount]
		categories = append(categories, fmt.Sprintf("and %d more...", cnt-catInfoCount))
	}

	// Progress
	var treesize, found int
	err = e.db.QueryRow(`WITH RECURSIVE parents AS (
		(select parents, id from elements WHERE id=$2 and guild=$1)
	UNION
		(SELECT b.parents, b.id FROM elements b INNER JOIN parents p ON b.id=ANY(p.parents) where guild=$1)
	) SELECT COUNT(id), (SELECT COUNT(el) FROM (SELECT UNNEST(inv) el FROM inventories WHERE guild=$1 AND "user"=$3) sub INNER JOIN parents ON parents.id=sub.el) FROM parents`, c.Guild(), el, c.Author().User.ID).Scan(&treesize, &found)
	if err != nil {
		e.base.Error(c, err)
		return
	}
	var dbtreesize int
	e.db.QueryRow(`SELECT treesize FROM elements WHERE id=$1 AND guild=$2`, elem.ID, c.Guild()).Scan(&dbtreesize)
	if dbtreesize != treesize {
		e.db.Exec(`UPDATE elements SET treesize=$3 WHERE id=$1 AND guild=$2`, elem.ID, c.Guild(), treesize)
	}

	// Element ID
	description = fmt.Sprintf("Element **#%d**\n", elem.ID) + description

	// Embed
	emb := sevcord.NewEmbed().
		Title(elem.Name+" Info").
		Description(description).
		Color(elem.Color).
		AddField("🧑 Creator", fmt.Sprintf("<@%s>", elem.Creator), true).
		AddField("📅 Created On", fmt.Sprintf("<t:%d>", elem.CreatedOn.Unix()), true).
		AddField("🌲 Tree Size", humanize.Comma(int64(treesize)), true).
		AddField("📊 Progress", humanize.FormatFloat("##.#", float64(found)/float64(treesize)*100)+"%", true).
		AddField("🔨 Made With", humanize.Comma(int64(madewith)), true).
		AddField("🧰 Used In", humanize.Comma(int64(usedin)), true).
		AddField("🔍 Found By", humanize.Comma(int64(foundby)), true).
		AddField("🎨 Color", util.FormatHex(elem.Color), true)

	// Optional things
	if elem.Image != "" {
		emb = emb.Thumbnail(elem.Image)
	}
	if elem.Commenter != "" {
		emb = emb.AddField("💬 Commenter", fmt.Sprintf("<@%s>", elem.Commenter), true)
	}
	if elem.Colorer != "" {
		emb = emb.AddField("🖌️ Colorer", fmt.Sprintf("<@%s>", elem.Colorer), true)
	}
	if elem.Imager != "" {
		emb = emb.AddField("🖼️ Imager", fmt.Sprintf("<@%s>", elem.Imager), true)
	}
	if len(categories) > 0 {
		emb = emb.AddField("📁 Categories", strings.Join(categories, ", "), false)
	}

	// Respond
	msg := sevcord.NewMessage("").AddEmbed(emb)
	c.Respond(msg)
}
