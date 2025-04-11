package base

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/Nv7-Github/sevcord/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/lib/pq"
)

func (b *Base) Give(c sevcord.Ctx, opts []any) {
	c.Acknowledge()

	user := opts[0].(*discordgo.User).ID
	q, ok := b.CalcQuery(c, opts[1].(string))
	if !ok {
		return
	}

	// Combine to inv
	_, err := b.db.Exec(`UPDATE inventories SET inv=inv | $1 WHERE guild=$2 AND "user"=$3`, pq.Array(q.Elements), c.Guild(), user)
	if err != nil {
		b.Error(c, err)
		return
	}

	// Respond
	c.Respond(sevcord.NewMessage(fmt.Sprintf("Succesfully gave elements to <@%s>!", user)))
}
func (b *Base) GiveFile(c sevcord.Ctx, opts []any) {

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

	var elems pq.Int32Array
	for i := 0; i < len(elemsstr); i++ {
		if elemsstr[i] == "" {
			continue
		}
		id, _ := strconv.Atoi(strings.TrimPrefix(elemsstr[i], "#"))
		elems = append(elems, int32(id))
	}
	if len(elems) == 0 {
		return
	}
	// Combine to inv
	_, err = b.db.Exec(`UPDATE inventories SET inv=inv | $1 WHERE guild=$2 AND "user"=$3`, elems, c.Guild(), user)
	if err != nil {
		b.Error(c, err)
		return
	}
	c.Respond(sevcord.NewMessage(fmt.Sprintf("Succesfully gave elements to <@%s>!", user)))

}
