package eod

import (
	"database/sql"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/bwmarrin/discordgo"
)

const catInfoCount = 3

var infoChoices []*discordgo.ApplicationCommandOptionChoice
var infoQuerys = map[string]string{
	"Name":         "SELECT name FROM eod_elements WHERE guild=? ORDER BY name %s LIMIT ? OFFSET ?",
	"Date Created": "SELECT name FROM eod_elements WHERE guild=? ORDER BY createdon %s LIMIT ? OFFSET ?",
	"Complexity":   "SELECT name FROM eod_elements WHERE guild=? ORDER BY complexity %s LIMIT ? OFFSET ?",
	"Difficulty":   "SELECT name FROM eod_elements WHERE guild=? ORDER BY difficulty %s LIMIT ? OFFSET ?",
	"Used In":      `SELECT name FROM eod_elements WHERE guild=? ORDER BY usedin %s LIMIT ? OFFSET ?`,
	// Ones below are commented out due to being extremely slow
	//"Made By":      "SELECT name FROM eod_elements WHERE guild=? ORDER BY (SELECT COUNT(1) AS cnt FROM eod_combos WHERE elem3 LIKE name AND guild=?) %s LIMIT ? OFFSET ?",
	//"Found By":     `SELECT name FROM eod_elements WHERE guild=?  ORDER BY (SELECT COUNT(1) as cnt FROM eod_inv WHERE guild=? AND (JSON_EXTRACT(inv, CONCAT('$."', LOWER(name), '"')) IS NOT NULL)) %s LIMIT ? OFFSET ?`,
	"Creator": "SELECT name FROM eod_elements WHERE guild=? ORDER BY creator %s LIMIT ? OFFSET ?",
	"Length":  `SELECT name FROM eod_elements WHERE guild=? ORDER BY LENGTH(name) %s LIMIT ? OFFSET ?`,
}

func (b *EoD) initInfoChoices() {
	infoChoices = make([]*discordgo.ApplicationCommandOptionChoice, len(infoQuerys))
	i := 0
	for k := range infoQuerys {
		infoChoices[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  k,
			Value: k,
		}
		i++
	}
}

func (b *EoD) sortPageGetter(p types.PageSwitcher) (string, int, int, error) {
	length := int(math.Floor(float64(p.Length-1) / float64(p.PageLength)))
	if p.PageLength*p.Page > (p.Length - 1) {
		return "", 0, length, nil
	}
	if p.Page < 0 {
		return "", length, length, nil
	}
	var res *sql.Rows
	var err error
	cnt := strings.Count(p.Query, "?")
	if cnt == 3 {
		res, err = b.db.Query(p.Query, p.Guild, p.PageLength, p.Page*p.PageLength)
		if err != nil {
			return "", length, length, err
		}
	} else {
		res, err = b.db.Query(p.Query, p.Guild, p.Guild, p.PageLength, p.Page*p.PageLength)
		if err != nil {
			return "", length, length, err
		}
	}
	defer res.Close()
	out := ""
	var name string
	for res.Next() {
		err = res.Scan(&name)
		if err != nil {
			return "", length, length, err
		}
		out += name + "\n"
	}
	return out, p.Page, length, nil
}

func (b *EoD) sortCmd(query string, order bool, m types.Msg, rsp types.Rsp) {
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		return
	}
	quer, exists := infoQuerys[query]
	if !exists {
		rsp.ErrorMessage("Invalid query type!")
		return
	}
	ord := "DESC"
	if order {
		ord = "ASC"
	}
	quer = fmt.Sprintf(quer, ord)
	b.newPageSwitcher(types.PageSwitcher{
		Kind:       types.PageSwitchElemSort,
		Title:      "Element Sort",
		PageGetter: b.sortPageGetter,
		Query:      quer,
		Length:     len(dat.ElemCache),
	}, m, rsp)
}

type catSortInfo struct {
	Name string
	Cnt  int
}

func (b *EoD) info(elem string, id int, isId bool, m types.Msg, rsp types.Rsp) {
	if len(elem) == 0 {
		return
	}
	lock.RLock()
	dat, exists := b.dat[m.GuildID]
	lock.RUnlock()
	if !exists {
		rsp.ErrorMessage("Guild isn't setup yet!")
		return
	}

	// Get Element name from ID
	if isId {
		if id > len(dat.ElemCache) {
			rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", elem))
			return
		}

		hasFound := false
		dat.Lock.RLock()
		for _, el := range dat.ElemCache {
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

	if isFoolsMode && !isFool(elem) {
		rsp.ErrorMessage(makeFoolResp(elem))
		return
	}

	// Get Element
	dat.Lock.RLock()
	el, exists := dat.ElemCache[strings.ToLower(elem)]
	dat.Lock.RUnlock()
	if !exists {
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
		rsp.ErrorMessage(fmt.Sprintf("Element **%s** doesn't exist!", elem))
		return
	}
	rsp.Acknowledge()

	// Get whether has element
	has := ""
	exists = false
	if dat.InvCache != nil {
		dat.Lock.RLock()
		_, exists = dat.InvCache[m.Author.ID]
		if exists {
			_, exists = dat.InvCache[m.Author.ID][strings.ToLower(el.Name)]
		}
		dat.Lock.RUnlock()
	}
	if !exists {
		has = "don't "
	}

	// Get Categories
	catsMap := make(map[catSortInfo]types.Empty)
	dat.Lock.RLock()
	for _, cat := range dat.CatCache {
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

	// Get SQL Stats
	quer := `SELECT a.cnt, b.cnt FROM (SELECT COUNT(1) AS cnt FROM eod_combos WHERE elem3 LIKE ? AND guild=?) a, (SELECT COUNT(1) as cnt FROM eod_inv WHERE guild=? AND (JSON_EXTRACT(inv, CONCAT('$."', LOWER(?), '"')) IS NOT NULL)) b`
	if isASCII(elem) {
		quer = `SELECT a.cnt, b.cnt FROM (SELECT COUNT(1) AS cnt FROM eod_combos WHERE CONVERT(elem3 USING utf8mb4) LIKE CONVERT(? USING utf8mb4) AND guild=CONVERT(? USING utf8mb4) COLLATE utf8mb4_general_ci) a, (SELECT COUNT(1) as cnt FROM eod_inv WHERE guild=? AND (JSON_EXTRACT(inv, CONCAT('$."', LOWER(?), '"')) IS NOT NULL)) b`
	}
	if isWildcard(elem) {
		quer = strings.ReplaceAll(quer, " LIKE ", "=")
	}

	row := b.db.QueryRow(quer, el.Name, el.Guild, el.Guild, el.Name)
	var madeby int
	var foundby int
	err := row.Scan(&madeby, &foundby)
	if rsp.Error(err) {
		return
	}

	emb := &discordgo.MessageEmbed{
		Title:       el.Name + " Info",
		Description: fmt.Sprintf("Element **#%d**\n<@%s> **You %shave this.**", el.ID, m.Author.ID, has),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Creator Mark", Value: el.Comment, Inline: false},
			{Name: "Used In", Value: strconv.Itoa(el.UsedIn), Inline: true},
			{Name: "Made With", Value: strconv.Itoa(madeby), Inline: true},
			{Name: "Found By", Value: strconv.Itoa(foundby), Inline: true},
			{Name: "Created By", Value: fmt.Sprintf("<@%s>", el.Creator), Inline: true},
			{Name: "Created On", Value: fmt.Sprintf("<t:%d>", el.CreatedOn.Unix()), Inline: true},
			{Name: "Complexity", Value: strconv.Itoa(el.Complexity), Inline: true},
			{Name: "Difficulty", Value: strconv.Itoa(el.Difficulty), Inline: true},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: el.Image,
		},
	}
	if len(cats) > 0 {
		emb.Fields = append(emb.Fields, &discordgo.MessageEmbedField{Name: "Categories", Value: catTxt.String(), Inline: false})
	}

	rsp.RawEmbed(emb)
}

func (b *EoD) infoCmd(elem string, m types.Msg, rsp types.Rsp) {
	if elem[0] == '#' {
		number, err := strconv.Atoi(elem[1:])
		if err != nil {
			rsp.ErrorMessage("Invalid Element ID!")
			return
		}
		b.info("", number, true, m, rsp)
		return
	}
	b.info(elem, 0, false, m, rsp)
}
