package elements

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/base"
	"github.com/Nv7-Github/Nv7Haven/eod/eodsort"
	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/bwmarrin/discordgo"
)

const catInfoCount = 3

func (b *Elements) SortCmd(sort string, postfix bool, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}

	rsp.Acknowledge()

	ids := make([]int, len(db.Elements))
	i := 0
	db.RLock()
	for _, el := range db.Elements {
		ids[i] = el.ID
		i++
	}
	db.RUnlock()

	var items []string
	if postfix {
		items = eodsort.SortElemList(ids, sort, db)
	} else {
		items = eodsort.SortElemList(ids, sort, db, true)
	}

	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      "Element Sort",
		PageGetter: b.base.InvPageGetter,
		Items:      items,
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
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	// Get Element name from ID
	var el types.Element
	if isId {
		el, res = db.GetElement(id)
		if !res.Exists {
			rsp.ErrorMessage(res.Message)
		}
	}

	if base.IsFoolsMode && !base.IsFool(elem) {
		rsp.ErrorMessage(base.MakeFoolResp(elem))
		return
	}

	// Get Element
	if !isId {
		el, res = db.GetElementByName(elem)
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
	}
	rsp.Acknowledge()

	// Get whether has element
	has := ""
	inv := db.GetInv(m.Author.ID)
	exists := inv.Contains(el.ID)
	if !exists {
		has = "don't "
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
	db.RLock()
	for _, comb := range db.Combos() {
		if comb == el.ID {
			madeby++
		}
	}

	// Get foundby
	foundby := 0
	for _, inv := range db.Invs() {
		if inv.Contains(el.ID) {
			foundby++
		}
	}
	db.RUnlock()

	suc, msg, tree := trees.CalcElemInfo(el.ID, m.Author.ID, db)
	if !suc {
		rsp.ErrorMessage(msg)
		return
	}
	// Make empty marks into "None"
	mark := el.Comment
	if len(strings.trimSpace(mark)) < 1 {
		mark = "None"
	}

	emb := &discordgo.MessageEmbed{
		Title:       el.Name + " Info",
		Description: fmt.Sprintf("Element **#%d**\n<@%s> **You %shave this.**", el.ID, m.Author.ID, has),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Mark", Value: mark, Inline: false},
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

	data, _ := b.GetData(m.GuildID)
	data.SetMsgElem(msgId, el.ID)
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
