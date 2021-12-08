package categories

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/eodsort"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

func (b *Categories) CategoriesCmd(elem string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}

	rsp.Acknowledge()

	el, res := db.GetElementByName(elem)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	// Get Categories
	catsMap := make(map[catSortInfo]types.Empty)
	db.RLock()
	for _, cat := range db.Cats() {
		_, exists := cat.Elements[el.ID]
		if exists {
			catsMap[catSortInfo{
				Name: cat.Name,
				Cnt:  len(cat.Elements),
			}] = types.Empty{}
		}
	}
	db.RUnlock()
	cats := make([]catSortInfo, len(catsMap))
	i := 0
	for k := range catsMap {
		cats[i] = k
		i++
	}

	// Sort categories by count
	sort.Slice(cats, func(i, j int) bool {
		return cats[i].Cnt > cats[j].Cnt
	})

	// Convert to array
	out := make([]string, len(cats))
	for i, cat := range cats {
		out[i] = cat.Name
	}

	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      fmt.Sprintf("%s Categories (%d)", el.Name, len(out)),
		PageGetter: b.base.InvPageGetter,
		Items:      out,
		User:       m.Author.ID,
	}, m, rsp)
}

func (b *Categories) DownloadCatCmd(catName string, sort string, postfix bool, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}

	cat, res := db.GetCat(catName)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	db.RLock()
	elems := make([]int, len(cat.Elements))
	i := 0

	for elem := range cat.Elements {
		elems[i] = elem
		i++
	}
	db.RUnlock()

	var vals []string
	if postfix {
		vals = eodsort.SortElemList(elems, sort, db)
	} else {
		vals = eodsort.SortElemList(elems, sort, db, true)
	}

	out := &strings.Builder{}
	for _, elem := range vals {
		out.WriteString(elem + "\n")
	}
	buf := strings.NewReader(out.String())

	channel, err := b.dg.UserChannelCreate(m.Author.ID)
	if rsp.Error(err) {
		return
	}

	_, err = b.dg.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Category **%s**:", cat.Name),
		Files: []*discordgo.File{
			{
				Name:        "cat.txt",
				ContentType: "text/plain",
				Reader:      buf,
			},
		},
	})
	if rsp.Error(err) {
		return
	}
	rsp.Message("Sent category in DMs!")
}
