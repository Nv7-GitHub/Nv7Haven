package elements

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/eodb"
	"github.com/Nv7-Github/Nv7Haven/eod/eodsort"
	"github.com/Nv7-Github/Nv7Haven/eod/trees"
	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/Nv7-Github/Nv7Haven/eod/util"
	"github.com/bwmarrin/discordgo"
)

const catInfoCount = 3
const catInfoCountExpanded = 6

func (b *Elements) SortCmd(sort string, postfix bool, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		return
	}

	rsp.Acknowledge()

	type sortItem struct {
		id   int
		name string
	}
	items := make([]sortItem, len(db.Elements))

	db.RLock()
	for i, el := range db.Elements {
		items[i] = sortItem{el.ID, el.Name}
	}
	db.RUnlock()

	eodsort.Sort(items, len(items), func(index int) int {
		return items[index].id
	}, func(index int) string {
		return items[index].name
	}, func(index int, val string) {
		items[index].name = val
	}, sort, m.Author.ID, db, postfix)

	text := make([]string, len(items))
	for i, val := range items {
		text[i] = val.name
	}

	b.base.NewPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchInv,
		Title:      db.Config.LangProperty("ElemSort", nil),
		PageGetter: b.base.InvPageGetter,
		Items:      text,
	}, m, rsp)
}

func newCmpCollapsed(db *eodb.DB) discordgo.ActionsRow {
	return discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    db.Config.LangProperty("InfoExpand", nil),
				CustomID: "expand",
				Style:    discordgo.SuccessButton,
				Emoji: discordgo.ComponentEmoji{
					Name:     "expand",
					ID:       "932829946706006046",
					Animated: false,
				},
			},
		},
	}
}

func newCmpExpanded(db *eodb.DB) discordgo.ActionsRow {
	return discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    db.Config.LangProperty("InfoCollapse", nil),
				CustomID: "collapse",
				Style:    discordgo.DangerButton,
				Emoji: discordgo.ComponentEmoji{
					Name:     "collapse",
					ID:       "932831405640155176",
					Animated: false,
				},
			},
		},
	}
}

type infoComponent struct {
	Expand   *discordgo.MessageEmbed
	Collapse *discordgo.MessageEmbed
	db       *eodb.DB
	Expanded bool

	b *Elements
}

func (c *infoComponent) Handler(_ *discordgo.Session, i *discordgo.InteractionCreate) {
	c.Expanded = !c.Expanded
	var emb *discordgo.MessageEmbed
	var cmp discordgo.ActionsRow
	if c.Expanded {
		emb = c.Expand
		cmp = newCmpExpanded(c.db)
	} else {
		emb = c.Collapse
		cmp = newCmpCollapsed(c.db)
	}
	c.b.dg.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{emb},
			Components: []discordgo.MessageComponent{cmp},
		},
	})
}

func (b *Elements) Info(elem string, m types.Msg, rsp types.Rsp) {
	db, res := b.GetDB(m.GuildID)
	if !res.Exists {
		rsp.ErrorMessage(res.Message)
		return
	}

	// Get Element
	el, res := db.GetElementByName(elem)
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

	// Get Categories
	cats := b.base.ElemCategories(el.ID, db, false)

	// Make text for collapsed
	catTxt := &strings.Builder{}
	for i := 0; i < catInfoCount && i < len(cats); i++ {
		catTxt.WriteString(cats[i])
		if i != catInfoCount-1 && i != len(cats)-1 {
			catTxt.WriteString(", ")
		}
	}
	if len(cats) > catInfoCount {
		catTxt.WriteString(db.Config.LangProperty("InfoAdditionalElemCats", len(cats)-catInfoCount))
	}

	// Make text for expanded
	catTxtExpanded := &strings.Builder{}
	for i := 0; i < catInfoCountExpanded && i < len(cats); i++ {
		catTxtExpanded.WriteString(cats[i])
		if i != catInfoCountExpanded-1 && i != len(cats)-1 {
			catTxtExpanded.WriteString(", ")
		}
	}
	if len(cats) > catInfoCountExpanded {
		catTxtExpanded.WriteString(db.Config.LangProperty("InfoAdditionalElemCats", len(cats)-catInfoCountExpanded))
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

	if len(el.Comment) == 0 {
		el.Comment = db.Config.LangProperty("DefaultComment", nil)
	}

	shortcomment := el.Comment
	if len(el.Comment) > 100 {
		shortcomment = el.Comment[:99]
	}

	if len(strings.ReplaceAll(shortcomment, "\n", ""))+4 < len(shortcomment) {
		shortcomment = strings.Join(strings.Split(shortcomment, "\n")[:4], "\n")
	}

	if shortcomment != el.Comment {
		shortcomment = strings.TrimSpace(shortcomment) + "..."
	}

	createdOn := fmt.Sprintf("<t:%d>", el.CreatedOn.Unix())
	if el.CreatedOn.Unix() <= 4 {
		createdOn = db.Config.LangProperty("StarterElemCreateTime", nil)
	}

	infoFields := make([]*discordgo.MessageEmbedField, 0)
	if el.Commenter != "" {
		infoFields = append(infoFields, &discordgo.MessageEmbedField{Name: db.Config.LangProperty("InfoCommenter", nil), Value: fmt.Sprintf("<@%s>", el.Commenter), Inline: true})
	}
	if el.Imager != "" {
		infoFields = append(infoFields, &discordgo.MessageEmbedField{Name: db.Config.LangProperty("InfoImager", nil), Value: fmt.Sprintf("<@%s>", el.Imager), Inline: true})
	}
	if el.Colorer != "" {
		infoFields = append(infoFields, &discordgo.MessageEmbedField{Name: db.Config.LangProperty("InfoColorer", nil), Value: fmt.Sprintf("<@%s>", el.Colorer), Inline: true})
	}

	// Get Air/Earth/Fire/Water names
	names := make([]string, 0)
	for i := 1; i <= 4; i++ {
		el, _ := db.GetElement(i)
		names = append(names, el.Name)
	}

	// Make fields
	fullFields := []*discordgo.MessageEmbedField{
		{Name: db.Config.LangProperty("InfoComment", nil), Value: el.Comment, Inline: false},
		{Name: db.Config.LangProperty("InfoCombosUsedIn", nil), Value: strconv.Itoa(el.UsedIn), Inline: true},
		{Name: db.Config.LangProperty("InfoCombosMadeWith", nil), Value: strconv.Itoa(madeby), Inline: true},
		{Name: db.Config.LangProperty("InfoUsersFoundBy", nil), Value: strconv.Itoa(foundby), Inline: true},
		{Name: db.Config.LangProperty("InfoCreator", nil), Value: fmt.Sprintf("<@%s>", el.Creator), Inline: true},
		{Name: db.Config.LangProperty("InfoCreateTime", nil), Value: createdOn, Inline: true},
		{Name: db.Config.LangProperty("InfoColor", nil), Value: util.FormatHex(el.Color), Inline: true},
		{Name: db.Config.LangProperty("InfoTreeSize", nil), Value: strconv.Itoa(tree.Total), Inline: true},
		{Name: db.Config.LangProperty("InfoComplexity", nil), Value: strconv.Itoa(el.Complexity), Inline: true},
		{Name: db.Config.LangProperty("InfoDifficulty", nil), Value: strconv.Itoa(el.Difficulty), Inline: true},
		{Name: names[0], Value: util.FormatBigInt(el.Air), Inline: true},
		{Name: names[1], Value: util.FormatBigInt(el.Earth), Inline: true},
		{Name: names[2], Value: util.FormatBigInt(el.Fire), Inline: true},
		{Name: names[3], Value: util.FormatBigInt(el.Water), Inline: true},
	}
	fullFields = append(fullFields, infoFields...)

	// Collapsed fields
	fields := []*discordgo.MessageEmbedField{
		{Name: db.Config.LangProperty("InfoComment", nil), Value: shortcomment, Inline: false},
		{Name: db.Config.LangProperty("InfoCreator", nil), Value: fmt.Sprintf("<@%s>", el.Creator), Inline: true},
		{Name: db.Config.LangProperty("InfoCreateTime", nil), Value: createdOn, Inline: true},
		{Name: db.Config.LangProperty("InfoTreeSize", nil), Value: strconv.Itoa(tree.Total), Inline: true},
	}

	// Get whether has element
	hasPars := map[string]any{
		"ID":   el.ID,
		"User": m.Author.ID,
	}
	has := db.Config.LangProperty("InfoElemIDUserHasElem", hasPars)
	inv := db.GetInv(m.Author.ID)
	exists := inv.Contains(el.ID)
	if !exists {
		has = db.Config.LangProperty("InfoElemIDUserNoHasElem", hasPars)
	}

	// Embed
	emb := &discordgo.MessageEmbed{
		Title:       db.Config.LangProperty("InfoTitle", el.Name),
		Description: has,
		Fields:      fields,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: el.Image,
		},
		Color: el.Color,
	}
	if !exists {
		fullFields = append(fullFields, &discordgo.MessageEmbedField{
			Name:   db.Config.LangProperty("InfoElemProgress", nil),
			Value:  fmt.Sprintf("%s%%", util.FormatFloat(float32(tree.Found)/float32(tree.Total)*100, 2)),
			Inline: true,
		})
	}
	if len(cats) > 0 {
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{Name: db.Config.LangProperty("InfoElemCats", nil), Value: catTxt.String(), Inline: false})
		fullFields = append(fullFields, &discordgo.MessageEmbedField{Name: db.Config.LangProperty("InfoElemCats", nil), Value: catTxtExpanded.String(), Inline: false})
	}

	if m.Author.ID == "567132457820749842" {
		for _, elem := range types.StarterElements {
			if elem.Name == el.Name {
				emb.Thumbnail.URL = elem.Image
			}
		}
	}

	// Collapsed
	full := *emb
	full.Fields = fullFields

	if len(el.Comment) > 1024 {
		full.Fields = full.Fields[1:]
		full.Description = fmt.Sprintf("%s\n\n**%s**\n%s", emb.Description, db.Config.LangProperty("InfoComment", nil), el.Comment)
	}

	// Send
	msgId := rsp.RawEmbed(emb, newCmpCollapsed(db))

	// Component
	cmp := &infoComponent{
		b:        b,
		Expand:   &full,
		Collapse: emb,
		db:       db,
	}

	// Data
	data, _ := b.GetData(m.GuildID)
	data.SetMsgElem(msgId, el.ID)
	data.AddComponentMsg(msgId, cmp)
}

func (b *Elements) InfoCmd(elem string, m types.Msg, rsp types.Rsp) {
	elem = strings.TrimSpace(elem)
	b.Info(elem, m, rsp)
}
