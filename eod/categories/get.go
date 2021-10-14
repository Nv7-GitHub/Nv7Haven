package categories

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/bwmarrin/discordgo"
)

func (b *Categories) CategoriesCmd(elem string, m types.Msg, rsp types.Rsp) {
	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		return
	}

	rsp.Acknowledge()

	el, res := dat.GetElement(elem)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	// Get Categories
	catsMap := make(map[catSortInfo]types.Empty)
	dat.Lock.RLock()
	for _, cat := range dat.Categories {
		_, exists := cat.Elements[el.Name]
		if exists {
			catsMap[catSortInfo{
				Name: cat.Name,
				Cnt:  len(cat.Elements),
			}] = types.Empty{}
		}
	}
	dat.Lock.RUnlock()
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
	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		return
	}

	cat, res := dat.GetCategory(catName)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	dat.Lock.RLock()
	elems := make([]string, len(cat.Elements))
	i := 0

	for elem := range cat.Elements {
		elems[i] = elem
		i++
	}
	dat.Lock.RUnlock()

	if postfix {
		util.SortElemList(elems, sort, dat)
	} else {
		util.SortElemList(elems, sort, dat, true)
	}

	out := &strings.Builder{}
	for _, elem := range elems {
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
