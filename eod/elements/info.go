package elements

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/bwmarrin/discordgo"
)

const catInfoCount = 3

func (b *Elements) SortCmd(sort string, m types.Msg, rsp types.Rsp) {
	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		return
	}

	rsp.Acknowledge()

	items := make([]string, len(dat.Elements))
	i := 0
	for _, el := range dat.Elements {
		items[i] = el.Name
		i++
	}
	util.SortElemList(items, sort, dat)

	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      "Element Sort",
		PageGetter: b.base.InvPageGetter,
		Items:      items,
		Length:     len(dat.Elements),
	}, m, rsp)
}

type catSortInfo struct {
	Name string
	Cnt  int
}

func (b *Elements) Info(elem string, id int, isId bool, m types.Msg, rsp types.Rsp) {
	if len(elem) == 0 && !isId {
		return
	}
	b.lock.RLock()
	dat, exists := b.dat[m.GuildID]
	b.lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("Guild isn't setup yet!")
		return
	}

	// Get Element name from ID
	if isId {
		if id > len(dat.Elements) {
			rsp.ErrorMessage(fmt.Sprintf("Element **#%d** doesn't exist!", id))
			return
		}

		hasFound := false
		dat.Lock.RLock()
		for _, el := range dat.Elements {
			if el.ID == id {
				hasFound = true
				elem = el.Name
			}
		}
		dat.Lock.RUnlock()

		if !hasFound {
			rsp.ErrorMessage(fmt.Sprintf("Element **#%d** doesn't exist!", id))
			return
		}
	}

	if base.IsFoolsMode && !base.IsFool(elem) {
		rsp.ErrorMessage(base.MakeFoolResp(elem))
		return
	}

	// Get Element
	el, res := dat.GetElement(elem)
	if !res.Exists {
		// If what you said was "????", then stop
		if strings.Contains(elem, "?") {
			isValid := false
			for _, letter := range elem {
				if letter != '?' {
					isValid = true
					break
				}
			}
			if !isValid {
				return
			}
		}
		rsp.ErrorMessage(res.Message)
		return
	}
	rsp.Acknowledge()

	// Get whether has element
	has := ""
	exists = false
	inv, res := dat.GetInv(m.Author.ID, true)
	if res.Exists {
		exists = inv.Elements.Contains(el.Name)
	}
	if !exists {
		has = "don't "
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

	// Sort by count
	sort.Slice(cats, func(i, j int) bool {
		return cats[i].Cnt > cats[j].Cnt
	})

	// Make text
	catTxt := &strings.Builder{}
	for i := 0; i < catInfoCount && i < len(cats); i++ {
		catTxt.WriteString(cats[i].Name)
		if i != catInfoCount-1 && i != len(cats)-1 {
			catTxt.WriteString(", ")
		}
	}
	if len(cats) > catInfoCount {
		fmt.Fprintf(catTxt, ", and %d more...", len(cats)-catInfoCount)
	}

	// Get Madeby
	madeby := 0
	for _, comb := range dat.Combos {
		if strings.EqualFold(comb, el.Name) {
			madeby++
		}
	}

	// Get foundby
	foundby := 0
	for _, inv := range dat.Inventories {
		if inv.Elements.Contains(el.Name) {
			foundby++
		}
	}

	suc, msg, tree := trees.CalcElemInfo(elem, m.Author.ID, dat)
	if !suc {
		rsp.ErrorMessage(msg)
		return
	}

	emb := &discordgo.MessageEmbed{
		Title:       el.Name + " Info",
		Description: fmt.Sprintf("Element **#%d**\n<@%s> **You %shave this.**", el.ID, m.Author.ID, has),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Mark", Value: el.Comment, Inline: false},
			{Name: "Used In", Value: strconv.Itoa(el.UsedIn), Inline: true},
			{Name: "Made With", Value: strconv.Itoa(madeby), Inline: true},
			{Name: "Found By", Value: strconv.Itoa(foundby), Inline: true},
			{Name: "Created By", Value: fmt.Sprintf("<@%s>", el.Creator), Inline: true},
			{Name: "Created On", Value: fmt.Sprintf("<t:%d>", el.CreatedOn.Unix()), Inline: true},
			{Name: "Color", Value: util.FormatHex(el.Color), Inline: true},
			{Name: "Tree Size", Value: strconv.Itoa(tree.Total), Inline: true},
			{Name: "Complexity", Value: strconv.Itoa(el.Complexity), Inline: true},
			{Name: "Difficulty", Value: strconv.Itoa(el.Difficulty), Inline: true},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: el.Image,
		},
		Color: el.Color,
	}
	if m.Author.ID == "567132457820749842" {
		for _, elem := range base.StarterElements {
			if elem.Name == el.Name {
				emb.Thumbnail.URL = elem.Image
			}
		}
	}
	if has != "" {
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{
			Name:   "Progress",
			Value:  fmt.Sprintf("%s%%", util.FormatFloat(float32(tree.Found)/float32(tree.Total)*100, 2)),
			Inline: true,
		})
	}
	if len(cats) > 0 {
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{Name: "Categories", Value: catTxt.String(), Inline: false})
	}
	if len(el.Comment) > 1024 {
		emb.Fields = emb.Fields[1:]
		emb.Description = fmt.Sprintf("%s\n\n**Mark**\n%s", emb.Description, el.Comment)
	}

	msgId := rsp.RawEmbed(emb)
	dat.SetMsgElem(msgId, el.Name)
}

func (b *Elements) InfoCmd(elem string, m types.Msg, rsp types.Rsp) {
	elem = strings.TrimSpace(elem)
	if elem[0] == '#' {
		number, err := strconv.Atoi(elem[1:])
		if err != nil {
			rsp.ErrorMessage("Invalid Element ID!")
			return
		}
		b.Info("", number, true, m, rsp)
		return
	}
	b.Info(elem, 0, false, m, rsp)
}
