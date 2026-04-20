package base

import (
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/sevcord/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func (b *Base) Take(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	user := opts[0].(*discordgo.User).ID
	q, ok := b.CalcQuery(c, opts[1].(string))
	if !ok {
		return
	}
	var tx *sqlx.Tx
	var err error
	tx, err = b.db.Beginx()
	if err != nil {
		b.Error(c, err)
		return
	}
	// remove from inv
	_, err = tx.Exec(`UPDATE inventories SET inv=inv-$1 WHERE guild=$2 AND "user"=$3`, pq.Array(q.Elements), c.Guild(), user)
	if err != nil {
		tx.Rollback()
		b.Error(c, err)
		return
	}
	//check if inv is 0 length after operation
	var len int
	b.db.Get(&len, `SELECT cardinality(inv) from inventories WHERE guild=$1 AND "user"=$2`, c.Guild(), user)
	if len == 0 {
		tx.Rollback()
		c.Respond(sevcord.NewMessage("Cannot remove all elements from the inventory! " + types.RedCircle))
		return
	}

	// Respond
	c.Respond(sevcord.NewMessage(fmt.Sprintf("Succesfully removed elements from <@%s>!", user)))
}
func (b *Base) Set(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	user := opts[0].(*discordgo.User).ID
	q, ok := b.CalcQuery(c, opts[1].(string))
	if !ok {
		return
	}
	if len(q.Elements) == 0 {
		c.Respond(sevcord.NewMessage("Can't set an invetory to empty! " + types.RedCircle))
		return
	}

	// set to inv
	_, err := b.db.Exec(`UPDATE inventories SET inv=$1 WHERE guild=$2 AND "user"=$3`, pq.Array(q.Elements), c.Guild(), user)
	if err != nil {
		b.Error(c, err)
		return
	}

	// Respond
	c.Respond(sevcord.NewMessage(fmt.Sprintf("Succesfully set elements to <@%s>!", user)))
}
func (b *Base) SetFile(c sevcord.Ctx, opts []any) {
	c.Acknowledge()
	c.Acknowledge()
	user := opts[0].(*discordgo.User).ID
	file := opts[1].(*sevcord.SlashCommandAttachment)
	URL := file.URL
	resp, err := http.Get(URL)

	if err != nil {
		b.Error(c, err)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		b.Error(c, err)
		return
	}
	list := string(data)
	elemsstr := strings.Split(list, "\n")

	var elems []int
	for i := 0; i < len(elemsstr); i++ {

		if elemsstr[i] == "" {
			continue
		}
		id, err := strconv.Atoi(elemsstr[i])
		if err != nil {
			c.Respond(sevcord.NewMessage("Invalid element ID! " + types.RedCircle))
			return
		}
		if !slices.Contains(elems, id) {
			elems = append(elems, int(id))
		}

	}
	if len(elems) == 0 {
		c.Respond(sevcord.NewMessage("Cannot set an inventory to empty! " + types.RedCircle))
		return
	}
	// set to inv
	_, err = b.db.Exec(`UPDATE inventories SET inv=$1 WHERE guild=$2 AND "user"=$3`, pq.Array(elems), c.Guild(), user)
	if err != nil {
		b.Error(c, err)
		return
	}

	// Respond
	c.Respond(sevcord.NewMessage(fmt.Sprintf("Succesfully set elements to <@%s>!", user)))

}
